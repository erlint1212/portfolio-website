package routing

const (
	GameLogSlug = "game_logs"
)

const (
	ExchangePortfolioDirect = "portfolio_direct"
	ExchangePortfolioTopic  = "portfolio_topic"
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

