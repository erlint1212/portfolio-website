package routing

import "time"

type GameLog struct {
	CurrentTime time.Time
	Message     string
	// Username    string
}

type Project struct {
	Title       string
	Category    string
	SourceURL   string
	IsLive      bool
	GamePath    string
	Description string
	Tags        []string
}
