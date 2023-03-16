# Sync
thread safe generic data structures implemented with go routines

## Syncmap
a thread safe generic map implemented with go routines

### how to use

```go
import "github.com/emilesalem/sync/pkg/syncmap"

s := syncmap.NewSyncmap(ctx, map[string]string{
    `key1`: `value1`,
    `key2`: `value2`,
    `key3`: `value3`,
})

s.Set(`key4`, `value4`)


v := s.Get(`key4`)
// v == `value4`

s.Delete(`key4`)

v = s.Get(`key4`)
// v == ``
```

## Syncqueue

a thread safe generic queue implemented with go routines

### how to use

```go
import "github.com/emilesalem/sync/pkg/syncmap"

s := syncqueue.NewSyncqueue[string](ctx, syncqueue.Options{Capacity: capacity})

ok := syncQueue.Add(m)
if !ok {
    // queue capacity exceeded
}

msg := <-syncQueue.Read()
```