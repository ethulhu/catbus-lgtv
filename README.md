<!--
SPDX-FileCopyrightText: 2020 Ethel Morgan

SPDX-License-Identifier: MIT
-->

# Catbus LGTV

Control WebOS-based LG TVs with [Catbus](https://ethulhu.co.uk/catbus), a home automation platform built around MQTT.

## MQTT Topics

The control of each parameter of the TV is split into its own topic:

- power, either `on` or `off`.
- volume, as a percentage, from 0 to 100.
- app, as a set of user-provided values, or App IDs on the TV (e.g. `com.webos.app.hdmi1`).

## Configuration

The bridge is configured with a JSON file, containing:

- the broker host & port.
- the TV's host, MAC address, and a key.
- the topics for power, volume, and app.
- a set of meaningful names for App IDs.

For example,

```json
{
	"broker_host": "home-server.local",
	"broker_port": 1883,

	"tv": {
		"host": "192.168.0.42",
		"mac": "00:de:ed:ab:ca:ff",
		"key": "a key from the TV",
	},

	"topic_power": "home/living-room/tv/power",
	"topic_volume": "home/living-room/tv/volume_percent",
	"topic_app": "home/living-room/tv/input",

	"apps": {
		"XBMC": "com.webos.app.hdmi1",
		"MiraCast": "com.webos.app.miracast"
	}
}
```

## Keys

Without a key, you will need to approve the server's connection on the TV every time the server starts.

You can generate a key using `cmd/generate-key`, which will print one out. Put it in the config, and the server is pre-authorized from then on.
