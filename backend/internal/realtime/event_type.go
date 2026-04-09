package realtime

type EventType string

const (
	EventTypeMarketCreated EventType = "market.created"
	EventTypeMarketUpdated EventType = "market.updated"
	EventTypeMarketSettled EventType = "market.settled"
)
