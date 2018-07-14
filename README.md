# Golang events bus generator

This package is a fork from `reddec/gbus`, with significant enhancements and bug fixes. It allows you to generate events bus
according to your events description.

## Why?
There are a lot of different events bus implementations. Most of them work with interfaces and force you to cast event
messages to desired structures. This can lead to mistakes and unexpected behaviour. 


The aim of this package is to help you to maintain explicit event subscribe/trigger interface and to avoid
problems that you could face while working with abstract interfaces.

## Installation
```bash
$ go get github.com/oliosinter/go-events-bus/cmd/events-bus-gen
```

## Usage
Suppose, you have a file `events.go`, containing following `Events` interface:
```go
package notifications

type Events interface {
	UserRegistration(*UserInfo)
}

type UserInfo struct {
	ID    uint64
	Name  string
	Email string
}
```

`Events` interface is a description of your events. Each method is considered to be a single event.
After executing:
```bash
$ events-bus-gen Events events.go
```
you will get a generated events bus implementation with trigger/subscribe methods for each event. In current case it will
have following interface:
```go
package notifications

type EventsBus interface {
	UserRegistration(*UserInfo)
	// UserRegistration adds event listener for event 'UserRegistration'
	OnUserRegistration(handler func(arg0 *UserInfo))
	// RemoveUserRegistration excludes event listener
	RemoveUserRegistration(handler func(arg0 *UserInfo))
}
```

## Examples
Check our `./examples` folder.