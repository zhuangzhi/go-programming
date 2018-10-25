# Options

Golang support optional parameters in function. The parameters would support default value.

## Solution

A typical solution is use parameter build functions as the input parameter of fuction and build parameter values in the function.

```go

type Options struct {
    RetryAttemps int
    RetryInterval time.Duration
    Timeout time.Duration
    Operate string
    Url string
}

type Option func(*Options) error

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
func Request(opts ... Option) {
    option := &Options{}
    for _, opt := range opts {
        opt(option)
    }
    // .. doing works.
}

```

## Default Value

To support default value, we can define a default.

```go


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


```