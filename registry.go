package clif

import (
	"reflect"
	"sort"
	"sync"
)

// ReduceMethod is type of the callback of the `Reduce` and `ReduceAsync` methods
type ReduceMethod func(name string, value interface{}) bool

// Registry is a small container holding objects which are injected into command calls.
type Registry struct {
	Container map[string]reflect.Value
}

// NewRegistry constructs new, empty registry
func NewRegistry() *Registry {
	return &Registry{
		Container: make(map[string]reflect.Value),
	}
}

// Alias registers an object under a different name. Eg for interfaces.
func (this *Registry) Alias(alias string, v interface{}) {
	r := reflect.ValueOf(v)
	this.Container[alias] = r
}

// Get returns registered object by their type name (`reflect.TypeOf(..).String()`)
func (this *Registry) Get(s string) reflect.Value {
	if v, ok := this.Container[s]; ok {
		return v
	} else {
		return reflect.ValueOf(nil)
	}
}

// Has checks whether a requested type is registered
func (this *Registry) Has(s string) bool {
	if _, ok := this.Container[s]; ok {
		return true
	} else {
		return false
	}
}

// Names returns sorted (asc) list of registered names
func (this *Registry) Names() []string {
	names := make([]string, len(this.Container))
	i := 0
	for n, _ := range this.Container {
		names[i] = n
		i++
	}
	sort.Strings(names)
	return names
}

// Reduce calls a bool-returning function on all registered objects in alphanumerical
// order and returns all selected as a slice
func (this *Registry) Reduce(cb ReduceMethod) []interface{} {
	res := make([]interface{}, 0)
	for _, name := range this.Names() {
		value := this.Container[name]
		if cb(name, value.Interface()) {
			res = append(res, value.Interface())
		}
	}
	return res
}

// ReduceAsync calls a bool-returning function on all registered objects
// concurrently and writes all selected objects in a return channel
func (this *Registry) ReduceAsync(cb ReduceMethod) chan interface{} {
	res := make(chan interface{})
	var wg sync.WaitGroup
	for name, value := range this.Container {
		wg.Add(1)
		go func(name string, value interface{}) {
			defer wg.Done()
			if cb(name, value) {
				res <- value
			}
		}(name, value.Interface())
	}
	go func() {
		wg.Wait()
		close(res)
	}()
	return res
}

// Register adds new object to registry. An existing object of same type would be replaced.
func (this *Registry) Register(v interface{}) {
	r := reflect.ValueOf(v)
	this.Container[r.Type().String()] = r
}


