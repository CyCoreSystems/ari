package ari

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	err := DefaultClient.SendMessage("PJSIP/101", "PJSIP", "101", "Hello", nil)
	assert.Nil(t, err, "Test message failed to send")
}

func TestSendMessageByUri(t *testing.T) {
	vars := map[string]string{"testme": "testmeVal"}
	err := DefaultClient.SendMessageByUri("PJSIP/101", "sip:101@localhost", "Hello", vars)
	assert.Nil(t, err, "Error sending message to specific endpoint")
}
