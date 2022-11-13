package tests

import (
    "testing"
    "paxos/servers"
)

func start(acceptorIds []int, learnerIds []int) ([]*servers.Acceptor, []*servers.Learner) {
    acceptors := make([]*servers.Acceptor, 0)
    for _, acceptorId := range acceptorIds {
        acceptor := servers.NewAcceptor(acceptorId, learnerIds)
        acceptors = append(acceptors, acceptor)
    }

    learners := make([]*servers.Learner, 0)
    for _, learnerId := range learnerIds {
        learner := servers.NewLearner(learnerId, acceptorIds)
        learners = append(learners, learner)
    }

    return acceptors, learners
}

func cleanup(acceptors []*servers.Acceptor, learners []*servers.Learner) {
    for _, acceptor := range acceptors {
        acceptor.Close()
    }

    for _, learner := range learners {
        learner.Close()
    }
}

func TestSingleProposer(t *testing.T) {
    acceptorIds := []int{1001, 1002, 1003}
    learnerIds := []int{2001}

    acceptors, learners := start(acceptorIds, learnerIds)
    defer cleanup(acceptors, learners)

    proposer := servers.NewProposer(1, acceptorIds)

    value := proposer.Propose("hello world")
    if value != "hello world" {
        t.Errorf("Expected value to be 'hello world', got '%s'", value)
    }

    learnValue := learners[0].Chosen()
    if learnValue != "hello world" {
        t.Errorf("Expected learn value to be 'hello world', got '%s'", learnValue)
    }
}

func TestTwoProposers(t *testing.T) {
    acceptorIds := []int{1001, 1002, 1003}
    learnerIds := []int{2001}

    acceptors, learners := start(acceptorIds, learnerIds)
    defer cleanup(acceptors, learners)

    proposer1 := servers.NewProposer(1, acceptorIds)

    proposer2 := servers.NewProposer(2, acceptorIds)

    value1 := proposer1.Propose("hello world")
    value2 := proposer2.Propose("hi world")

    if value1 != value2 {
        t.Errorf("Expected value1 to be equal to value2, got '%s' and '%s'", value1, value2)
    }

    learnValue := learners[0].Chosen()
    if learnValue != value1 {
        t.Errorf("Expected learn value to be '%s', got '%s'", value1, learnValue)
    }
}
