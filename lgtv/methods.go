package lgtv

import (
	"context"
	"encoding/json"
	"log"
)

func (c *client) ListApps(ctx context.Context) ([]App, error) {
	if !c.IsConnected() {
		return nil, ErrNotConnected
	}

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
	if !c.IsConnected() {
		return App{}, ErrNotConnected
	}

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
func (c *client) SubscribeApp() func() {
	id, rspChan, cancel := c.newRequest()

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
			c.appHandler(app)
		}
	}()
	return cancel
}
func (c *client) SetApp(ctx context.Context, app App) error {
	if !c.IsConnected() {
		return ErrNotConnected
	}

	id, rspChan, cancel := c.newRequest()
	defer cancel()

	req := &request{
		ID:   id,
		Type: requestTypeRequest,
		URI:  setApp,
		Payload: struct {
			ID string `json:"id"`
		}{app.ID},
	}
	c.requestChannel <- req

	_, err := c.receive(ctx, rspChan)
	return err
}

func (c *client) Volume(ctx context.Context) (Volume, error) {
	if !c.IsConnected() {
		return Volume{}, ErrNotConnected
	}

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
func (c *client) SubscribeVolume() func() {
	id, rspChan, cancel := c.newRequest()

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
			c.volumeHandler(payload)
		}
	}()
	return cancel
}
func (c *client) SetVolume(ctx context.Context, volume int) error {
	if !c.IsConnected() {
		return ErrNotConnected
	}

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
	if !c.IsConnected() {
		return ErrNotConnected
	}

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
