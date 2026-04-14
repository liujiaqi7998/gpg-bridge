package protocol

import (
	"bytes"
	"testing"
)

func TestParseSocketMetaParsesStandardFormat(t *testing.T) {
	nonce := []byte("0123456789ABCDEF")
	data := append([]byte("4321\n"), nonce...)

	meta, err := ParseSocketMeta(data)
	if err != nil {
		t.Fatalf("ParseSocketMeta returned error: %v", err)
	}
	if meta.Port != 4321 {
		t.Fatalf("unexpected port: %d", meta.Port)
	}
	if !bytes.Equal(meta.Nonce[:], nonce) {
		t.Fatalf("unexpected nonce: %v", meta.Nonce)
	}
}

func TestParseSocketMetaParsesCygwinFormat(t *testing.T) {
	data := []byte("!<socket >4321 s 0001-0002-0003-0004x")

	meta, err := ParseSocketMeta(data)
	if err != nil {
		t.Fatalf("ParseSocketMeta returned error: %v", err)
	}
	if meta.Port != 4321 {
		t.Fatalf("unexpected port: %d", meta.Port)
	}

	expected := [16]byte{0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00}
	if meta.Nonce != expected {
		t.Fatalf("unexpected nonce: %#v", meta.Nonce)
	}
}

func TestParseSocketMetaRejectsInvalidPayload(t *testing.T) {
	_, err := ParseSocketMeta([]byte("bad"))
	if err == nil {
		t.Fatal("expected invalid payload to fail")
	}
}
