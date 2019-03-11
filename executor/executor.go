package executor

import "errors"

var (
	// ErrTaskRejected task was rejected because task queue is reach to max.
	ErrTaskRejected = errors.New("executor: task rejected")

	// ErrTaskCanceled task was canceled by timeout or by caller (with context.Context).
	ErrTaskCanceled = errors.New("executor: task canceled")
)

type Local interface {
	ID() int64
	Name() string
	Put(interface{}, interface{})
	Get(interface{}) interface{}
}

type local struct {
	id   int64
	name string
	ctx  map[interface{}]interface{}
}

func (l *local) ID() int64 {
	return l.id
}
func (l *local) Name() string {
	return l.name
}

func (l *local) Put(key, value interface{}) {
	l.ctx[key] = value
}

func (l *local) Get(key interface{}) interface{} {
	return l.ctx[key]
}

func NewLocal(id int64, name string) Local {
	l := new(local)
	l.id = id
	l.name = name
	l.ctx = make(map[interface{}]interface{}, 32)
	return l
}

type LocalRunner func(Local)

type LocalCaller func(Local) (interface{}, error)

type Handler func(interface{}, error)

type Executor interface {
	Submit(LocalRunner) error
	Invoke(LocalCaller) (interface{}, error)
	InvokeAsync(LocalCaller, callback Handler)
}
