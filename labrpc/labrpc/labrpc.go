package labrpc

import (
    "bytes"
    "fmt"
    "lab-rpc/labgob"
    "log"
    "math/rand"
    "reflect"
    "strings"
    "sync"
    "sync/atomic"
    "time"
)

type requestMessage struct {
    endName             interface{} // name of sending ClientEnd
    serviceMethod       string      // e.g. "Raft.AppendEntries"
    argsType            reflect.Type
    args                []byte
    responseMessageChan chan responseMessage
}

type responseMessage struct {
    ok    bool
    reply []byte
}

type ClientEnd struct {
    endName            interface{}         // this end-point's name
    requestMessageChan chan requestMessage // copy of Network.requestMessageChan
    done               chan struct{}       // closed when Network is cleaned up
}

// send an RPC, wait for the reply.
// the return value indicates success; false means that
// no reply was received from the server.
func (clientEnd *ClientEnd) Call(serviceMethod string, args interface{}, reply interface{}) bool {
    req := requestMessage{}
    req.endName = clientEnd.endName
    req.serviceMethod = serviceMethod
    req.argsType = reflect.TypeOf(args)
    req.responseMessageChan = make(chan responseMessage)

    queryBuffer := new(bytes.Buffer)
    queryEncoder := labgob.NewEncoder(queryBuffer)
    if err := queryEncoder.Encode(args); err != nil {
        panic(err)
    }
    req.args = queryBuffer.Bytes()

    //
    // send the request.
    //
    select {
    case clientEnd.requestMessageChan <- req:
        // the request has been sent.
    case <-clientEnd.done:
        // entire Network has been destroyed.
        return false
    }

    //
    // wait for the reply.
    //
    res := <-req.responseMessageChan
    if res.ok {
        replyBuffer := bytes.NewBuffer(res.reply)
        replyEncoder := labgob.NewDecoder(replyBuffer)
        if err := replyEncoder.Decode(reply); err != nil {
            log.Fatalf("ClientEnd.Call(): decode reply: %v\n", err)
        }
        return true
    } else {
        return false
    }
}

type Network struct {
    mu                 sync.Mutex
    reliable           bool
    longDelays         bool                        // pause a long time on send on disabled connection
    longReordering     bool                        // sometimes delay replies a long time
    ends               map[interface{}]*ClientEnd  // endName -> ClientEnd
    enabled            map[interface{}]bool        // endName -> enabled
    servers            map[interface{}]*Server     // serverName -> Server
    connections        map[interface{}]interface{} // endName -> serverName
    requestMessageChan chan requestMessage
    done               chan struct{} // closed when Network is cleaned up
    count              int32         // total RPC count, for statistics
    bytes              int64         // total bytes send, for statistics
}

func MakeNetwork() *Network {
    network := &Network{}
    network.reliable = true
    network.ends = map[interface{}]*ClientEnd{}
    network.enabled = map[interface{}]bool{}
    network.servers = map[interface{}]*Server{}
    network.connections = map[interface{}](interface{}){}
    network.requestMessageChan = make(chan requestMessage)
    network.done = make(chan struct{})

    // single goroutine to handle all ClientEnd.Call()s
    go func() {
        for {
            select {
            case xreq := <-network.requestMessageChan:
                atomic.AddInt32(&network.count, 1)
                atomic.AddInt64(&network.bytes, int64(len(xreq.args)))
                go network.processRequest(xreq)
            case <-network.done:
                return
            }
        }
    }()

    return network
}

func (network *Network) Cleanup() {
    close(network.done)
}

func (network *Network) Reliable(yes bool) {
    network.mu.Lock()
    defer network.mu.Unlock()

    network.reliable = yes
}

func (network *Network) LongReordering(yes bool) {
    network.mu.Lock()
    defer network.mu.Unlock()

    network.longReordering = yes
}

func (network *Network) LongDelays(yes bool) {
    network.mu.Lock()
    defer network.mu.Unlock()

    network.longDelays = yes
}

func (network *Network) readEndNameInfo(endName interface{}) (
    enabled bool, serverName interface{}, server *Server, reliable bool, longreordering bool,
) {
    network.mu.Lock()
    defer network.mu.Unlock()

    enabled = network.enabled[endName]
    serverName = network.connections[endName]
    if serverName != nil {
        server = network.servers[serverName]
    }
    reliable = network.reliable
    longreordering = network.longReordering
    return
}

