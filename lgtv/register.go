package lgtv

import (
	"context"
	"encoding/json"
)

type (
	registerRequest struct {
		PairingType string                  `json:"pairingType"`
		Manifest    registerRequestManifest `json:"manifest"`
		ClientKey   string                  `json:"client-key"`
	}
	registerRequestManifest struct {
		Permissions []string `json:"permissions"`
	}
	registerResponse struct {
		ClientKey string `json:"client-key"`
	}
)

var (
	permissions = []string{
		"CONTROL_AUDIO",
		"CONTROL_INPUT_MEDIA_PLAYBACK",
		"CONTROL_INPUT_TEXT",
		"CONTROL_INPUT_TV",
		"CONTROL_MOUSE_AND_KEYBOARD",
		"CONTROL_POWER",
		"LAUNCH",
		"READ_CURRENT_CHANNEL",
		"READ_INPUT_DEVICE_LIST",
		"READ_INSTALLED_APPS",
		"READ_LGE_SDX",
		"READ_LGE_TV_INPUT_EVENTS",
		"READ_NOTIFICATIONS",
		"READ_RUNNING_APPS",
		"READ_TV_CHANNEL_LIST",
		"READ_TV_CURRENT_TIME",
		"READ_UPDATE_INFO",
		"SEARCH",
		"TEST_SECURE",
		"UPDATE_FROM_REMOTE_APP",
		"WRITE_NOTIFICATION_TOAST",
		"WRITE_SETTINGS",
	}
)

func (c *client) Register(ctx context.Context, key string) (string, error) {
	id, rspChan, cancel := c.newRequest()
	defer cancel()

	req := &request{
		ID:   id,
		Type: requestTypeRegister,
		Payload: registerRequest{
			PairingType: "PROMPT",
			Manifest: registerRequestManifest{
				Permissions: permissions,
			},
			ClientKey: key,
		},
	}
	c.requestChannel <- req

	rsp, err := c.receive(ctx, rspChan)
	if err != nil {
		return key, err
	}
	if rsp.Type == responseTypeRegistered {
		return key, nil
	}

	rsp, err = c.receive(ctx, rspChan)
	if err != nil {
		return key, err
	}
	payload := registerResponse{}
	if err := json.Unmarshal(rsp.Payload, &payload); err != nil {
		return key, err
	}
	return payload.ClientKey, nil
}
