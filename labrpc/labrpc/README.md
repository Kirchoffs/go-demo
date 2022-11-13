# Notes
From __MIT 6.824__

## Code Analysis
### Basic Flow
When we call `clientEnd.Call`, it prepares a request with `responseMessageChan` and send the request to `clientEnd.requestMessageChan`. `clientEnd.requestMessageChan` is a copy of the network's channel. So, in the goroutine `network.processRequest`, the network extracts and checks the request (from the network's `requestMessageChan`)'s destination and dispatches it to the corresponding server. Also, the network creates another `responseMessage` channel for the result of server's processing.

After the request is processed, the network sends the reply back to the client through `req.responseMessageChannel`, and the client receives the reply from `responseMessageChannel`.

- client --> requestMessage channel --> network --> server.service.method
- server.service.method --> responseMessage channel --> network --> responseMessage channel --> client

### Goroutine with Sleep
```go
time.AfterFunc(time.Duration(ms)*time.Millisecond, func() {
    atomic.AddInt64(&rn.bytes, int64(len(reply.reply)))
    req.replyCh <- reply
})
```

```go
go func() {
    time.Sleep(time.Duration(ms) * time.Millisecond)
    atomic.AddInt64(&rn.bytes, int64(len(reply.reply)))
    req.replyCh <- reply
}()
```

```go
time.Sleep(time.Duration(ms) * time.Millisecond)
atomic.AddInt64(&rn.bytes, int64(len(reply.reply)))
req.replyCh <- reply
```

`time.AfterFunc` way is the most efficient and preferred approach. `time.AfterFunc` schedules the function to run after the specified duration in a separate goroutine managed by the time package. This avoids blocking the current goroutine.
