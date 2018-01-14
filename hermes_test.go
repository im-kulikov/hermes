package hermes

import (
	"testing"
	"time"
)

func TestNewWebSocket(t *testing.T) {
	ws, err := NewWebSocket(
		Deadline(5 * time.Second),
	)
	if err != nil {
		t.Fatal(err)
	}

	defer ws.Close()

	if err = ws.Ping(); err != nil {
		t.Fatal(err)
	}

	if err = ws.Pong(); err != nil {
		t.Fatal(err)
	}

	if err = ws.Subscribe("test-channel"); err != nil {
		t.Fatal(err)
	}

	if err = ws.SendTextMessage([]byte("test-message")); err != nil {
		t.Fatal(err)
	}
}
