# Notes
## Code Analysis
### Analysis 1
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
