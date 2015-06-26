package clif

import "reflect"

// Registry is a small container holding objects which are injected into command calls.
type Registry struct {
	reg map[string]reflect.Value
}

// NewRegistry constructs new, empty registry
func NewRegistry() *Registry {
	return &Registry{
		reg: make(map[string]reflect.Value),
	}
}

// Get returns registered object by their type name (`reflect.TypeOf(..).String()`)
func (this *Registry) Get(s string) reflect.Value {
	if v, ok := this.reg[s]; ok {
		return v
	} else {
		return reflect.ValueOf(nil)
	}
}

// Has checks whether a requested type is registered
func (this *Registry) Has(s string) bool {
	if _, ok := this.reg[s]; ok {
		return true
	} else {
		return false
	}
}

// Register adds new object to registry. An existing object of same type would be replaced.
func (this *Registry) Register(v interface{}) {
	r := reflect.ValueOf(v)
	this.reg[r.Type().String()] = r
}

// Alias registers an object under a different name. Eg for interfaces.
func (this *Registry) Alias(alias string, v interface{}) {
	r := reflect.ValueOf(v)
	this.reg[alias] = r
}
