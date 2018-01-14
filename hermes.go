package hermes

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocket struct {
	ws     *websocket.Conn
	quit   chan bool
	Stream chan *Event
	Errors chan error
}

type Event struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func (s *WebSocket) Close() {
	s.quit <- true
}

func (s *WebSocket) Subscribe(channel string) error {
	a := &Event{
		Event: "pusher:subscribe",
		Data: map[string]interface{}{
			"channel": channel,
		},
	}
	return s.ws.WriteJSON(a)
}

func (s *WebSocket) SendTextMessage(message []byte) error {
	return s.ws.WriteMessage(websocket.TextMessage, message)
}

func (s *WebSocket) Ping() error {
	a := &Event{
		Event: "pusher:ping",
	}
	return s.ws.WriteJSON(a)
}

func (s *WebSocket) Pong() error {
	a := &Event{
		Event: "pusher:pong",
	}
	return s.ws.WriteJSON(a)
}

func NewWebSocket(opts ...Option) (*WebSocket, error) {
	var (
		err     error
		options = newOptions(opts...)
		s       = &WebSocket{
			quit:   make(chan bool, 1),
			Stream: make(chan *Event),
			Errors: make(chan error),
		}
	)

	// set up websocket
	s.ws, _, err = websocket.DefaultDialer.Dial(options.url, nil)
	if err != nil {
		return nil, fmt.Errorf("error dialing websocket: %s", err)
	}

	go func() {
		defer func() {
			if err = s.ws.Close(); err != nil {
				s.Errors <- err
			}
		}()
		for {
			runtime.Gosched()
			if err = s.ws.SetReadDeadline(time.Now().Add(options.deadline)); err != nil {
				s.Errors <- err
			}
			select {
			case <-s.quit:
				return
			default:
				var message []byte
				var err error
				_, message, err = s.ws.ReadMessage()
				if err != nil {
					s.Errors <- err
					continue
				}
				e := &Event{}
				err = json.Unmarshal(message, e)
				if err != nil {
					s.Errors <- err
					continue
				}
				s.Stream <- e
			}
		}
	}()

	return s, nil
}
