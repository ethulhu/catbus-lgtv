// SPDX-FileCopyrightText: 2020 Ethel Morgan
//
// SPDX-License-Identifier: MIT

package lgtv

import (
	"context"
	"encoding/json"
	"log"
)

func (c *client) ListApps(ctx context.Context) ([]App, error) {
	id, rspChan, cancel := c.newRequest()
	defer cancel()

	req := &request{
		ID:   id,
		Type: requestTypeRequest,
		URI:  listApps,
	}
	c.requestChannel <- req

	rsp, err := c.receive(ctx, rspChan)
	if err != nil {
		return nil, err
	}

	payload := listAppsResponse{}
	if err := json.Unmarshal(rsp.Payload, &payload); err != nil {
		return nil, err
	}
	return payload.Apps, nil
}
func (c *client) App(ctx context.Context) (App, error) {
	id, rspChan, cancel := c.newRequest()
	defer cancel()

	req := &request{
		ID:   id,
		Type: requestTypeRequest,
		URI:  getApp,
	}
	c.requestChannel <- req

	rsp, err := c.receive(ctx, rspChan)
	if err != nil {
		return App{}, err
	}
	payload := getAppResponse{}
	if err := json.Unmarshal(rsp.Payload, &payload); err != nil {
		return App{}, err
	}
	// TODO: cache a copy of all App names in the client object.
	return App{ID: payload.ID}, nil
}
func (c *client) SubscribeApp(f func(App)) {
	id, rspChan, _ := c.newRequest()

	req := &request{
		ID:   id,
		Type: requestTypeSubscribe,
		URI:  getApp,
	}
	c.requestChannel <- req

	go func() {
		for rsp := range rspChan {
			if err := rsp.Err(); err != nil {
				log.Printf("error recieved from TV waiting for App events: %v", err)
				continue

			}

			payload := getAppResponse{}
			if err := json.Unmarshal(rsp.Payload, &payload); err != nil {
				log.Printf("could not unmarshal App payload: %v", err)
				continue
			}
			app := App{ID: payload.ID}
			go f(app)
		}
	}()
}
func (c *client) SetApp(ctx context.Context, appID string) error {
	id, rspChan, cancel := c.newRequest()
	defer cancel()

	req := &request{
		ID:   id,
		Type: requestTypeRequest,
		URI:  setApp,
		Payload: struct {
			ID string `json:"id"`
		}{appID},
	}
	c.requestChannel <- req

	_, err := c.receive(ctx, rspChan)
	return err
}

func (c *client) Volume(ctx context.Context) (Volume, error) {
	id, rspChan, cancel := c.newRequest()
	defer cancel()

	req := &request{
		ID:   id,
		Type: requestTypeRequest,
		URI:  getVolume,
	}
	c.requestChannel <- req

	rsp, err := c.receive(ctx, rspChan)
	if err != nil {
		return Volume{}, err
	}

	payload := Volume{}
	if err := json.Unmarshal(rsp.Payload, &payload); err != nil {
		return Volume{}, err
	}
	return payload, nil
}
func (c *client) SubscribeVolume(f func(Volume)) {
	id, rspChan, _ := c.newRequest()

	req := &request{
		ID:   id,
		Type: requestTypeSubscribe,
		URI:  getVolume,
	}
	c.requestChannel <- req

	go func() {
		for rsp := range rspChan {
			if err := rsp.Err(); err != nil {
				log.Printf("error recieved from TV waiting for Volume events: %v", err)
				continue

			}

			payload := Volume{}
			if err := json.Unmarshal(rsp.Payload, &payload); err != nil {
				log.Printf("could not unmarshal Volume payload: %v", err)
				continue
			}
			go f(payload)
		}
	}()
}
func (c *client) SetVolume(ctx context.Context, volume int) error {
	id, rspChan, cancel := c.newRequest()
	defer cancel()

	req := &request{
		ID:      id,
		Type:    requestTypeRequest,
		URI:     setVolume,
		Payload: setVolumeRequest{volume},
	}
	c.requestChannel <- req

	_, err := c.receive(ctx, rspChan)
	return err
}

func (c *client) TurnOff(ctx context.Context) error {
	id, rspChan, cancel := c.newRequest()
	defer cancel()

	req := &request{
		ID:   id,
		Type: requestTypeRequest,
		URI:  turnOff,
	}
	c.requestChannel <- req

	_, err := c.receive(ctx, rspChan)
	return err
}
