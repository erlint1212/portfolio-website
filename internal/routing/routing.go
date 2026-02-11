package routing

const (
	ExchangePortfoilioDirect = "portfolio_direct"
	ExchangePortfoilioTopic  = "portfolio_topic"
)

type SimpleQueueType int

const (
	Transient SimpleQueueType = iota
	Durable
)

type AckType int

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

