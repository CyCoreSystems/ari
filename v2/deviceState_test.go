package ari

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type DeviceTests struct {
	suite.Suite
	cancel context.CancelFunc
}

func (s *DeviceTests) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	DefaultClient.Listen(ctx)
	time.Sleep(1 * time.Second)
}

func (s *DeviceTests) TearDownSuite() {
	s.cancel()
}

func (s *DeviceTests) TestListDeviceStates() {
	blah, err := DefaultClient.SubscribeApplication("default", "endpoint:PJSIP/101")
	s.Nil(err, "Oops")
	fmt.Printf("%+v", blah)

	time.Sleep(2 * time.Second)

	list, err := DefaultClient.ListDeviceStates()
	var dev DeviceState
	s.Nil(err, "Error getting list of deviceStates")
	for _, element := range list {
		fmt.Println("deviceState: ", element)
	}
	if len(list) > 0 {
		dev = list[0]
	} else {
		return
	}

	device, err := DefaultClient.GetDeviceState(dev.Name)
	s.Nil(err, "Unable to get DeviceState for "+dev.Name)
	s.Equal(dev.State, device.State)

	err = DefaultClient.ChangeDeviceState(dev.Name, "busy")
	s.Nil(err, "Error changing deviceState.")

	err = DefaultClient.DeleteDeviceState(dev.Name)
	s.Nil(err, "Error deleting deviceState.")

	_, err = DefaultClient.UnsubscribeApplication("default", "endpoint:PJSIP/101")
}

func TestDeviceSuite(t *testing.T) {
	suite.Run(t, new(DeviceTests))
}
