package servers

import (
    "fmt"
    "log"
    "net"
    "net/rpc"
    "paxos/message"
)

type Learner struct {
    rpcListener net.Listener
    id int
    acceptedMsg map[int]message.MsgArgs
}

func (learner *Learner) Learn(args *message.MsgArgs, reply *message.MsgReply) error {
    acceptedMsg := learner.acceptedMsg[args.From]
    if acceptedMsg.Number < args.Number {
        learner.acceptedMsg[args.From] = *args
        reply.Ok = true
    } else {
        reply.Ok = false
    }
    return nil
}

func (learner *Learner) Chosen() interface{} {
    acceptCounts := make(map[int]int)
    acceptMsg := make(map[int]message.MsgArgs)

    for _, accepted := range learner.acceptedMsg {
        if accepted.Number != 0 {
            acceptCounts[accepted.Number]++
            acceptMsg[accepted.Number] = accepted
        }
    }

    for n, count := range acceptCounts {
        if count >= learner.majority() {
            return acceptMsg[n].Value
        }
    }

    return nil
}

func (learner *Learner) majority() int {
    return len(learner.acceptedMsg) / 2 + 1
}

func NewLearner(id int, acceptorIds []int) *Learner {
    learner := &Learner{
        id: id,
        acceptedMsg: make(map[int]message.MsgArgs),
    }

    for _, acceptorId := range acceptorIds {
        learner.acceptedMsg[acceptorId] = message.MsgArgs{
            Number: 0,
            Value: nil,
        }
    }

    learner.server(id)
    return learner
}

func (learner *Learner) server(id int) {
    rpcs := rpc.NewServer()
    rpcs.Register(learner)
    addr := fmt.Sprintf(":%d", id)
    l, e := net.Listen("tcp", addr)
    if e != nil {
        log.Fatal("listen error:", e)
    }
    learner.rpcListener = l
    go func() {
        for {
            conn, err := learner.rpcListener.Accept()
            if err != nil {
                continue
            } else {
                go rpcs.ServeConn(conn)
            }
        }
    }()
}

func (learner *Learner) Close() {
    learner.rpcListener.Close()
}