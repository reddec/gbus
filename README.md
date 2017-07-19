# GBUS - golang event bus generator

Generates an event emitter based on specified interface.

All functions without return value in the interface will be interpreted as event. For example (sample.go):

```go

type Rocket interface {
    Launched()
    LaunchFailed(reason error)
    Landed(lat, lon float64)
}

// go:generate Rocket sample.go

```


will produce file sample_eventbus.go with subscribe interface (RocketEvents) and implementation of event emitter (RocketEventEmitter).

Usually used as

```go

type Something struct {
    //...
    events  RocketEventEmitter
}

//...

func (s *Something) Events() RocketEvents { return &s.events }

```


