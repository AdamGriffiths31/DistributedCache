package protocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSetCommand(t *testing.T) {
	cmd := &CommandSet{
		Key:   []byte("Foo"),
		Value: []byte("Bar"),
		TTL:   2,
	}

	r := bytes.NewReader(cmd.Bytes())
	parsedCmd, err := ParseCommand(r)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, cmd, parsedCmd)
}

func TestParseGetCommand(t *testing.T) {
	cmd := &CommandGet{
		Key: []byte("Foo"),
	}

	r := bytes.NewReader(cmd.Bytes())
	parsedCmd, err := ParseCommand(r)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, cmd, parsedCmd)
}

func TestInvalidCommand(t *testing.T) {
	cmd := []byte{0x00, 0x00, 0x00, 0x00}
	r := bytes.NewReader(cmd)
	_, err := ParseCommand(r)
	assert.Error(t, err)
}
