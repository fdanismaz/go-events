# go-events
```
go get github.com/fdanismaz/go-events/event
```

## Initialization
The event listener starts as the `github.com/fdanismaz/go-events/event` is imported because
the package has an `init` function which starts the listener until you explicitly call 
`event.Stop()` function.

## Subscribing to Events
To subscribe to an event you can use the `Subscribe` or `SubscribeMultiple` functions. The
`Subscribe` function subscribe your callback function to a specific event, `SubscribeMultiple`
function, on the other hand, subscribe your callback function to multiple events, meaning that
whenever any of the subscribed event is emitted, your callback function will be executed.

```go
var UserCreated event.Type = "user-created"

handlerId := event.Subscribe(UserCreated, func(args ...interface{}) {
    id := args[0].(int)
    name := args[1].(string)
    
    fmt.Printf("A new user is created. ID: %d, Name: %s\n", id, name)
})
```

## Unsubscribing from Events
Both the `Subscribe` and the `SubscribeMultiple` creates returns a handler id for your 
callback function. When you want to unsubscribe from a specific event or multiple events
you can use `Unsubscribe` or `UnsubscribeMultiple`.

```go
event.Unsubscribe(UserCreated, handlerId)
```

## Emitting Events
When you emit an event, an event object including the event type, and the arguments is created.
That event object is written to the event channel which is catched by the event listener. The event 
listener executes each registered handler function for that event type in a separate go routine.

```go
event.Emit(UserCreated, userId, userName)
```