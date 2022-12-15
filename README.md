# mud

[![Go Reference](https://pkg.go.dev/badge/github.com/nvlled/mud.svg)](https://pkg.go.dev/github.com/nvlled/mud)

A go library for generic objects pools. Basically a thin wrapper
over `sync.Map` that allows using a single pool for different
type parameters, or even different types.

## Installation

```
go get github.com/nvlled/mud
```

## Usage && Example

```
pool := mud.NewPool()
obj1 := mud.Alloc(pool, func() *SomeType { new(SomeType) })
obj2 := mud.Alloc(pool, func() *OtherType { new(OtherType) })
obj3 := mud.Alloc(pool, func() *GenericType[int] { new(GenericType[int]) })
obj4 := mud.Alloc(pool, func() *GenericType[string] { new(GenericType[string]) })

// ...

mud.Free(pool, obj1)
mud.Free(pool, obj2)
mud.Free(pool, obj3)
mud.Free(pool, obj4)
```
