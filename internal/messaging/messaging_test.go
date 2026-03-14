package messaging

import (
	"bytes"
	"encoding/gob"
	"testing"
	"time"

	"github.com/erlint1212/portfolio/internal/routing"
)

// UnmarshalGob
func gobEncode(t *testing.T, v any) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(v); err != nil {
		t.Fatalf("gobEncode helper: %v", err)
	}
	return buf.Bytes()
}

func TestUnmarshalGob_GameLog_RoundTrip(t *testing.T) {
	t.Parallel()

	original := routing.GameLog{
		CurrentTime: time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC),
		Message:     "player scored",
	}
	data := gobEncode(t, original)

	decoded, err := UnmarshalGob[routing.GameLog](data)
	if err != nil {
		t.Fatalf("UnmarshalGob returned error: %v", err)
	}
	if !decoded.CurrentTime.Equal(original.CurrentTime) {
		t.Errorf("CurrentTime: want %v, got %v", original.CurrentTime, decoded.CurrentTime)
	}
	if decoded.Message != original.Message {
		t.Errorf("Message: want %q, got %q", original.Message, decoded.Message)
	}
}

func TestUnmarshalGob_EmptyMessage(t *testing.T) {
	t.Parallel()

	original := routing.GameLog{
		CurrentTime: time.Now().UTC(),
		Message:     "",
	}
	data := gobEncode(t, original)

	decoded, err := UnmarshalGob[routing.GameLog](data)
	if err != nil {
		t.Fatalf("UnmarshalGob returned error: %v", err)
	}
	if decoded.Message != "" {
		t.Errorf("expected empty message, got %q", decoded.Message)
	}
}

func TestUnmarshalGob_StringType(t *testing.T) {
	t.Parallel()

	original := "hello world"
	data := gobEncode(t, original)

	decoded, err := UnmarshalGob[string](data)
	if err != nil {
		t.Fatalf("UnmarshalGob returned error: %v", err)
	}
	if decoded != original {
		t.Errorf("want %q, got %q", original, decoded)
	}
}

func TestUnmarshalGob_IntType(t *testing.T) {
	t.Parallel()

	original := 42
	data := gobEncode(t, original)

	decoded, err := UnmarshalGob[int](data)
	if err != nil {
		t.Fatalf("UnmarshalGob returned error: %v", err)
	}
	if decoded != original {
		t.Errorf("want %d, got %d", original, decoded)
	}
}

func TestUnmarshalGob_InvalidData(t *testing.T) {
	t.Parallel()

	_, err := UnmarshalGob[routing.GameLog]([]byte("not valid gob"))
	if err == nil {
		t.Fatal("expected error for invalid gob data, got nil")
	}
}

func TestUnmarshalGob_EmptySlice(t *testing.T) {
	t.Parallel()

	_, err := UnmarshalGob[routing.GameLog]([]byte{})
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestUnmarshalGob_TypeMismatch(t *testing.T) {
	t.Parallel()

	// Encode a string, try to decode as GameLog
	data := gobEncode(t, "just a string")
	_, err := UnmarshalGob[routing.GameLog](data)
	if err == nil {
		t.Fatal("expected error for type mismatch, got nil")
	}
}

// WriteLog  (just verifies no error – output goes to log, not stdout)
func TestWriteLog_NoError(t *testing.T) {
	t.Parallel()

	gl := routing.GameLog{
		CurrentTime: time.Now().UTC(),
		Message:     "test event",
	}
	if err := WriteLog(gl); err != nil {
		t.Fatalf("WriteLog returned unexpected error: %v", err)
	}
}

func TestWriteLog_EmptyMessage(t *testing.T) {
	t.Parallel()

	gl := routing.GameLog{
		CurrentTime: time.Now().UTC(),
		Message:     "",
	}
	if err := WriteLog(gl); err != nil {
		t.Fatalf("WriteLog returned unexpected error: %v", err)
	}
}

func TestWriteLog_ZeroTime(t *testing.T) {
	t.Parallel()

	gl := routing.GameLog{
		Message: "zero time entry",
	}
	if err := WriteLog(gl); err != nil {
		t.Fatalf("WriteLog returned unexpected error: %v", err)
	}
}

// HandlerWriteLog
func TestHandlerWriteLog_ReturnsAckOnSuccess(t *testing.T) {
	t.Parallel()

	handler := HandlerWriteLog()
	gl := routing.GameLog{
		CurrentTime: time.Now().UTC(),
		Message:     "handler test",
	}

	result := handler(gl)
	if result != routing.Ack {
		t.Errorf("expected Ack (%d), got %d", routing.Ack, result)
	}
}

func TestHandlerWriteLog_ReturnsFunctionType(t *testing.T) {
	t.Parallel()

	handler := HandlerWriteLog()
	if handler == nil {
		t.Fatal("HandlerWriteLog returned nil")
	}
}

// Gob encode -> decode round-trip (simulates publish -> subscribe path)
func TestGobRoundTrip_GameLog(t *testing.T) {
	t.Parallel()

	original := routing.GameLog{
		CurrentTime: time.Date(2026, 3, 12, 14, 0, 0, 0, time.UTC),
		Message:     "round trip test",
	}

	// Encode (same logic as PublishGob)
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(original); err != nil {
		t.Fatalf("encode: %v", err)
	}

	// Decode (same logic as Subscribe's unmarshaller)
	decoded, err := UnmarshalGob[routing.GameLog](buf.Bytes())
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	if !decoded.CurrentTime.Equal(original.CurrentTime) {
		t.Errorf("CurrentTime mismatch: want %v, got %v", original.CurrentTime, decoded.CurrentTime)
	}
	if decoded.Message != original.Message {
		t.Errorf("Message mismatch: want %q, got %q", original.Message, decoded.Message)
	}
}

func TestGobRoundTrip_LargeMessage(t *testing.T) {
	t.Parallel()

	longMsg := ""
	for i := 0; i < 1000; i++ {
		longMsg += "a"
	}

	original := routing.GameLog{
		CurrentTime: time.Now().UTC(),
		Message:     longMsg,
	}
	data := gobEncode(t, original)

	decoded, err := UnmarshalGob[routing.GameLog](data)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if decoded.Message != longMsg {
		t.Errorf("message length: want %d, got %d", len(longMsg), len(decoded.Message))
	}
}
