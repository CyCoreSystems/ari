package ari

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type BridgeTests struct {
	suite.Suite
	c chan *Event
}

func (s *BridgeTests) SetupSuite() {
	DefaultClient.Go()
	time.Sleep(1 * time.Second)
}

func (s *BridgeTests) TearDownSuite() {
	DefaultClient.Close()
}

func (s *BridgeTests) TestBridgeCreate() {
	//Creation of empty bridge.
	bridge, err := DefaultClient.CreateBridge(CreateBridgeRequest{
		BridgeId: "testBridge",
		Type:     "mixing",
	})
	s.Nil(err, "Bridge creation failed.")

	//the channel created for User1's phone
	u1Chan := OriginateRequest{
		Endpoint:  "PJSIP/102",
		App:       "default",
		AppArgs:   "",
		ChannelId: "Chan1",
	}

	_, err = DefaultClient.CreateChannel(u1Chan)
	s.Nil(err, "Channel for User1 not created")

	e := <-DefaultClient.Events
	for e.Type != "StasisStart" {
		e = <-DefaultClient.Events
	}

	err = DefaultClient.AnswerChannel("Chan1")
	s.Nil(err)
	//User1 has picked up his phone, Asterisk has finished the channel creation

	//the channel created for User2's phone
	u2Chan := OriginateRequest{
		Endpoint:  "PJSIP/101",
		App:       "default",
		AppArgs:   "",
		ChannelId: "Chan2",
	}

	_, err = DefaultClient.CreateChannel(u2Chan)
	s.Nil(err, "Channel for User2 not created")

	e = <-DefaultClient.Events
	for e.Type != "StasisStart" {
		e = <-DefaultClient.Events
	}

	err = DefaultClient.AnswerChannel("Chan2")
	s.Nil(err)
	//User2 has picked up his phone, Asterisk has finished the channel creation

	req := AddChannelRequest{
		ChannelId: "Chan1,Chan2",
	}

	//Incorporated 'Get Bridge' test:
	getBridge, err := DefaultClient.GetBridge(bridge.Id)
	s.Nil(err)
	s.NotNil(getBridge)
	fmt.Println("Bridge retrieved: ", getBridge)

	err = DefaultClient.AddChannel(bridge.Id, req)

	s.Nil(err, "Channels not added to bridge.")

}

func (s *BridgeTests) TestPlayMoh() {
	err := DefaultClient.PlayMusicOnHold("testBridge", "default")
	s.Nil(err, "Couldn't play MOH to bridge")
	fmt.Println("Sleeping 1 second to play MOH to bridge.")
	time.Sleep(1 * time.Second)
	err = DefaultClient.BridgeStopMoh("testBridge")
	s.Nil(err, "Couldn't stop MOH on bridge.")

}

func (s *BridgeTests) TestZDeleteBridge() {
	err := DefaultClient.BridgeDelete("testBridge")
	s.Nil(err, "Couldn't delete bridge")
}

func (s *BridgeTests) TestBridgeList() {
	list, err := DefaultClient.ListBridges()
	s.Nil(err, "Could not attain list of bridges.")
	fmt.Println("List of bridges: ")
	for _, element := range list {
		fmt.Println(element)
	}
	fmt.Println("End list of bridges.")
}

func (s *BridgeTests) TestRemoveChannel() {
	err := DefaultClient.RemoveChannel("testBridge", "Chan1")
	s.Nil(err, "Failure removing Chan1 from bridge.")
	fmt.Println("Chan1 should be removed from bridge. Sleeping for 3 seconds.")
	time.Sleep(3 * time.Second)
}

func (s *BridgeTests) TestBridgeRecordAndPlay() {
	req := RecordRequest{
		Name:               "name",
		Format:             "wav",
		Beep:               true,
		IfExists:           "overwrite",
		MaxDurationSeconds: 4,
	}
	rec, err := DefaultClient.RecordBridge("testBridge", req)
	fmt.Println("Allowing 5 seconds for liveRecording")
	time.Sleep(5 * time.Second)
	s.Nil(err, "Could not start live recording 'name'")
	s.Equal(rec.Name, "name", "Name retrieved does not match for live recording")
	s.Equal(rec.Format, "wav", "Format retrieved does not match for live recording")

	//Playback the recording.
	req2 := PlayMediaRequest{
		Media: "recording:name",
	}
	_, err = DefaultClient.PlayToBridge("testBridge", req2)
	fmt.Println("Waiting 5 seconds to play back recording")
	time.Sleep(5 * time.Second)
	s.Nil(err, "Couldn't play recording back to 'testBridge'")
}

func TestBridgeSuite(t *testing.T) {
	suite.Run(t, new(BridgeTests))
}
