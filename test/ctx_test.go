package test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func Dial(url string) (net.Conn, error) {
	return net.Dial("tcp", url)
}

type dialResult struct {
	Conn net.Conn
	Err  error
}

func (d dialResult) Close() error {
	if d.Conn != nil {
		return d.Conn.Close()
	}
	return nil
}

func DialCtx(ctx context.Context, url string, done chan struct{}) (net.Conn, error) {
	ch := make(chan dialResult)
	go func() {
		conn, err := Dial(url)
		defer func() {
			if recover() != nil {
				if conn != nil {
					fmt.Println("timeout!!!")
					conn.Close()
				}
			}
			if done != nil {
				close(done)
			}
		}()
		ch <- dialResult{Conn: conn, Err: err}
	}()

	select {
	case <-ctx.Done():
		close(ch)
		return nil, ErrTimeout
	case r := <-ch:
		return r.Conn, r.Err
	}
}

var (
	ErrTimeout = errors.New("error timeout")
)

func RunWithCtx(ctx context.Context, fn func() io.Closer) (interface{}, error) {
	ch := make(chan interface{})
	go func() {
		r := fn()
		defer func() {
			if recover() != nil {
				if r != nil {
					r.Close()
				}
			}
		}()
		ch <- r
	}()
	select {
	case <-ctx.Done():
		close(ch)
		return nil, ErrTimeout
	case r := <-ch:
		return r, nil
	}
}

// DialCtx2
func DialCtx2(ctx context.Context, url string) (net.Conn, error) {
	r, err := RunWithCtx(ctx, func() io.Closer {
		conn, err := Dial(url)
		return dialResult{conn, err}
	})
	if err != nil {
		return nil, err
	}
	dr := r.(dialResult)
	return dr.Conn, dr.Err
}

type server struct {
	l net.Listener
}

func (s *server) Start(port int, sleep time.Duration) error {
	l, err := net.Listen("tcp", fmt.Sprint("127.0.0.1:", port))
	if err != nil {
		return err
	}
	s.l = l
	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}
		c.Close()
		time.Sleep(sleep)
	}
}

func (s server) Close() {
	l.Close()
}

func TestDial(t *testing.T) {
	s := server{}
	sina := "www.sina.com:80"
	start := time.Now()
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	co, err := DialCtx(ctx, sina, nil)
	fmt.Printf("elapsed: %v\n", time.Since(start))
	if err != nil {
		fmt.Println(err)
	}
	if co != nil {
		co.Close()
	}
	// <-end

	DialTimout(sina, time.Second/10)
	DialTimout(sina, 10*time.Second)
}

func DialTimout(url string, timeout time.Duration) {
	start := time.Now()
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	co, err := DialCtx2(ctx, "www.sina.com:80")
	fmt.Printf("elapsed: %v\n", time.Since(start))
	if err != nil {
		fmt.Println(err)
	}
	if co != nil {
		co.Close()
	}
}
