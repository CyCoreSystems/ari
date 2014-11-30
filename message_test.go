package ari

import (
	"io/ioutil"
	"testing"

	"github.com/golang/glog"
	"github.com/stretchr/testify/assert"
)

var ariStasisStartEvent []byte
var stateChange []byte

func init() {
	var err error
	ariStasisStartEvent, err = ioutil.ReadFile("test_events/stasis_start.json")
	if err != nil {
		glog.Fatalln("Failed to open test event file:", err.Error())
	}
	stateChange, err = ioutil.ReadFile("test_events/channelStateChange.json")
	if err != nil {
		glog.Fatalln("Failed to open test event file:", err.Error())
	}
}

func TestDecodeAs(t *testing.T) {
	assert := assert.New(t)

	m, err := NewMessage(ariStasisStartEvent)
	assert.Nil(err, "Construction of new message must succeed")
	assert.Equal("StasisStart", m.Type, "Message type should be StatisStart")

	// Decode as StasisStart
	var sm StasisStart
	err = m.DecodeAs(&sm)
	assert.Nil(err, "DecodeAs StasisStart of message must succeed")
	assert.Equal("test", sm.Channel.Dialplan.Context, "Channel context decoded")
}

func TestDecodeAsChannelStateChange(t *testing.T) {
	assert := assert.New(t)

	m, err := NewMessage(stateChange)
	assert.Nil(err, "Construction of new message must succeed")
	assert.Equal("ChannelStateChange", m.Type, "Message type should be ChannelStateChange")

	// Decode as StasisStart
	var sc ChannelStateChange
	err = m.DecodeAs(&sc)
	assert.Nil(err, "DecodeAs ChannelStateChange of message must succeed")
	assert.Equal("Ringing", sc.Channel.State, "Channel context decoded")
}
