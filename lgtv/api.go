// Package lgtv is a client for WebOS LG TVs.
package lgtv

import (
	"context"
	"errors"
	"time"
)

type (
	// Client is a reusable WebOS LG TV client.
	Client interface {
		// Connect connects to the TV, with optional Key parameter.
		Connect(context.Context) error
		// Register registers the Client with the TV.
		// This may may require manual approval on the TV itself.
		// If registration fails, Register() will return the original key.
		Register(context.Context, string) (string, error)
		// Err blocks until a connection error is returned.
		// Once an error is returned, the Client must be reconnected.
		Err() error
		// IsConnected returns whether the Client is connected to the TV.
		IsConnected() bool

		// ListApp lists all apps on the TV.
		ListApps(context.Context) ([]App, error)
		// App gets the current app on the TV.
		App(context.Context) (App, error)
		// SetApp sets the current app on the TV.
		SetApp(context.Context, App) error
		// SetAppHandler sets the hander for App event subscriptions.
		SetAppHandler(func(App))
		// SubscribeApp starts listening for App events, and returns a cancel function.
		SubscribeApp() func()

		// Volume gets the current volume on the TV.
		Volume(context.Context) (Volume, error)
		// SetVolume sets the current volume on the TV.
		SetVolume(context.Context, int) error
		// SetVolumeHandler sets the hander for Volume event subscriptions.
		SetVolumeHandler(func(Volume))
		// SubscribeVolume starts listening for Volume events, and returns a cancel function.
		SubscribeVolume() func()

		// TurnOff turns off the TV.
		TurnOff(context.Context) error
	}

	Options struct {
		PongTimeout time.Duration
	}

	App struct {
		Name string `json:"title"`
		ID   string `json:"id"`
	}

	Volume struct {
		Level int  `json:"volume"`
		Muted bool `json:"muted"`
	}

	// TVError is an error returned by the TV, e.g. about invalid messages.
	// Connection errors will always return via Client.Err().
	TVError struct {
		message string
	}
)

var (
	DefaultOptions = Options{
		PongTimeout: 10 * time.Second,
	}

	ErrNotConnected = errors.New("not connected to TV")
)

// pingPeriod must be less than PongTimeout.
func (o *Options) pingPeriod() time.Duration {
	return (o.PongTimeout * 9) / 10
}

func (e *TVError) Error() string {
	return e.message
}
