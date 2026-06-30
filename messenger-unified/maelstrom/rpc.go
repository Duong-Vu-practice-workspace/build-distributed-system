package maelstrom

import (
	"fmt"
	"sync"
	"time"
)

// pendingRequest tracks an in-flight synchronous RPC.
type pendingRequest struct {
	ch    chan Message
	timer *time.Timer
}

// rpcState holds pending synchronous RPC requests keyed by the outgoing
// message id.
type rpcState struct {
	mu       sync.Mutex
	requests map[int]*pendingRequest
	timeout  time.Duration
}

func newRPCState(timeout time.Duration) *rpcState {
	if timeout <= 0 {
		timeout = 1 * time.Second
	}
	return &rpcState{
		requests: make(map[int]*pendingRequest),
		timeout:  timeout,
	}
}

// SyncRPC sends body to dest and blocks until a matching reply arrives or the
// configured timeout elapses. The returned Message is the reply body wrapped
// in a Message envelope.
func (n *Node) SyncRPC(dest string, body map[string]interface{}) (Message, error) {
	return n.SyncRPCTimeout(dest, body, 0)
}

// SyncRPCTimeout sends body to dest and blocks until a matching reply arrives
// or timeout elapses. A zero timeout uses the default (1 second).
func (n *Node) SyncRPCTimeout(dest string, body map[string]interface{}, timeout time.Duration) (Message, error) {
	if timeout <= 0 {
		timeout = 1 * time.Second
	}

	msgID := n.nextMsgID()

	if n.rpc == nil {
		n.rpcMu.Lock()
		if n.rpc == nil {
			n.rpc = newRPCState(timeout)
		}
		n.rpcMu.Unlock()
	}

	ch := make(chan Message, 1)
	timer := time.AfterFunc(timeout, func() {
		n.rpc.mu.Lock()
		if pending, ok := n.rpc.requests[msgID]; ok {
			close(pending.ch)
			delete(n.rpc.requests, msgID)
		}
		n.rpc.mu.Unlock()
	})

	n.rpc.mu.Lock()
	n.rpc.requests[msgID] = &pendingRequest{ch: ch, timer: timer}
	n.rpc.mu.Unlock()

	n.sendWithID(dest, body, msgID)

	resp, ok := <-ch
	if !ok {
		return Message{}, fmt.Errorf("sync_rpc timeout waiting for reply to msg_id %d", msgID)
	}
	return resp, nil
}

// resolveRPC delivers an incoming reply to a pending SyncRPC call. It returns
// true when the reply was consumed by an RPC waiter.
func (n *Node) resolveRPC(replyTo interface{}, msg Message) bool {
	replyToFloat, ok := replyTo.(float64)
	if !ok {
		return false
	}
	replyToID := int(replyToFloat)

	n.rpcMu.Lock()
	if n.rpc == nil {
		n.rpcMu.Unlock()
		return false
	}
	n.rpc.mu.Lock()
	pending, ok := n.rpc.requests[replyToID]
	if ok {
		delete(n.rpc.requests, replyToID)
	}
	n.rpc.mu.Unlock()
	n.rpcMu.Unlock()

	if !ok {
		return false
	}

	if pending.timer != nil {
		pending.timer.Stop()
	}
	pending.ch <- msg
	return true
}
