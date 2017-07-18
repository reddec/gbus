package main

import (
	"os"
	"github.com/reddec/astools"
	"text/template"
	"strings"
)

func Signature(m atool.Method, f *atool.File) string {
	t := ""
	for i, arg := range m.In {
		if i != 0 {
			t += ", "
		}
		t += arg.Name + " " + f.Printer.ToString(arg.Type)
	}
	return t
}

func Call(m atool.Method) string {
	t := ""
	for i, arg := range m.In {
		if i != 0 {
			t += ", "
		}
		t += arg.Name
	}
	return t
}

type params struct {
	Package string
	Name    string
	Events  []atool.Method
}

func (p *params) EventsList() string {
	t := ""
	for _, event := range p.Events {
		if t != "" {
			t += ", "
		}
		t += event.Name
	}
	return t
}

func main() {

	interfaces := os.Args[1:]
	if len(interfaces) == 0 {
		return
	}

	interfaceNames := map[string]bool{}
	for _, iface := range interfaces[:len(interfaces)-1] {
		interfaceNames[iface] = true
	}
	file := interfaces[len(interfaces)-1]
	info, err := atool.Scan(file)

	if err != nil {
		panic(err)
	}

	var busTemplate = template.Must(template.New("").Funcs(template.FuncMap{
		"title": strings.Title,
		"signature": func(m atool.Method) string {
			return Signature(m, info)
		},
		"call": Call,
	}).Parse(text))

	for _, iface := range info.Interfaces {
		if interfaceNames[iface.Name] {
			params := params{
				Package: info.Package,
				Name:    iface.Name,
				Events:  []atool.Method{},
			}

			for _, method := range iface.Methods {
				if len(method.Out) == 0 {
					params.Events = append(params.Events, method)
				}
			}
			if len(params.Events) > 0 {
				err := busTemplate.Execute(os.Stdout, &params)
				if err != nil {
					panic(err)
				}
			}
		}
	}

}
