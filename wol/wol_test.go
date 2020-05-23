// SPDX-FileCopyrightText: 2020 Ethel Morgan
//
// SPDX-License-Identifier: MIT

package wol

import (
	"encoding/hex"
	"net"
	"reflect"
	"testing"
)

func TestPacket(t *testing.T) {
	tests := []struct {
		mac    string
		packet string
	}{
		{
			"a8:23:22:ad:be:c7",
			"ffffffffffffa82322adbec7a82322adbec7a82322adbec7a82322adbec7a82322adbec7a82322adbec7a82322adbec7a82322adbec7a82322adbec7a82322adbec7a82322adbec7a82322adbec7a82322adbec7a82322adbec7a82322adbec7a82322adbec7",
		},
	}

	for i, test := range tests {
		mac, err := net.ParseMAC(test.mac)
		if err != nil {
			t.Fatalf("invalid MAC address %q: %v", test.mac, err)
		}

		if len(test.packet) != 102 {
			if err != nil {
				t.Fatalf("invalid packet length address: found %v, want 102", len(test.packet))
			}
		}
		want, err := hex.DecodeString(test.packet)
		if err != nil {
			t.Fatalf("invalid hex string: %v", err)
		}

		got := packet(mac)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("[test %d]\nwanted %v\nlen %d\ngot %v\nlen %d", i, want, len(want), got, len(got))
		}
	}
}