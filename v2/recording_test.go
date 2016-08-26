package ari

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type RecordingTests struct {
	suite.Suite
	cancel context.CancelFunc
}

func (s *RecordingTests) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	DefaultClient.Listen(ctx)
	time.Sleep(1 * time.Second)
}

func (s *RecordingTests) TearDownSuite() {
	s.cancel()
}

func (s *RecordingTests) TestLiveRecordingFunctions() {
	req2 := OriginateRequest{
		Endpoint:  "PJSIP/101",
		App:       "default",
		AppArgs:   "",
		ChannelId: "Recorder",
	}

	_, err := DefaultClient.CreateChannel(req2)
	s.Nil(err, "Channel by Application not made")

	<-DefaultClient.Bus.Once(context.TODO(), "StasisStart")

	err = DefaultClient.AnswerChannel("MyApp")

	req := RecordRequest{
		Name:     "record",
		Format:   "wav",
		Beep:     true,
		IfExists: "overwrite",
	}
	rec, err := DefaultClient.RecordChannel("Recorder", &req)
	fmt.Println("LiveRecording has begun for recording test")
	s.Nil(err, "Could not start live recording 'record'")
	fmt.Println("Waiting 3 seconds before pausing.")
	time.Sleep(3 * time.Second)

	err = DefaultClient.PauseLiveRecording("record")
	s.Nil(err, "Could not pause recording!")
	fmt.Println("Pausing for 2 seconds")
	time.Sleep(2 * time.Second)

	err = DefaultClient.ResumeLiveRecording("record")
	s.Nil(err, "Could not resume recording!")
	fmt.Println("Resuming recording, waiting 2 seconds before muting.")
	time.Sleep(2 * time.Second)

	err = DefaultClient.MuteLiveRecording("record")
	s.Nil(err, "Could not mute recording!")
	fmt.Println("Muting for 2 seconds")
	time.Sleep(2 * time.Second)

	err = DefaultClient.UnmuteLiveRecording("record")
	s.Nil(err, "Could not un-mute recording!")
	fmt.Println("Unmuted recording, waiting 2 seconds before ending.")

	rec, err = DefaultClient.GetLiveRecording("record")
	s.Nil(err, "Could not retrieve chosen recording!")
	fmt.Println("retrieving current test recording: ", rec)
	time.Sleep(2 * time.Second)

	err = DefaultClient.StopLiveRecording("record")
	s.Nil(err, "Could not stop live recording.")
	fmt.Println("Live recording stopped.")

	fmt.Println("Replaying the live recording...")
	req3 := PlayMediaRequest{
		Media: "recording:record",
	}
	_, err = DefaultClient.PlayToChannel("Recorder", req3)
	fmt.Println("Waiting 5 seconds to play back recording")
	time.Sleep(5 * time.Second)
	s.Nil(err, "Couldn't play recording back to 'Recorder'")

	fmt.Println("Starting new recording to be deleted.")
	req4 := RecordRequest{
		Name:     "record2",
		Format:   "wav",
		Beep:     true,
		IfExists: "overwrite",
	}
	rec, err = DefaultClient.RecordChannel("Recorder", &req4)
	s.Nil(err, "Could not start live recording 'record2'")
	if err == nil {
		fmt.Println("LiveRecording has begun for scrap test")
		fmt.Println("Waiting 2 seconds before scrapping.")
		time.Sleep(2 * time.Second)

		fmt.Println("Scrapping 'record2'")
		err = DefaultClient.ScrapLiveRecording("record2")
		s.Nil(err, "Could not scrap 'record2'")
		if err == nil {
			fmt.Println("'record2' now scrapped.")
		}
	}
}

func (s *RecordingTests) TestStoredRecordingFunctions() {
	//Added one second breaks between each test so that none happen concurrently. Otherwise the IO bogs down and we get errors.
	fmt.Println("Retrieving previous recording 'record'")
	rec, err := DefaultClient.GetStoredRecording("record")
	s.Nil(err, "Couldn't retrieve stored recording.")
	if err == nil {
		fmt.Println("Recording: ", rec)
	}
	time.Sleep(1 * time.Second)

	fmt.Println("Copying recording to 'record3'")
	rec, err = DefaultClient.CopyStoredRecording("record", "record3")
	s.Nil(err, "Couldn't copy stored recording.")
	if err == nil {
		fmt.Println("Copied recording: ", rec)
	}
	time.Sleep(1 * time.Second)

	fmt.Println("Deleting recording 'record'")
	err = DefaultClient.DeleteStoredRecording("record")
	s.Nil(err, "Couldn't delete stored recording.")
	time.Sleep(1 * time.Second)

	list, err := DefaultClient.ListStoredRecordings()
	s.Nil(err, "Couldn't list stored recordings.")
	if err == nil {
		fmt.Println("Listing recordings. 'record3' should exist. 'record2' and 'record' should not.")
		for _, element := range list {
			fmt.Println("Element: ", element)
		}
	}
	time.Sleep(1 * time.Second)

	fmt.Println("Deleting recording 'record3' for future tests.")
	err = DefaultClient.DeleteStoredRecording("record3")
	s.Nil(err, "Couldn't delete stored recording.")
}

func TestRecordingSuite(t *testing.T) {
	suite.Run(t, new(RecordingTests))
}
