package notifications

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testBusEventType struct {
	counter int
}

func TestBus_WrongEventTypes_ExpectError(t *testing.T) {
	var nilInterface interface{}
	checks := []interface{}{
		[]int{},
		123,
		"string",
		nilInterface,
	}

	for _, eventType := range checks {
		types := EventTypes{"name": eventType}

		_, err := NewBus(types)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "wrong event type")
	}
}

func TestBus_AddHandler_UnknownEventType_ExpectError(t *testing.T) {
	bus, _ := NewBus(EventTypes{ // nolint: gas
		"event": (*testBusEventType)(nil),
	})
	var event EventName = "some event"

	err := bus.AddHandler(event, func() {})

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unknown event")
}

func TestBus_AddHandler_WrongType_ExpectError(t *testing.T) {
	var event EventName = "event"
	bus, _ := NewBus(EventTypes{ // nolint: gas
		event: (*testBusEventType)(nil),
	})

	err := bus.AddHandler(event, &testBusEventType{})

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "wrong handler type")
}

func TestBus_AddHandler_WrongHandlerInArguments_ExpectError(t *testing.T) {
	var event EventName = "event"
	bus, _ := NewBus(EventTypes{ // nolint: gas
		event: (*testBusEventType)(nil),
	})

	err := bus.AddHandler(event, func() {})

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "in coming arguments count")
}

func TestBus_AddHandler_WrongHandlerOutArguments_ExpectError(t *testing.T) {
	var event EventName = "event"
	bus, _ := NewBus(EventTypes{ // nolint: gas
		event: (*testBusEventType)(nil),
	})

	err := bus.AddHandler(event, func(e string) error { return nil })

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "out coming arguments count")
}

func TestBus_AddHandler_WrongHandlerEventType_ExpectError1(t *testing.T) {
	var event EventName = "event"
	bus, _ := NewBus(EventTypes{ // nolint: gas
		event: (*testBusEventType)(nil),
	})

	err := bus.AddHandler(event, func(e string) {})

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "is not equal to expected event type")
}

func TestBus_AddHandler_WrongHandlerEventType_ExpectError2(t *testing.T) {
	var event EventName = "event"
	bus, _ := NewBus(EventTypes{ // nolint: gas
		event: (*testBusEventType)(nil),
	})

	err := bus.AddHandler(event, func(e testBusEventType) {})

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "is not equal to expected event type")
}

func TestBus_AddHandler_WrongHandlerEventType_ExpectError3(t *testing.T) {
	var event EventName = "event"
	bus, _ := NewBus(EventTypes{ // nolint: gas
		event: testBusEventType{},
	})

	err := bus.AddHandler(event, func(e *testBusEventType) {})

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "is not equal to expected event type")
}

func TestBus_AddHandler_Trigger_ExpectOk1(t *testing.T) {
	var event EventName = "event"
	bus, _ := NewBus(EventTypes{ // nolint: gas
		event: (*testBusEventType)(nil),
	})

	err1 := bus.AddHandler(event, func(e *testBusEventType) {})
	err2 := bus.TriggerEvent(event, &testBusEventType{})

	assert.Nil(t, err1)
	assert.Nil(t, err2)
}

func TestBus_AddHandler_Trigger_ExpectOk2(t *testing.T) {
	var event EventName = "event"
	bus, _ := NewBus(EventTypes{ // nolint: gas
		event: testBusEventType{},
	})

	err1 := bus.AddHandler(event, func(e testBusEventType) {})
	err2 := bus.TriggerEvent(event, testBusEventType{})

	assert.Nil(t, err1)
	assert.Nil(t, err2)
}

func TestBus_Trigger_ExpectError1(t *testing.T) {
	var event EventName = "event"
	bus, _ := NewBus(EventTypes{}) // nolint: gas

	err := bus.TriggerEvent(event, testBusEventType{})

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unknown event")
}

func TestBus_Trigger_ExpectError2(t *testing.T) {
	var event EventName = "event"
	bus, _ := NewBus(EventTypes{ // nolint: gas
		event: (*testBusEventType)(nil),
	})

	err := bus.TriggerEvent(event, testBusEventType{})

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "wrong event type")
}

func TestBus_Trigger_ExpectOk(t *testing.T) {
	var event EventName = "event"
	bus, _ := NewBus(EventTypes{ // nolint: gas
		event: (*testBusEventType)(nil),
	})

	handler := func(e *testBusEventType) {
		e.counter++
	}

	expected := 5
	for i := 0; i < expected; i++ {
		err := bus.AddHandler(event, handler)
		assert.Nil(t, err)
	}

	eventData := &testBusEventType{}
	err := bus.TriggerEvent(event, eventData)

	assert.Nil(t, err)
	assert.Equal(t, expected, eventData.counter)
}