func (network *Network) isServerDead(endName interface{}, serverName interface{}, server *Server) bool {
    network.mu.Lock()
    defer network.mu.Unlock()

    if !network.enabled[endName] || network.servers[serverName] != server {
        return true
    }
    return false
}

func (network *Network) processRequest(req requestMessage) {
    enabled, serverName, server, reliable, longreordering := network.readEndNameInfo(req.endName)

    if enabled && serverName != nil && server != nil {
        if !reliable {
            // short delay
            ms := (rand.Int() % 27)
            time.Sleep(time.Duration(ms) * time.Millisecond)
        }

        if !reliable && (rand.Int()%1000) < 100 {
            // drop the request, return as if timeout
            req.responseMessageChan <- responseMessage{false, nil}
            return
        }

        // execute the request (call the RPC handler).
        // in a separate thread so that we can periodically check
        // if the server has been killed and the RPC should get a
        // failure reply.
        responseMessageChan := make(chan responseMessage)
        go func() {
            r := server.dispatch(req)
            responseMessageChan <- r
        }()

        // wait for handler to return,
        // but stop waiting if DeleteServer() has been called,
        // and return an error.
        var reply responseMessage
        replyOK := false
        serverDead := false
        for !replyOK && !serverDead {
            select {
            case reply = <-responseMessageChan:
                replyOK = true
            case <-time.After(100 * time.Millisecond):
                serverDead = network.isServerDead(req.endName, serverName, server)
                if serverDead {
                    go func() {
                        <-responseMessageChan // drain channel to let the goroutine created earlier terminate
                    }()
                }
            }
        }

        // do not reply if DeleteServer() has been called, i.e.
        // the server has been killed. this is needed to avoid
        // situation in which a client gets a positive reply
        // to an Append, but the server persisted the update
        // into the old Persister. config.go is careful to call
        // DeleteServer() before superseding the Persister.
        serverDead = network.isServerDead(req.endName, serverName, server)

        if !replyOK || serverDead {
            // server was killed while we were waiting; return error.
            req.responseMessageChan <- responseMessage{false, nil}
        } else if !reliable && (rand.Int()%1000) < 100 {
            // drop the reply, return as if timeout
            req.responseMessageChan <- responseMessage{false, nil}
        } else if longreordering && rand.Intn(900) < 600 {
            // delay the response for a while
            ms := 200 + rand.Intn(1+rand.Intn(2000))
            // Russ points out that this timer arrangement will decrease
            // the number of goroutines, so that the race
            // detector is less likely to get upset.
            time.AfterFunc(time.Duration(ms)*time.Millisecond, func() {
                atomic.AddInt64(&network.bytes, int64(len(reply.reply)))
                req.responseMessageChan <- reply
            })
        } else {
            atomic.AddInt64(&network.bytes, int64(len(reply.reply)))
            req.responseMessageChan <- reply
        }
    } else {
        // simulate no reply and eventual timeout.
        ms := 0
        if network.longDelays {
            // let Raft tests check that leader doesn't send
            // RPCs synchronously.
            ms = (rand.Int() % 7000)
        } else {
            // many kv tests require the client to try each
            // server in fairly rapid succession.
            ms = (rand.Int() % 100)
        }
        time.AfterFunc(time.Duration(ms)*time.Millisecond, func() {
            req.responseMessageChan <- responseMessage{false, nil}
        })
    }
}

// create a client end-point.
// start the thread that listens and delivers.
func (network *Network) MakeEnd(endName interface{}) *ClientEnd {
    network.mu.Lock()
    defer network.mu.Unlock()

    if _, ok := network.ends[endName]; ok {
        log.Fatalf("MakeEnd: %v already exists\n", endName)
    }

    clientEnd := &ClientEnd{}
    clientEnd.endName = endName
    clientEnd.requestMessageChan = network.requestMessageChan
    clientEnd.done = network.done
    network.ends[endName] = clientEnd
    network.enabled[endName] = false
    network.connections[endName] = nil

    return clientEnd
}

func (network *Network) AddServer(serverName interface{}, server *Server) {
    network.mu.Lock()
    defer network.mu.Unlock()

    network.servers[serverName] = server
}

func (network *Network) DeleteServer(serverName interface{}) {
    network.mu.Lock()
    defer network.mu.Unlock()

    network.servers[serverName] = nil
}

// connect a ClientEnd to a server.
// a ClientEnd can only be connected once in its lifetime.
func (network *Network) Connect(endName interface{}, serverName interface{}) {
    network.mu.Lock()
    defer network.mu.Unlock()

    network.connections[endName] = serverName
}

