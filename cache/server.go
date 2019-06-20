package cache

import (
	"github.com/golang/groupcache"
)

type Server struct {
}



func test() {
	group := groupcache.NewGroup("hello", 1000, nil)
	group.Get()
}
