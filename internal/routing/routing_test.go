package routing

import (
	"testing"
	"time"
)

// Constants
func TestSlugAndExchangeConstants(t *testing.T) {
	t.Parallel()

	if GameLogSlug == "" {
		t.Fatal("GameLogSlug must not be empty")
	}
	if ExchangePortfolioDirect == "" {
		t.Fatal("ExchangePortfolioDirect must not be empty")
	}
	if ExchangePortfolioTopic == "" {
		t.Fatal("ExchangePortfolioTopic must not be empty")
	}
}

func TestExchangeNamesAreDifferent(t *testing.T) {
	t.Parallel()

	if ExchangePortfolioDirect == ExchangePortfolioTopic {
		t.Fatalf("Direct and Topic exchange names must differ, both are %q", ExchangePortfolioDirect)
	}
}

// SimpleQueueType enum
func TestSimpleQueueType_IotaValues(t *testing.T) {
	t.Parallel()

	if Transient != 0 {
		t.Errorf("expected Transient == 0, got %d", Transient)
	}
	if Durable != 1 {
		t.Errorf("expected Durable == 1, got %d", Durable)
	}
}

func TestSimpleQueueType_AreDistinct(t *testing.T) {
	t.Parallel()

	if Transient == Durable {
		t.Error("Transient and Durable must have different underlying values")
	}
}

// AckType enum
func TestAckType_IotaValues(t *testing.T) {
	t.Parallel()

	if Ack != 0 {
		t.Errorf("expected Ack == 0, got %d", Ack)
	}
	if NackRequeue != 1 {
		t.Errorf("expected NackRequeue == 1, got %d", NackRequeue)
	}
	if NackDiscard != 2 {
		t.Errorf("expected NackDiscard == 2, got %d", NackDiscard)
	}
}

func TestAckType_AllDistinct(t *testing.T) {
	t.Parallel()

	seen := map[AckType]string{
		Ack:         "Ack",
		NackRequeue: "NackRequeue",
		NackDiscard: "NackDiscard",
	}
	if len(seen) != 3 {
		t.Error("AckType constants must all be distinct")
	}
}

// GameLog struct
func TestGameLog_FieldsRoundTrip(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	gl := GameLog{
		CurrentTime: now,
		Message:     "player joined",
	}

	if gl.CurrentTime != now {
		t.Errorf("CurrentTime: want %v, got %v", now, gl.CurrentTime)
	}
	if gl.Message != "player joined" {
		t.Errorf("Message: want %q, got %q", "player joined", gl.Message)
	}
}

func TestGameLog_ZeroValue(t *testing.T) {
	t.Parallel()

	var gl GameLog
	if !gl.CurrentTime.IsZero() {
		t.Error("zero-value GameLog.CurrentTime should be zero time")
	}
	if gl.Message != "" {
		t.Error("zero-value GameLog.Message should be empty")
	}
}

// Project struct
func TestProject_AllFieldsPopulated(t *testing.T) {
	t.Parallel()

	p := Project{
		Title:       "My Project",
		Category:    "/backend/api",
		SourceURL:   "https://github.com/example",
		IsLive:      true,
		GamePath:    "/games/test",
		Description: "A test project",
		Tags:        []string{"Go", "Docker"},
	}

	if p.Title != "My Project" {
		t.Errorf("Title mismatch: %q", p.Title)
	}
	if p.Category != "/backend/api" {
		t.Errorf("Category mismatch: %q", p.Category)
	}
	if p.SourceURL != "https://github.com/example" {
		t.Errorf("SourceURL mismatch: %q", p.SourceURL)
	}
	if !p.IsLive {
		t.Error("IsLive should be true")
	}
	if p.GamePath != "/games/test" {
		t.Errorf("GamePath mismatch: %q", p.GamePath)
	}
	if p.Description != "A test project" {
		t.Errorf("Description mismatch: %q", p.Description)
	}
	if len(p.Tags) != 2 || p.Tags[0] != "Go" || p.Tags[1] != "Docker" {
		t.Errorf("Tags mismatch: %v", p.Tags)
	}
}

func TestProject_ZeroValue(t *testing.T) {
	t.Parallel()

	var p Project
	if p.IsLive {
		t.Error("zero-value Project.IsLive should be false")
	}
	if p.Tags != nil {
		t.Error("zero-value Project.Tags should be nil")
	}
}
