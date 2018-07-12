package notifications

import (
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

// EventName stores event name
type EventName string

// EventTypes stores information about event type for each event
type EventTypes map[EventName]interface{}

// Bus is an another events bus implementation
type Bus struct {
	events   map[EventName]reflect.Type
	handlers map[EventName][]reflect.Value
	lock     sync.RWMutex
}

// NewBus creates a new event bus.
// It takes a map of events with their data types and checks, that all event data types are structs or
// struct pointers.
//
// If any event has a wrong data type, error is returned.
func NewBus(types EventTypes) (*Bus, error) {
	bus := &Bus{
		events:   make(map[EventName]reflect.Type),
		handlers: make(map[EventName][]reflect.Value),
	}

	for name, eventType := range types {
		typeOf := reflect.TypeOf(eventType)
		if typeOf == nil {
			return nil, errors.Errorf("wrong event type: nil interface as a type for event %s", name)
		}
		if typeOf.Kind() != reflect.Struct && typeOf.Kind() != reflect.Ptr {
			return nil, errors.Errorf("wrong event type %s for event %s", typeOf.Kind(), name)
		}
		if typeOf.Kind() == reflect.Ptr && typeOf.Elem().Kind() != reflect.Struct {
			return nil, errors.Errorf("wrong event type %s for event %s", typeOf.Elem().Kind(), name)
		}

		bus.events[name] = typeOf
	}

	return bus, nil
}

// AddHandler registers new handler for an event. It checks that:
// - event exists;
// - handler is a func;
// - handler has one incoming argument and zero out coming;
// - handler has an expected interface for this event data type.
//
// If any of that requirements are wrong, error is returned.
func (b *Bus) AddHandler(name EventName, handler interface{}) error {
	typeOf, ok := b.getEventType(name)
	if !ok {
		return errors.Errorf("unknown event %s", name)
	}

	if err := checkHandlerForEventType(handler, typeOf); err != nil {
		return errors.Wrapf(err, "cannot add new handler for event %s", name)
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	if _, ok := b.handlers[name]; !ok {
		b.handlers[name] = make([]reflect.Value, 0, 1)
	}
	b.handlers[name] = append(b.handlers[name], reflect.ValueOf(handler))

	return nil
}

// TriggerEvent notifies all listeners.
// It makes sure, that:
// - event has been registered;
// - `data` has an expected type for this event.
//
// If all requirements are met, handlers will be notified, else error is returned.
func (b *Bus) TriggerEvent(name EventName, data interface{}) error {
	expectedEventType, ok := b.getEventType(name)
	if !ok {
		return errors.Errorf("unknown event %s", name)
	}

	if err := compareTypes(expectedEventType, reflect.TypeOf(data)); err != nil {
		return errors.Wrapf(err, "wrong event type for %s event", name)
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	if handlers, ok := b.handlers[name]; ok {
		for _, handler := range handlers {
			handler.Call([]reflect.Value{reflect.ValueOf(data)})
		}
	}

	return nil
}

func (b *Bus) getEventType(name EventName) (reflect.Type, bool) {
	b.lock.RLock()
	desc, ok := b.events[name]
	b.lock.RUnlock()

	return desc, ok
}

func checkHandlerForEventType(handler interface{}, expectedEventType reflect.Type) error {
	handlerRef := reflect.TypeOf(handler)
	if handlerRef.Kind() != reflect.Func {
		return errors.Errorf("wrong handler type %s", handlerRef.Kind())
	}

	if handlerRef.NumIn() != 1 {
		return errors.Errorf("wrong handler in coming arguments count: %d instead of 1", handlerRef.NumIn())
	}

	if handlerRef.NumOut() != 0 {
		return errors.Errorf("wrong handler out coming arguments count: %d instead of 0", handlerRef.NumOut())
	}

	return compareTypes(expectedEventType, handlerRef.In(0))
}

func compareTypes(expected, actual reflect.Type) error {
	if expected.Kind() != actual.Kind() || expected != actual {
		return errors.Errorf("actual event type %s is not equal to expected event type %s", actual, expected)
	}

	return nil
}
