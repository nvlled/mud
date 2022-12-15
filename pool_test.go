package mud_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/nvlled/mud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Maybe[T any] struct {
	value T
}

func MaybeIntCtor() *Maybe[int] {
	return &Maybe[int]{value: 0}
}
func MaybeStrCtor() *Maybe[string] {
	return &Maybe[string]{value: ""}
}

func TestGenericPoolSingle(t *testing.T) {
	genPool := mud.NewPool()
	y := mud.Alloc(genPool, MaybeIntCtor)

	assert.Equal(t, y.value, 0)
	y.value = 999
	assert.Equal(t, y.value, 999)

	mud.Free(genPool, y)

	y = mud.Alloc(genPool, MaybeIntCtor)
	assert.Equal(t, y.value, 999)

	y2 := mud.Alloc(genPool, MaybeIntCtor)
	y2.value = 123
	assert.Equal(t, y2.value, 123)

	mud.Free(genPool, y)
	mud.Free(genPool, y2)

	z := mud.Alloc(genPool, MaybeIntCtor)
	assert.Equal(t, z.value, 999)

	z = mud.Alloc(genPool, MaybeIntCtor)
	assert.Equal(t, z.value, 123)

	z = mud.Alloc(genPool, MaybeIntCtor)
	assert.Equal(t, z.value, 0)
}

func TestGenericPoolMixed(t *testing.T) {
	genPool := mud.NewPool()
	x := mud.Alloc(genPool, MaybeIntCtor)
	str := mud.Alloc(genPool, MaybeStrCtor)
	assert.Equal(t, x.value, 0)
	assert.Equal(t, str.value, "")

	x.value = 1
	str.value = "owl"

	mud.Free(genPool, x)
	mud.Free(genPool, str)

	x = mud.Alloc(genPool, MaybeIntCtor)
	str = mud.Alloc(genPool, MaybeStrCtor)

	assert.Equal(t, x.value, 1)
	assert.Equal(t, str.value, "owl")
}

func TestGenericGet(t *testing.T) {
	var maybeNil *Maybe[int]
	genPool := mud.NewPool()

	x := mud.Get[Maybe[int]](genPool)
	require.Equal(t, maybeNil, x)

	mud.PreAlloc(genPool, MaybeIntCtor, 10)

	x = mud.Get[Maybe[int]](genPool)
	require.NotEqual(t, maybeNil, x)
	x.value = 100

	mud.Free(genPool, x)
	x = mud.Get[Maybe[int]](genPool)
	require.NotEqual(t, x, maybeNil)
	require.Equal(t, 100, x.value)

	z := mud.Get[Maybe[int]](genPool)
	require.NotEqual(t, maybeNil, x)
	require.NotEqual(t, x.value, z.value)
}

func inc[T any](m *sync.Map, x T) {
	typeOf := reflect.TypeOf(x)
	n, ok := m.Load(typeOf)
	if !ok {
		n = 0
	}
	m.Store(typeOf, n.(int)+1)
}
func get[T any](m *sync.Map, x T) int {
	typeOf := reflect.TypeOf(x)
	n, ok := m.Load(typeOf)
	if !ok {
		return 0
	}
	return n.(int)
}
func TestMap(t *testing.T) {

	m := new(sync.Map)
	inc(m, 100)
	inc(m, 123)
	inc(m, "foo")

	assert.Equal(t, 2, get(m, 456))
	assert.Equal(t, 1, get(m, "foo"))
	assert.Equal(t, 0, get(m, 1.2))
}

func BenchmarkGenPool(b *testing.B) {
	genPool := mud.NewPool()
	mud.PreAlloc(genPool, MaybeIntCtor, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := mud.Alloc(genPool, MaybeIntCtor)
		mud.Free(genPool, x)
	}
}
