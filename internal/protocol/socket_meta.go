package protocol

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

type SocketMeta struct {
	Port  uint16
	Nonce [16]byte
}

func ParseSocketMeta(data []byte) (SocketMeta, error) {
	if strings.HasPrefix(string(data), "!<socket >") {
		return parseCygwinSocketMeta(strings.TrimPrefix(string(data), "!<socket >"))
	}
	return parseStandardSocketMeta(data)
}

func parseStandardSocketMeta(data []byte) (SocketMeta, error) {
	if len(data) < 17 {
		return SocketMeta{}, fmt.Errorf("invalid socket metadata length: %d", len(data))
	}
	left := strings.TrimSpace(string(data[:len(data)-16]))
	port, err := strconv.ParseUint(left, 10, 16)
	if err != nil {
		return SocketMeta{}, fmt.Errorf("parse socket port: %w", err)
	}

	var nonce [16]byte
	copy(nonce[:], data[len(data)-16:])
	return SocketMeta{Port: uint16(port), Nonce: nonce}, nil
}

func parseCygwinSocketMeta(raw string) (SocketMeta, error) {
	parts := strings.SplitN(raw, " ", 3)
	if len(parts) != 3 {
		return SocketMeta{}, fmt.Errorf("invalid cygwin socket metadata")
	}
	port, err := strconv.ParseUint(parts[0], 10, 16)
	if err != nil {
		return SocketMeta{}, fmt.Errorf("parse cygwin port: %w", err)
	}
	if parts[1] != "s" {
		return SocketMeta{}, fmt.Errorf("invalid cygwin socket marker: %q", parts[1])
	}

	tail := parts[2]
	if !strings.HasSuffix(tail, "x") {
		return SocketMeta{}, fmt.Errorf("invalid cygwin nonce terminator")
	}
	segments := strings.Split(strings.TrimSuffix(tail, "x"), "-")
	if len(segments) != 4 {
		return SocketMeta{}, fmt.Errorf("invalid cygwin nonce segments")
	}

	var nonce [16]byte
	for idx, segment := range segments {
		value, err := strconv.ParseUint(segment, 16, 32)
		if err != nil {
			return SocketMeta{}, fmt.Errorf("parse cygwin nonce segment %d: %w", idx, err)
		}
		binary.LittleEndian.PutUint32(nonce[idx*4:], uint32(value))
	}
	return SocketMeta{Port: uint16(port), Nonce: nonce}, nil
}
