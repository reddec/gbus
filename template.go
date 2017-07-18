package main

var text = `
package {{.Package}}

import (
    "sync"
)


// How to use it:
//
// type Sample struct {
//      emitter {{.Name}}EventEmitter
// }
//
// func (s *Sample) Events() {{.Name}}Events { return &s.emitter }
//
// ...
//
// func (s *Sample) SomeJob() {
//    {{- range $event := .Events}}
// }
//


type (
{{- range $event := .Events}}
	 {{$.Name}}{{$event.Name}}HandlerFunc func({{$event | signature}}) // Listener handler function for event '{{$event.Name}}'
{{- end}}
)

// {{.Name}}EventEmitter implements events listener and events emitter operations
// for events {{.EventsList}}
type {{.Name}}EventEmitter struct {
	{{.Name}}Events // Implements listener operations
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

	// AddHandler adds handler for events ({{.EventsList}})
    AddHandler(handler {{.Name}})

    // RemoveHandler remove handler for events
    RemoveHandler(handler {{.Name}})
}

{{range $event := .Events}}
// On{{$event.Name}} adds event listener for event '{{$event.Name}}'
func (bus *{{$.Name}}EventEmitter) On{{$event.Name}}(handler {{$.Name}}{{$event.Name}}HandlerFunc) {
    bus.lock{{$event.Name}}.Lock()
    defer bus.lock{{$event.Name}}.Unlock()
    bus.on{{$event.Name}} = append(bus.on{{$event.Name}}, handler)
}

// No{{$event.Name}} excludes event listener
func (bus *{{$.Name}}EventEmitter)  No{{$event.Name}}(handler {{$.Name}}{{$event.Name}}HandlerFunc) {
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
func (bus *{{$.Name}}EventEmitter)  {{$event.Name}}({{$event | signature}}) {
    bus.lock{{$event.Name}}.RLock()
    defer bus.lock{{$event.Name}}.RUnlock()
    for _, f := range bus.on{{$event.Name}} {
        f({{$event | call}})
    }
}
{{end}}

// AddHandler adds handler for events ({{.EventsList}})
func (bus *{{$.Name}}EventEmitter) AddHandler(handler {{$.Name}}) {
	{{- range $event := .Events}}
	bus.On{{$event.Name}}(handler.{{$event.Name}})
	{{- end}}
}

// RemoveHandler remove handler for events
func (bus *{{$.Name}}EventEmitter) RemoveHandler(handler {{$.Name}}) {
	{{- range $event := .Events}}
	bus.No{{$event.Name}}(handler.{{$event.Name}})
	{{- end}}
}
`
