// This packages is thin wrapper on sync.Pool
// with the purpose of adding type safety using generics.
// Any type of object can be allocated and free using
// a unified object pool.
package mud

import (
	"reflect"
	"sync"
)

// A Constructor is a function that creates
// an object and returns a pointer to it.
// Example:
//
//	type Banana struct{}
//	var bananaCtor Constructor[Banana] = func() *Banana { return &Banana{} }
type Constructor[T any] func() *T

// A Pool is a generic pool that can store any type.
type Pool struct {
	billiard *sync.Map
	//billiard map[reflect.Type]*sync.Pool
}

// Creates a new pool.
func NewPool() *Pool {
	return &Pool{
		billiard: &sync.Map{},
	}
}

func getPool(mudPool *Pool, t reflect.Type) (*sync.Pool, bool) {
	obj, ok := mudPool.billiard.Load(t)
	if !ok {
		return nil, false
	}
	return obj.(*sync.Pool), ok
}

func getOrCreatePool[T any](mudPool *Pool, t reflect.Type, ctor Constructor[T]) *sync.Pool {
	obj, ok := mudPool.billiard.Load(t)

	if !ok {
		obj = &sync.Pool{
			New: func() any {
				if ctor != nil {
					return ctor()
				}
				var x T
				return &x
			},
		}
		mudPool.billiard.Store(t, obj)
	}
	return obj.(*sync.Pool)
}

// Pre-allocates a pool of the particular type with the given size.
// Note: this will increase the size, not set size to numObjects.
// Successive N calls with numObjects
// is equivalent to setting the size to (at least) N*numObjects.
func PreAlloc[T any](mudPool *Pool, ctor Constructor[T], numObjects int) {
	typeOf := reflect.TypeOf(ctor())
	pool := getOrCreatePool(mudPool, typeOf, ctor)
	for i := 0; i < numObjects; i++ {
		x := ctor()
		pool.Put(x)
	}
}

// Returns an object of type T, either an existing one from the pool,
// or a new one using the constructor. Alloc() will never nil
// as long as ctor doesn't return a nil.
func Alloc[T any](mudPool *Pool, ctor Constructor[T]) *T {
	var ptr *T
	typeOf := reflect.TypeOf(ptr)

	pool := getOrCreatePool(mudPool, typeOf, ctor)
	object := pool.Get().(*T)

	return object
}

// Similar to Alloc(), but will return nil if there
// are no at least one previous call to Free() or PreAlloc() objects of the given type.
func Get[T any](mudPool *Pool) *T {
	var ptr *T
	typeOf := reflect.TypeOf(ptr)

	pool, ok := getPool(mudPool, typeOf)
	if !ok {
		var nilObject *T
		return nilObject
	}
	return pool.Get().(*T)
}

// Frees the object to the pool and makes it available for
// later Alloc()'s. If the object wasn't
// allocated from the pool, it will still be added.
func Free[T any](mudPool *Pool, object *T) {
	typeOf := reflect.TypeOf(object)
	pool, ok := getPool(mudPool, typeOf)
	if ok {
		pool.Put(object)
	}
}
