package realtime

type Subscription struct {
	Events <-chan StreamEvent
	Close  func()
}
