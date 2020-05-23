// SPDX-FileCopyrightText: 2020 Ethel Morgan
//
// SPDX-License-Identifier: MIT

package lgtv

import "encoding/json"

type (
	requestType  string
	responseType string

	uri string

	request struct {
		ID      int         `json:"id"`
		Type    requestType `json:"type"`
		URI     uri         `json:"uri,omitempty"`
		Payload interface{} `json:"payload,omitempty"`
	}

	response struct {
		ID      int             `json:"id"`
		Type    responseType    `json:"type"`
		Error   string          `json:"error"`
		Payload json.RawMessage `json:"payload"`
	}

)

const (
	requestTypeRegister  = requestType("register")
	requestTypeRequest   = requestType("request")
	requestTypeSubscribe = requestType("subscribe")

	responseTypeError      = responseType("error")
	responseTypeRegistered = responseType("registered")

	listApps  = uri("ssap://com.webos.applicationManager/listApps")
	getApp    = uri("ssap://com.webos.applicationManager/getForegroundAppInfo")
	setApp    = uri("ssap://system.launcher/launch")
	getVolume = uri("ssap://audio/getVolume")
	setVolume = uri("ssap://audio/setVolume")
	turnOff   = uri("ssap://system/turnOff")
)

func (rsp *response) Err() error {
	if rsp.Type == responseTypeError {
		return &TVError{message: rsp.Error}
	}
	return nil
}

// Below here lie the many many JSON definitions.
type (
	listAppsResponse struct {
		Apps []App `json:"apps"`
	}
	getAppResponse struct {
		ID string `json:"appId"`
	}
	setVolumeRequest struct {
		Level int `json:"volume"`
	}
)