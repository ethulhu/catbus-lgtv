// SPDX-FileCopyrightText: 2020 Ethel Morgan
//
// SPDX-License-Identifier: MIT

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
		// Register registers the Client with the TV.
		// This may may require manual approval on the TV itself.
		// If registration fails, Register() will return the original key.
		Register(context.Context, string) (string, error)

		// ListApp lists all apps on the TV.
		ListApps(context.Context) ([]App, error)
		// App gets the current app on the TV.
		App(context.Context) (App, error)
		// SetApp sets the current app ID on the TV.
		SetApp(context.Context, string) error
		// SubscribeApp listens for App events.
		SubscribeApp(func(App))

		// Volume gets the current volume on the TV.
		Volume(context.Context) (Volume, error)
		// SetVolume sets the current volume percentage on the TV.
		SetVolume(context.Context, int) error
		// SubscribeVolume listens for Volume events.
		SubscribeVolume(func(Volume))

		// TurnOff turns off the TV.
		TurnOff(context.Context) error

		// Wait blocks until the connection to the TV is closed.
		Wait() error

		// Close closes the connection to the TV.
		Close() error
	}

	Options struct {
		PongTimeout time.Duration
	}

	App struct {
		Name string `json:"title"`
		ID   string `json:"id"`
	}

	Volume struct {
		Percent int  `json:"volume"`
		Muted   bool `json:"muted"`
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

func (e TVError) Error() string {
	return e.message
}
