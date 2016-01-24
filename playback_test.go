package ari

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type PlaybackTests struct {
	suite.Suite
	list   []Channel
	cancel context.CancelFunc
}

func (s *PlaybackTests) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	DefaultClient.Listen(ctx)
	s.cancel = cancel
	time.Sleep(1 * time.Second)
}

func (s *PlaybackTests) TearDownSuite() {
	s.cancel()
}

func (s *PlaybackTests) TestGetPlaybackDetails() {
	//This bit sets up a channel.

	req := OriginateRequest{
		Endpoint:  "PJSIP/101",
		App:       "default",
		AppArgs:   "",
		ChannelId: "MyApp",
	}

	Chan2, err := DefaultClient.CreateChannel(req)
	s.Nil(err, "Channel by Application not made")
	s.NotEmpty(Chan2.Id)

	s.Equal("MyApp", Chan2.Id)

	//Wait until we receive "StasisStart" from our channel, meaning it has been answered (on far end) and may be answered by Asterisk.

	<-DefaultClient.Bus.Once(context.TODO(), "StasisStart")

	err = DefaultClient.AnswerChannel("MyApp")
	s.Nil(err)
	//Channel has been answered and set up properly.

	//Creating a recording of the channel.
	req3 := RecordRequest{
		Name:               "name",
		Format:             "wav",
		Beep:               true,
		IfExists:           "overwrite",
		MaxDurationSeconds: 4,
	}
	rec, err := DefaultClient.RecordChannel("MyApp", req3)
	fmt.Println("Allowing 6 seconds for liveRecording")
	time.Sleep(6 * time.Second)
	s.Nil(err, "Could not start live recording 'name'")
	s.Equal(rec.Name, "name", "Name retrieved does not match for live recording")
	s.Equal(rec.Format, "wav", "Format retrieved does not match for live recording")

	req2 := PlayMediaRequest{
		Media: "recording:name",
	}

	playback, err := DefaultClient.PlayToChannel("MyApp", req2)
	fmt.Println("Returned playback (playtochannel): ", playback)
	fmt.Println("Waiting 8 seconds to play back recording")
	s.Nil(err, "Couldn't play recording back to 'MyApp'")
	//Playing the recording back to the channel, creating a playback session.

	//Trying to retrieve that playback session.
	playback2, err := DefaultClient.GetPlaybackDetails(playback.Id)
	fmt.Printf("%+v", playback2)
	s.Nil(err, "Error retrieving playback.")
	s.Equal(playback.Id, playback2.Id, "Playbacks not equal.")
	fmt.Println("Playback returned: ", playback2)
	time.Sleep(3 * time.Second)

	//Testing the control function using the 'restart' command.
	err = DefaultClient.ControlPlayback(playback.Id, "restart")
	fmt.Println("Restarting playback...")
	time.Sleep(3 * time.Second)
	s.Nil(err, "Couldn't restart playback.")

	//Lastly, let's delete the playback session prematurely.
	fmt.Println("Cutting off playback.")
	err = DefaultClient.StopPlayback(playback.Id)
	s.Nil(err, "Couldn't stop playback.")
}

func TestPlaybackSuite(t *testing.T) {
	suite.Run(t, new(PlaybackTests))
}
