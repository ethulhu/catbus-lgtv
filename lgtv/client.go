// SPDX-FileCopyrightText: 2020 Ethel Morgan
//
// SPDX-License-Identifier: MIT

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

		sync.Mutex
		sequence         int
		requestChannel   chan *request
		responseChannels map[int]chan *response
		connectionClosed chan struct{}
		errors           chan error

		appNameForID map[string]string
	}
)

func Dial(ctx context.Context, host string, opts Options) (Client, error) {
	uri := fmt.Sprintf("ws://%v:3000", host)

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("could not dial %v: %w", uri, err)
	}
	conn.SetPongHandler(func(_ string) error {
		return conn.SetReadDeadline(time.Now().Add(opts.PongTimeout))
	})

	c := &client{
		conn: conn,

		requestChannel:   make(chan *request),
		responseChannels: map[int]chan *response{},
		connectionClosed: make(chan struct{}),
		errors:           make(chan error),
	}
	go c.readLoop()
	go c.writeLoop(opts.pingPeriod())

	return c, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) Wait() error {
	defer close(c.errors)
	return <-c.errors
}

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
func (c *client) writeLoop(pingPeriod time.Duration) {
	ping := time.NewTicker(pingPeriod)
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
