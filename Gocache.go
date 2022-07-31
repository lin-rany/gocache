package GOcache

import (
	"fmt"
	"log"
	"sync"
)

//
type Getter interface {
	Get(key string) ([]byte, error)
}

type Getterfunc func(key string) ([]byte, error)

func (f Getterfunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	maincache cache
	getter    Getter
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cachebyte int64, getter Getter) *Group {
	if getter == nil {
		panic("getter is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		getter:    getter,
		maincache: cache{cachebytes: cachebyte},
		name:      name,
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

func (g *Group) Get(key string) (Byteview, error) {
	if key == "" {
		return Byteview{}, fmt.Errorf("key is required")
	}
	if v, ok := g.maincache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	return g.load(key)

}

func (g *Group) load(key string) (Byteview, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (Byteview, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return Byteview{}, err
	}
	value := Byteview{b: Clonebyte(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value Byteview) {
	g.maincache.add(key, value)
}
