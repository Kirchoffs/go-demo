package servers

import (
    "fmt"
    "log"
    "net"
    "net/rpc"
    "paxos/message"
)

type Acceptor struct {
    rpcListener net.Listener
    id int
    receivedNumber int            // Max number of prepare request 
    acceptedNumber int            // Max number of accepted number
    acceptedValue interface{}
    learners []int
}

func (acceptor *Acceptor) Prepare(args *message.MsgArgs, reply *message.MsgReply) error {
    if args.Number > acceptor.receivedNumber {
        reply.Ok = true
        reply.Number = acceptor.acceptedNumber
        reply.Value = acceptor.acceptedValue
        acceptor.receivedNumber = args.Number
    } else {
        reply.Ok = false
    }
    return nil
}

func (acceptor *Acceptor) Accept(args *message.MsgArgs, reply *message.MsgReply) error {
    if args.Number >= acceptor.receivedNumber {
        reply.Ok = true
        acceptor.receivedNumber = args.Number
        acceptor.acceptedNumber = args.Number
        acceptor.acceptedValue = args.Value

        for _, learner_port := range acceptor.learners {
            go func(learner int) {
                addr := fmt.Sprintf("127.0.0.1:%d", learner)
                args.From = acceptor.id
                args.To = learner
                resp := new(message.MsgReply)
                ok := message.Call(addr, "Learner.Learn", args, resp) // args.Number is already set by proposer
                if !ok {
                    return
                }
            }(learner_port)
        }
    } else {
        reply.Ok = false
    }
    return nil
}

func NewAcceptor(id int, learners []int) *Acceptor {
    acceptor := &Acceptor{
        id: id,
        learners: learners,
    }

    acceptor.server()
    return acceptor
}

func (acceptor *Acceptor) server() {
    rpcs := rpc.NewServer()
    rpcs.Register(acceptor)
    addr := fmt.Sprintf(":%d", acceptor.id)
    l, e := net.Listen("tcp", addr)
    if e != nil {
        log.Fatal("listen error: ", e)
    }
    acceptor.rpcListener = l

    go func() {
        for {
            conn, err := acceptor.rpcListener.Accept()
            if err != nil {
                continue
            }
            go rpcs.ServeConn(conn)
        }
    }()
}

func (acceptor *Acceptor) Close() {
    acceptor.rpcListener.Close()
}
