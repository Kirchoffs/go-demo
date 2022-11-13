package message

import (
    "net/rpc"
)

type MsgArgs struct {
    Number int
    Value interface{}
    From int
    To int
}

type MsgReply struct {
    Number int
    Value interface{}
    Ok bool
}

func Call(srv string, name string, args interface{}, reply interface{}) bool {
    c, err := rpc.Dial("tcp", srv)
    if err != nil {
        return false
    }
    defer c.Close()

    err = c.Call(name, args, reply)
    return err == nil
}