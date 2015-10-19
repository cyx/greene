# greene

Connection recycler using [ConnState][1].

[1]: https://golang.org/pkg/net/http/#ConnState

## use case

Rotating clients regularly to avoid pinning
persistent clients to a single node. Especially
targetted for an `ELB => [EC2, EC2...]` scenario
where hyper active clients get pinned to specific
EC2 instances.

## example

```go
// Recycle connections after 5 minutes.
server := &http.Server{
	Addr:      ":8000",
	ConnState: greene.New(time.Second * 300),
}
server.ListenAndServe()
```

## license

MIT
