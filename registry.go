package cli

import "reflect"

type Registry struct {
	reg map[string]reflect.Value
}

func NewRegistry() *Registry {
	return &Registry{
		reg: make(map[string]reflect.Value),
	}
}

func (this *Registry) Get(s string) reflect.Value {
	if v, ok := this.reg[s]; ok {
		return v
	} else {
		return reflect.ValueOf(nil)
	}
}

func (this *Registry) Has(s string) bool {
	if _, ok := this.reg[s]; ok {
		return true
	} else {
		return false
	}
}

func (this *Registry) Register(v interface{}) {
	r := reflect.ValueOf(v)
	this.reg[r.Type().String()] = r
}

func (this *Registry) Alias(alias string, v interface{}) {
	r := reflect.ValueOf(v)
	this.reg[alias] = r
}
