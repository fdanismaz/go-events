package event

import (
	"fmt"
	"github.com/lithammer/shortuuid"
	"sync"
)

type Type string

type handler struct {
	id  string
	run func(args ...interface{})
}

type event struct {
	eventType Type
	args      []interface{}
}

var registry map[Type]map[string]handler
var eventChannel chan event
var stopChannel chan bool
var once sync.Once

func init() {
	once.Do(func() {
		registry = make(map[Type]map[string]handler)
		eventChannel = make(chan event, 300)
		stopChannel = make(chan bool)
		go startListening()
	})
}

func Stop() {
	stopChannel <- true
}

func startListening() {
	for {
		select {
		case e := <-eventChannel:
			callHandlers(e)
		case <-stopChannel:
			return
		}
	}
}

func callHandlers(e event) {
	handlerMap, ok := registry[e.eventType]
	if !ok {
		fmt.Printf("No handler function found for the emitted event type %s", string(e.eventType))
	} else {
		for handlerId := range handlerMap {
			h := handlerMap[handlerId]
			go executeHandler(h, e.args...)
		}
	}
}

// A decorator function to prevent panics
func executeHandler(h handler, args ...interface{}) {
	decoratedHandler := func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Panic recovered for the handler %s", h.id)
			}
		}()

		h.run(args...)
	}
	decoratedHandler()
}

// Registers the given event handler function for the given event type and returns an id for the handler
// so that the client can unsubscribe from that event type using that handler id.
// It is the client's (who registers a handler function for an event type) responsibility to call
// the unsubscribe function
func Subscribe(eventType Type, handlerFunction func(args ...interface{})) string {
	handlerId := generateHandlerId()
	addToRegistry(eventType, handlerFunction, handlerId)
	return handlerId
}

func addToRegistry(eventType Type, handlerFunction func(args ...interface{}), handlerId string) {
	handlerMap, ok := registry[eventType]
	if !ok {
		handlerMap = make(map[string]handler)
		registry[eventType] = handlerMap
	}
	eventHandler := handler{
		id:  handlerId,
		run: handlerFunction,
	}
	handlerMap[handlerId] = eventHandler
}

func SubscribeMultiple(eventTypes []Type, handlerFunction func(args ...interface{})) string {
	var handlerId string
	for i := range eventTypes {
		addToRegistry(eventTypes[i], handlerFunction, handlerId)
	}
	return handlerId
}

func generateHandlerId() string {
	return shortuuid.New()
}

func Emit(eventType Type, args ...interface{}) {
	e := event{
		eventType: eventType,
		args:      args,
	}
	fmt.Printf("Emitting event %s", string(eventType))
	go writeToEventChannel(e)
}

func writeToEventChannel(e event) {
	eventChannel <- e
}

func Unsubscribe(eventType Type, handlerId string) {
	handlerMap, ok := registry[eventType]
	if !ok {
		fmt.Printf("No handler function found for the emitted event type %s\n", string(eventType))
	} else {
		delete(handlerMap, handlerId)
	}
}

func UnsubscribeMultiple(eventTypes []Type, handlerId string) {
	for i := range eventTypes {
		Unsubscribe(eventTypes[i], handlerId)
	}
}
