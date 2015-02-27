package ari

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ApplicationTests struct {
	suite.Suite
	c    chan *Event
	list []Application
}

func (s *ApplicationTests) SetupSuite() {
	// Start event listener/logger
	go s.LogEvent()

	// Connect to Asterisk
	DefaultClient.Go()
}

func (s *ApplicationTests) TearDownSuite() {
	DefaultClient.Close()
}

func (s *ApplicationTests) LogEvent() {
	defer s.LogEvent()
	e := <-DefaultClient.Events
	s.NotNil(e.Application, "Event's application name must exist")
}

var List []Application

func (s *ApplicationTests) TestListApplications() {
	var err error
	s.list, err = DefaultClient.ListApplications()
	s.Nil(err, "ListApplications must not return an error")
	s.NotNil(s.list, "Application list must not be empty")
}

func (s *ApplicationTests) TestGetApplication() {
	list, err := DefaultClient.ListApplications()
	if err != nil {
		log.Println("Failed to get list of applications; skipping test")
		return
	}
	if len(list) == 0 {
		log.Println("No applications in list; skipping test")
		return
	}
	fmt.Println("Got application list:", list)
	a, err := DefaultClient.GetApplication(list[0].Name)
	s.Nil(err, "GetApplication must not return an error")
	s.NotNil(a.Name, "GetApplication must return a Name")
}

func (s *ApplicationTests) TestSubscribeApplication() {
	_, err := DefaultClient.SubscribeApplication("default", "endpoint:PJSIP/101")
	s.Nil(err, "SubscribeApplication must not return an error")
}

func (s *ApplicationTests) TestUnsubscribeApplication() {
	_, err := DefaultClient.UnsubscribeApplication("default", "endpoint:PJSIP/101")
	s.Nil(err, "UnsubscribeApplication must not return an error")
}

func TestApplicationSuite(t *testing.T) {
	suite.Run(t, new(ApplicationTests))
}