// enable/disable a ClientEnd.
func (network *Network) Enable(endName interface{}, enabled bool) {
    network.mu.Lock()
    defer network.mu.Unlock()

    network.enabled[endName] = enabled
}

// get a server's count of incoming RPCs.
func (network *Network) GetCount(serverName interface{}) int {
    network.mu.Lock()
    defer network.mu.Unlock()

    svr := network.servers[serverName]
    return svr.GetCount()
}

func (network *Network) GetTotalCount() int {
    x := atomic.LoadInt32(&network.count)
    return int(x)
}

func (network *Network) GetTotalBytes() int64 {
    x := atomic.LoadInt64(&network.bytes)
    return x
}

// a server is a collection of services, all sharing
// the same rpc dispatcher. so that e.g. both a Raft
// and a k/v server can listen to the same rpc endpoint.
type Server struct {
    mu       sync.Mutex
    services map[string]*Service
    count    int // incoming RPCs
}

func MakeServer() *Server {
    server := &Server{}
    server.services = map[string]*Service{}
    return server
}

func (server *Server) AddService(svc *Service) {
    server.mu.Lock()
    defer server.mu.Unlock()
    server.services[svc.name] = svc
}

func (server *Server) dispatch(req requestMessage) responseMessage {
    server.mu.Lock()

    server.count += 1

    // split Raft.AppendEntries into service and method
    dot := strings.LastIndex(req.serviceMethod, ".")
    serviceName := req.serviceMethod[:dot]
    methodName := req.serviceMethod[dot+1:]

    service, ok := server.services[serviceName]

    server.mu.Unlock()

    if ok {
        return service.dispatch(methodName, req)
    } else {
        choices := []string{}
        for k := range server.services {
            choices = append(choices, k)
        }
        log.Fatalf(
            "labrpc.Server.dispatch(): unknown service %v in %v.%v; expecting one of %v\n",
            serviceName, serviceName, methodName, choices,
        )
        return responseMessage{false, nil}
    }
}

func (server *Server) GetCount() int {
    server.mu.Lock()
    defer server.mu.Unlock()
    return server.count
}

// an object with methods that can be called via RPC.
// a single server may have more than one Service.
type Service struct {
    name     string
    receiver reflect.Value
    typ      reflect.Type
    methods  map[string]reflect.Method
}

//    type JunkServer struct {
//        mu     sync.Mutex
//        logStr []string
//        logInt []int
//    }
func MakeService(receiver interface{}) *Service {
    service := &Service{}
    service.typ = reflect.TypeOf(receiver) // labrpc.JunkServer
    service.receiver = reflect.ValueOf(receiver)
    service.name = reflect.Indirect(service.receiver).Type().Name() // JunkServer
    service.methods = map[string]reflect.Method{}

    for i := 0; i < service.typ.NumMethod(); i++ {
        method := service.typ.Method(i)
        methodType := method.Type
        methodName := method.Name

        if method.PkgPath != "" ||
            methodType.NumIn() != 3 ||
            methodType.In(2).Kind() != reflect.Ptr ||
            methodType.NumOut() != 0 {
            // the method is not suitable for a handler
            fmt.Printf("bad method: %v\n", methodName)
        } else {
            // the method looks like a handler
            service.methods[methodName] = method
        }
    }

    return service
}

func (service *Service) dispatch(methodName string, req requestMessage) responseMessage {
    if method, ok := service.methods[methodName]; ok {
        args := reflect.New(req.argsType)

        argsBuffer := bytes.NewBuffer(req.args)
        argsDecoder := labgob.NewDecoder(argsBuffer)
        argsDecoder.Decode(args.Interface())

        replyType := method.Type.In(2)
        replyType = replyType.Elem()
        reply := reflect.New(replyType)

        function := method.Func
        function.Call([]reflect.Value{service.receiver, args.Elem(), reply})

        replyBuffer := new(bytes.Buffer)
        replyEncoder := labgob.NewEncoder(replyBuffer)
        replyEncoder.EncodeValue(reply)

        return responseMessage{true, replyBuffer.Bytes()}
    } else {
        choices := []string{}
        for k := range service.methods {
            choices = append(choices, k)
        }
        log.Fatalf(
            "labrpc.Service.dispatch(): unknown method %v in %v; expecting one of %v\n",
            methodName, req.serviceMethod, choices,
        )
        return responseMessage{false, nil}
    }
}
