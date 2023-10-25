There is an external service that processes some abstract objects by batches. This service can process only a certain number of items n in a given time interval p. If the limit is exceeded, the service blocks further processing for a long time.
The task is to implement a client to this external service, which will allow it to process the maximum possible number of objects without blocking. It is not necessary to give the implementation of the external service!

````

Service definition:
package main
import (
	"context"
	"errors"
	"time"
)
// ErrBlocked reports if service is blocked.
var ErrBlocked = errors.New("blocked")
// Service defines external service that can process batches of items.
type Service interface {
	GetLimits() (n uint64, p time.Duration)
	Process(ctx context.Context, batch Batch) error
}
// Batch is a batch of items.
type Batch []Item

// Item is some abstract item.
type Item struct{}

````
So we make a go project with main.go
you can fist build the go project with 
```
go build
```
Then execute the Executable file with just ./GoCode or what ever executable file name
```
./GoCode
```

For the Serive Part we can implement a service where items can be process in DB or doing some magic by calling ```Process(ctx context.Context, batch Batch) error``` and set the limits and Get by the service method which called ```GetLimits()```
