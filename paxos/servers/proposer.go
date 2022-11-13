package servers

import (
    "fmt"
    "paxos/message"
)

type Proposer struct {
    id int
    round int
    number int
    acceptors []int
}

func (proposer *Proposer) Propose(v interface{}) interface{} {
    proposer.round++
    proposer.number = proposer.proposalNumber()

    prepareCount := 0
    maxNumber := 0
    for _, acceptor_port := range proposer.acceptors {
        args := message.MsgArgs {
            Number: proposer.number,
            From: proposer.id,
            To: acceptor_port,
        }

        reply := new(message.MsgReply)
        ok := message.Call(fmt.Sprintf("127.0.0.1:%d", acceptor_port), "Acceptor.Prepare", args, reply)
        if !ok {
            continue
        }

        if reply.Ok {
            prepareCount++
            if reply.Number > maxNumber {
                maxNumber = reply.Number
                v = reply.Value
            }
        }

        if prepareCount == proposer.majority() {
            break
        }
    }

    acceptCount := 0
    if prepareCount >= proposer.majority() {
        for _, acceptor_port := range proposer.acceptors {
            args := message.MsgArgs {
                Number: proposer.number,
                Value: v,
                From: proposer.id,
                To: acceptor_port,
            }

            reply := new(message.MsgReply)
            ok := message.Call(fmt.Sprintf("127.0.0.1:%d", acceptor_port), "Acceptor.Accept", args, reply)
            if !ok {
                continue
            }

            if reply.Ok {
                acceptCount++
            }
        }
    }

    if acceptCount >= proposer.majority() {
        return v
    }

    return nil
}

func (proposer *Proposer) majority() int {
    return len(proposer.acceptors) / 2 + 1
}

func (proposer *Proposer) proposalNumber() int {
    return proposer.round << 16 | proposer.id
}

func NewProposer(id int, acceptorIds []int) *Proposer {
    return &Proposer {
        id: id,
        acceptors: acceptorIds,
    }
}