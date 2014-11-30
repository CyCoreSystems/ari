package ari

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	err := DefaultClient.SendMessage("PJSIP/101", "PJSIP", "102", "Hello", nil)
	assert.Nil(t, err, "Test message failed to send")
}

func TestSendMessageByUri(t *testing.T) {
	vars := map[string]string{"testme": "testmeVal"}
	err := DefaultClient.SendMessageByUri("PJSIP/101", "http://localhost:8088/endpoints/PJSIP/102", "Hello", vars)
	assert.Nil(t, err, "Error sending message to specific endpoint")
}
