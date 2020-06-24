<!--
SPDX-FileCopyrightText: 2020 Ethel Morgan

SPDX-License-Identifier: MIT
-->

# Catbus LGTV

Control WebOS-based LG TVs with [Catbus](https://ethulhu.co.uk/catbus), a home automation platform built around MQTT.

## MQTT Topics

The control of each parameter of the TV is split into its own topic:

- Power, either `on` or `off`.
- Volume, as a percentage, from 0 to 100.
- App, as a set of user-provided values, or App IDs on the TV (e.g. `com.webos.app.hdmi1`).

## Configuration

The bridge is configured with a JSON file, containing:

- The broker host & port.
- The TV's host and a [key](#keys).
- The topics for power, volume, app, and app enum values.
- A set of meaningful names for App IDs.

For example,

```json
{
	"mqttBroker": "tcp://home-server.local:1883",

	"tv": {
		"host": "192.168.0.42",
		"key": "a key from the TV",
	},

        "topics": {
		"app": "home/living-room/tv/input_enum",
		"appValues": "home/living-room/tv/input_enum/values",
		"power": "home/living-room/tv/power",
		"volume": "home/living-room/tv/volume_percent"
        },

	"apps": {
		"XBMC": "com.webos.app.hdmi1",
		"MiraCast": "com.webos.app.miracast"
	}
}
```

## Keys

Without a key, you will need to approve the server's connection on the TV every time the server starts.

You can generate a key using `cmd/generate-key`, which will print one out.
Put it in the config, and the server is pre-authorized from then on.
