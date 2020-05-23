// SPDX-FileCopyrightText: 2020 Ethel Morgan
//
// SPDX-License-Identifier: MIT

// Package wol implements Wake-On-Lan.
package wol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

var (
	header = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

	broadcastAddr = &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: 9,
	}
)

// Wake sends the Magic Packet to wake the computer with the given MAC address.
func Wake(mac net.HardwareAddr) error {
	conn, err := net.DialUDP("udp", nil, broadcastAddr)
	if err != nil {
		return fmt.Errorf("could not create UDP broadcast: %w", err)
	}
	defer conn.Close()

	payload := packet(mac)

	if _, err := conn.Write(payload); err != nil {
		return fmt.Errorf("could not send Magic Packet: %w", err)
	}
	return nil
}

func packet(mac net.HardwareAddr) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, header)
	for i := 0; i < 16; i++ {
		binary.Write(&buf, binary.BigEndian, mac)
	}
	return buf.Bytes()
}