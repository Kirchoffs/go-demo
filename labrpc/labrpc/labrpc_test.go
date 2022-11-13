package labrpc

import (
    "runtime"
    "strconv"
    "sync"
    "testing"
)

type JunkArgs struct {
    X int
}

type JunksReply struct {
    X string
}

type JunkServer struct {
    mu     sync.Mutex
    logStr []string
    logInt []int
}

func (junkServer *JunkServer) HandlerStringToInt(args string, reply *int) {
    junkServer.mu.Lock()
    defer junkServer.mu.Unlock()

    junkServer.logStr = append(junkServer.logStr, args)
    *reply, _ = strconv.Atoi(args)
}

func (junkServer *JunkServer) HandlerIntToString(args int, reply *string) {
    junkServer.mu.Lock()
    defer junkServer.mu.Unlock()

    junkServer.logInt = append(junkServer.logInt, args)
    *reply = strconv.Itoa(args)
}

func (junkServer *JunkServer) HandlerWithPointer(args *JunkArgs, reply *JunksReply) {
    reply.X = "pointer"
}

func (junkServer *JunkServer) HandlerWithoutPointer(args JunkArgs, reply *JunksReply) {
    reply.X = "no pointer"
}

func TestBasic(t *testing.T) {
    runtime.GOMAXPROCS(4)

    network := MakeNetwork()
    defer network.Cleanup()

    clientEnd := network.MakeEnd("end-42")

    junkServer := &JunkServer{}
    service := MakeService(junkServer)

    server := MakeServer()
    server.AddService(service)
    network.AddServer("server-42", server)

    network.Connect("end-42", "server-42")
    network.Enable("end-42", true)

    {
        var reply string
        clientEnd.Call("JunkServer.HandlerIntToString", 42, &reply)
        if reply != "42" {
            t.Fatalf("expected reply to be 42, got %s", reply)
        }
    }

    {
        var reply int
        clientEnd.Call("JunkServer.HandlerStringToInt", "42", &reply)
        if reply != 42 {
            t.Fatalf("expected reply to be 42, got %d", reply)
        }

    }
}

func TestTypes(t *testing.T) {
    runtime.GOMAXPROCS(4)

    network := MakeNetwork()
    defer network.Cleanup()

    clientEnd := network.MakeEnd("end-42")

    junkServer := &JunkServer{}
    service := MakeService(junkServer)

    server := MakeServer()
    server.AddService(service)
    network.AddServer("server-42", server)

    network.Connect("end-42", "server-42")
    network.Enable("end-42", true)

    {
        var args JunkArgs
        var reply JunksReply

        clientEnd.Call("JunkServer.HandlerWithPointer", &args, &reply)
        if reply.X != "pointer" {
            t.Fatalf("expected reply to be pointer, got %s", reply.X)
        }
    }

    {
        var args JunkArgs
        var reply JunksReply

        clientEnd.Call("JunkServer.HandlerWithoutPointer", args, &reply)
        if reply.X != "no pointer" {
            t.Fatalf("expected reply to be no pointer, got %s", reply.X)
        }
    }
}
