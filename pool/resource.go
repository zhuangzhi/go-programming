package pool

import (
	"sync/atomic"
)

const (
	// Ready Resource ready to use.
	Ready int32 = iota

	// Locked resource locked by user.
	Locked

	// LockerRelease locker should release it when 'Unlock'
	LockerRelease

	// Released resource released.
	Released
)

// Resource a resource can be lock, released.
type Resource struct {
	Value       interface{}
	state       int32
	releaseFunc func()
}

// Lock the resource to use, if lock failed, don't used it. Invoke 'Unlock' after Lock successful.
func (r *Resource) Lock() bool {
	return atomic.CompareAndSwapInt32(&r.state, Ready, Locked)
}

// Unlock the resource, released the resource is LockerRelease.
// Caller should release any resource it generated when Unlock failed because the caller has give up this request.
func (r *Resource) Unlock() bool {
	for {
		if atomic.LoadInt32(&r.state) == Released {
			return false
		}
		if atomic.CompareAndSwapInt32(&r.state, Locked, Ready) {
			return true
		}
		if atomic.CompareAndSwapInt32(&r.state, LockerRelease, Released) {
			return false
		}
	}
}

// Release the resource. Set state to LockerRelease if some other go-routine locked the resource.
func (r *Resource) Release() bool {
	for {
		if atomic.LoadInt32(&r.state) == Released {
			return false
		}

		if atomic.CompareAndSwapInt32(&r.state, Ready, Released) {
			if r.releaseFunc != nil {
				r.releaseFunc()
			}
			return true
		}

		if atomic.CompareAndSwapInt32(&r.state, Locked, LockerRelease) {
			return false
		}
	}
}
