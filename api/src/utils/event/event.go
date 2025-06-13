package event

import "github.com/olebedev/emitter"

// ------------------------------------------------------------
// : Enums
// ------------------------------------------------------------
const (
	UserConnected    = "user.connected"
	UserDisconnected = "user.disconnected"

	ExtractorItemStarted = "extractor.item.started"
	ExtractorItemFailed  = "extractor.item.failed"
	ExtractorItemDone    = "extractor.item.done"
)

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var (
	e = emitter.Emitter{Cap: 1}
)

// ------------------------------------------------------------
// : Methods
// ------------------------------------------------------------
func Use(topic string, middlewares ...func(*emitter.Event)) {
	e.Use(topic, middlewares...)
}

func Emit(topic string, args ...interface{}) chan struct{} {
	return e.Emit(topic, args...)
}

func On(topic string, middlewares ...func(*emitter.Event)) <-chan emitter.Event {
	return e.On(topic, middlewares...)
}

func Once(topic string, middlewares ...func(*emitter.Event)) <- chan emitter.Event {
	return e.Once(topic, middlewares...)
}

func Off(topic string, channels ...<-chan emitter.Event) {
	e.Off(topic, channels...)
}
