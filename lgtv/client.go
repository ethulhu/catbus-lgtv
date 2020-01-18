package lgtv

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type (
	client struct {
		conn *websocket.Conn
		opts Options
		uri  string

		sync.Mutex
		sequence         int
		requestChannel   chan *request
		responseChannels map[int]chan *response
		connectionClosed chan struct{}
		errors           chan error

		appHandler    func(App)
		volumeHandler func(Volume)

		appNameForID map[string]string
	}
)

func NewClient(host string, opts Options) Client {
	return &client{
		uri:  fmt.Sprintf("ws://%v:3000", host),
		opts: opts,
	}
}

func (c *client) Connect(ctx context.Context) error {
	if c.IsConnected() {
		return nil
	}

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, c.uri, nil)
	if err != nil {
		return fmt.Errorf("could not dial %v: %w", c.uri, err)
	}

	c.Lock()
	defer c.Unlock()
	c.sequence = 0
	c.requestChannel = make(chan *request)
	c.responseChannels = map[int]chan *response{}
	c.connectionClosed = make(chan struct{})
	c.errors = make(chan error)

	conn.SetPongHandler(func(text string) error {
		return conn.SetReadDeadline(time.Now().Add(c.opts.PongTimeout))
	})
	c.conn = conn
	go c.readLoop()
	go c.writeLoop()

	return nil
}

func (c *client) IsConnected() bool {
	return c.conn != nil
}

func (c *client) Err() error {
	defer close(c.errors)
	return <-c.errors
}
func (c *client) SetAppHandler(f func(App))       { c.appHandler = f }
func (c *client) SetVolumeHandler(f func(Volume)) { c.volumeHandler = f }

func (c *client) readLoop() {
	defer close(c.connectionClosed)
	defer func() {
		c.Lock()
		defer c.Unlock()
		for _, ch := range c.responseChannels {
			close(ch)
		}
	}()
	for {
		data := &response{}
		if err := c.conn.ReadJSON(data); err != nil {
			c.errors <- fmt.Errorf("could not read from websocket: %v", err)
			return
		}
		c.Lock()
		// If noöne requested it, throw it away.
		if ch, ok := c.responseChannels[data.ID]; ok {
			ch <- data
		}
		c.Unlock()
	}
}
func (c *client) writeLoop() {
	ping := time.NewTicker(c.opts.pingPeriod())
	defer ping.Stop()
	defer close(c.requestChannel)

	for {
		select {
		case data := <-c.requestChannel:
			if err := c.conn.WriteJSON(data); err != nil {
				c.errors <- fmt.Errorf("could not write to websocket%v: %v", data, err)
				return
			}
		case <-ping.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.errors <- fmt.Errorf("could not ping websocket: %v", err)
				return
			}
		case <-c.connectionClosed:
			_ = c.conn.Close()
			c.conn = nil
			log.Print("connection closed")
			return
		}
	}
}

func (c *client) newRequest() (int, <-chan *response, func()) {
	c.Lock()
	defer c.Unlock()

	id := c.sequence
	c.sequence++

	ch := make(chan *response)
	c.responseChannels[id] = ch

	cancel := func() {
		c.Lock()
		defer c.Unlock()

		close(ch)
		delete(c.responseChannels, id)
	}

	return id, ch, cancel
}
func (c *client) receive(ctx context.Context, ch <-chan *response) (*response, error) {
	select {
	case rsp := <-ch:
		return rsp, rsp.Err()
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
