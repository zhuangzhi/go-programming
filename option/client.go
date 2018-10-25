package http

import (
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	DefaultUrl           = "http://localhost:8080"
	DefaultRetryAttemps  = 2
	DefaultRetryInterval = time.Second
	DefaultTimeout       = 2 * time.Second
	DefaultOperate       = "GET"
)

type Client struct {
	Timeout time.Duration
}

// Options ...
type Options struct {
	RetryAttemps  int
	RetryInterval time.Duration
	Timeout       time.Duration
	Operate       string
	Url           string
}

// Option ... http do option
type Option func(*Options)

func WithTimeout(timeout time.Duration) Option {
	//
	return func(op *Options) {
		op.Timeout = timeout
	}
}

func WithRetry(attemps int, interval time.Duration) Option {
	return func(op *Options) {
		op.RetryAttemps = attemps
		op.RetryInterval = interval
	}
}

func WithOperate(operate string) Option {
	return func(op *Options) {
		op.Operate = operate
	}
}

func WithUrl(url string) Option {
	return func(op *Options) {
		op.Url = url
	}
}

var defaultOption = Options{
	RetryAttemps:  DefaultRetryAttemps,
	RetryInterval: DefaultRetryInterval,
	Timeout:       DefaultTimeout,
	Operate:       DefaultOperate,
	Url:           DefaultUrl,
}

// Do ... http request with retry.
func Do(opts ...Option) {
	option := &Options{}
	// copy default values firstly.
	*option = defaultOption
	for _, opt := range opts {
		opt(option)
	}

	type httpResult struct {
		StatusCode int
		Body       []byte
		Error      error
	}

	// .. doing works.
	switch option.Operate {
	case "GET":
		{
			Retry(func() (interface{}, error) {
				result := &httpResult{}
				result.StatusCode, result.Body, result.Error = fasthttp.GetTimeout(nil, option.Url, option.Timeout)
				return result, result.Error
			}, option.RetryAttemps, option.RetryInterval)
		}
	}
}

// Call ...
type Call func() (interface{}, error)

// Retry ... retry running a blocking task.
func Retry(call Call, attemps int, interval time.Duration) (interface{}, error) {
	// ...
	var err error
	r, err := call()
	if err == nil {
		return r, nil
	}

	for i := 0; i < attemps; i++ {
		time.Sleep(interval)
		r, err = call()
		if err == nil {
			return r, nil
		}
	}

	return nil, fmt.Errorf("work faild after retry:%v, with error:%v", attemps, err)
}

//
