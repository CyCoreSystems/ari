package ari

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type BridgeTests struct {
	suite.Suite
	cancel context.CancelFunc
}

func (s *BridgeTests) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	DefaultClient.Listen(ctx)
	time.Sleep(1 * time.Second)
}

func (s *BridgeTests) TearDownSuite() {
	s.cancel()
}

func (s *BridgeTests) TestBridgeCreate() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Creation of empty bridge.
	bridge, err := DefaultClient.CreateBridge(CreateBridgeRequest{
		Id:   "testBridge",
		Type: "mixing",
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

	<-DefaultClient.Bus.Once(ctx, "StasisStart")

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

	<-DefaultClient.Bus.Once(ctx, "StasisStart")

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
	rec, err := DefaultClient.RecordBridge("testBridge", &req)
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

type BridgeTestsSplit struct {
	suite.Suite
	cancel context.CancelFunc
}

func (s *BridgeTestsSplit) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	DefaultClient.Listen(ctx)
	time.Sleep(1 * time.Second)
}

func (s *BridgeTestsSplit) TearDownSuite() {
	s.cancel()
}

func (s *BridgeTestsSplit) TestBridgeOtherApp() {
	// Create a channel with the default client
	ch1, err := DefaultClient.NewChannel("PJSIP/101", nil, nil)
	s.Nil(err, "Failed to create first channel")

	// Wait for answer
	<-DefaultClient.Bus.Once(context.TODO(), "StasisStart")

	// Create a bridge with the default client
	br, err := DefaultClient.NewBridge()
	s.Nil(err, "Failed to create bridge")

	// Create a new client
	secondCtx, cancel := context.WithCancel(context.Background())
	nc := NewClient(nil)
	nc.Listen(secondCtx)
	defer cancel()

	// Create channel on second client
	ch2, err := nc.NewChannel("PJSIP/102", nil, nil)
	s.Nil(err, "Failed to create second channel")

	// Wait for answer
	<-nc.Bus.Once(secondCtx, "StasisStart")

	// Add ch1 to bridge
	err = br.Add(ch1.Id)
	s.Nil(err, "Failed to add ch1 to bridge")

	// Add ch2 to the bridge
	err = nc.AddChannel(br.Id, AddChannelRequest{ChannelId: ch2.Id})
	s.Nil(err, "Failed to add ch2 to bridge")

	fmt.Println("Waiting 1 second for achievement of serenity")
	time.Sleep(1 * time.Second)

	// Tear down bridge by the second client
	// (the one which did NOT create it)
	//err = br.Delete()
	err = nc.BridgeDelete(br.Id)
	s.Nil(err, "Failed to delete bridge")

	// Hang up ch1
	err = ch1.Hangup()
	s.Nil(err, "Failed to hang up ch1")

	// Hang up ch2
	err = ch2.Hangup()
	s.Nil(err, "Failed to hang up ch2")
}

func TestBridgeSuite(t *testing.T) {
	suite.Run(t, new(BridgeTests))
}

func TestBridgeSplitSuite(t *testing.T) {
	suite.Run(t, new(BridgeTestsSplit))
}
