//go:generate ../bin/events-bus-gen Events events.go

// nolint
package main

import productModel "github.com/oliosinter/go-events-bus/examples/models"

// Events is a an interface that will be used to generate events bus.
// Each method of this interface is considered to be an event.
// For each event two separate methods will be generated:
//  - a method to subscribe on event notifications
//  - a method to remove handler subscription from event notifications
type Events interface {
	UserRegistration(*UserInfo)
	ProductUpdate(*productModel.Product)
	CombinedEvent(UserInfo, productModel.Product)
}

// UserInfo is a struct for UserRegistration example event
type UserInfo struct {
	ID    uint64
	Name  string
	Email string
}
