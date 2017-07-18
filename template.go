package main

var text = `
package {{.Package}}

import (
    "sync"
)

type (
{{- range $event := .Events}}
	 {{$.Name}}{{$event.Name}}HandlerFunc func({{$event | signature}}) // Listener handler function for event '{{$event.Name}}'
{{- end}}
)

type bus{{.Name}} struct {
    {{- range $event := .Events}}
    lock{{$event.Name}} sync.RWMutex
    on{{$event.Name}}   []{{$.Name}}{{$event.Name}}HandlerFunc
    {{- end}}
}

// {{.Name}}Events is a client side of event bus that allows subscribe to
// {{.EventsList}} events
type {{.Name}}Events interface {
    {{- range $event := .Events}}

    // {{$event.Name}} adds event listener for event '{{$event.Name}}'
    On{{$event.Name}}(handler {{$.Name}}{{$event.Name}}HandlerFunc)

    // No{{$event.Name}} excludes event listener
    No{{$event.Name}}(handler {{$.Name}}{{$event.Name}}HandlerFunc)
    {{- end}}
}

// New{{.Name}}Bus creates pair of emitter and listener manager
func New{{.Name}}Bus() ({{.Name}}Events, {{.Name}}) {
    v:= &bus{{.Name}} {}
    return &v, &v
}

{{range $event := .Events}}
func (bus *bus{{$.Name}}) On{{$event.Name}}(handler {{$.Name}}{{$event.Name}}HandlerFunc) {
    bus.lock{{$event.Name}}.Lock()
    defer bus.lock{{$event.Name}}.Unlock()
    bus.on{{$event.Name}} = append(bus.on{{$event.Name}}, handler)
}

func (bus *bus{{$.Name}})  No{{$event.Name}}(handler {{$.Name}}{{$event.Name}}HandlerFunc) {
    bus.lock{{$event.Name}}.Lock()
    defer bus.lock{{$event.Name}}.Unlock()
    var res []{{$.Name | title}}{{$event.Name}}HandlerFunc
    for _, f := range bus.on{{$event.Name}} {
        if f != handler {
            res = append(res, f)
        }
    }
    bus.on{{$event.Name}} = res
}

// {{$event.Name}} emits event with same name
func (bus *bus{{$.Name}})  {{$event.Name}}({{$event | signature}}) {
    bus.lock{{$event.Name}}.RLock()
    defer bus.lock{{$event.Name}}.RUnlock()
    for _, f := range bus.on{{$event.Name}} {
        f({{$event | call}})
    }
}
{{end}}`
