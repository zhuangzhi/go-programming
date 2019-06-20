package cache
import (
	"net"
	"context"
)
const (
	PartitionCount = 10000
)

type Cache struct {
	cache []map[string][]byte
}

func New() (*Cache, error) {
	return &Cache{
		cache = make([]map[string][]byte, PartitionCount)
	}, nil
}



func (c Cache) Get(ctx context.Context) chan net.Conn {
	ch := make(chan net.Conn)

	
	go func() {
		conn, err := net.Dial("xxx")
		if err!=nil {
			//xxx
		}

		select {
			ch <-conn:
		<-time.After(time.Second)
			conn.Close()
		}
	}()

	go func() {
		select{
			<-ctx.Done():
			close(ch)
		}
	}()

	return ch
}

func t() {
	c := Cache{}
	ctx, cancel:=context.WithCancel(context.Background())
	ctx, cancel = context.WithTimeout(ctx, time.Second)
	ch := c.Get() 
	conn, ok := <-ch
	if ok
}