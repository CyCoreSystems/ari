package ari

import (
	"errors"
	"io"
)

// StoredRecording represents a communication path interacting with an Asterisk
// server for stored recording resources
type StoredRecording interface {
	// List lists the recordings
	List(filter *Key) ([]*Key, error)

	// Get gets the Recording by type
	Get(key *Key) *StoredRecordingHandle

	// data gets the data for the stored recording
	Data(key *Key) (*StoredRecordingData, error)

	// Download retrieves the data for the stored recording
	// IMPORTANT: Don't forget to close the reader when done.
	Download(key *Key) (*StoredRecordingBinaryData, error)

	// Copy copies the recording to the destination name
	//
	// NOTE: because ARI offers no forced-copy, Copy should always return the
	// StoredRecordingHandle of the destination, even if the Copy fails.  Doing so
	// allows the user to Delete the existing StoredRecording before retrying.
	Copy(key *Key, dest string) (*StoredRecordingHandle, error)

	// Delete deletes the recording
	Delete(key *Key) error
}

// StoredRecordingData is the data for a stored recording
type StoredRecordingData struct {
	// Key is the cluster-unique identifier for this stored recording
	Key *Key `json:"key"`

	Format string `json:"format"`
	Name   string `json:"name"`
}

// ID returns the identifier for the stored recording.
func (d StoredRecordingData) ID() string {
	return d.Name // TODO: does the identifier include the Format and Name?
}

// StoredRecordingBinaryData is a binary reader for a stored recording file.
type StoredRecordingBinaryData struct {
	// Key is the cluster-unique identifier for this stored recording file
	Key *Key

	io.ReadCloser

	// ContentType is the MIME type of the stored recording file, e.g. "audio/wav"
	ContentType string
}

// Close closes the ReadCloser for the stored recording binary data.
func (d *StoredRecordingBinaryData) Close() error {
	if d.ReadCloser != nil {
		return d.ReadCloser.Close()
	}
	return nil
}

// Read reads data from the stored recording binary data.
func (d *StoredRecordingBinaryData) Read(p []byte) (n int, err error) {
	if d.ReadCloser == nil {
		return 0, errors.New("read closer is nil")
	}
	return d.ReadCloser.Read(p)
}

// A StoredRecordingHandle is a reference to a stored recording that can be operated on
type StoredRecordingHandle struct {
	key      *Key
	s        StoredRecording
	exec     func(a *StoredRecordingHandle) error
	executed bool
}

// NewStoredRecordingHandle creates a new stored recording handle
func NewStoredRecordingHandle(key *Key, s StoredRecording, exec func(a *StoredRecordingHandle) error) *StoredRecordingHandle {
	return &StoredRecordingHandle{
		key:  key,
		s:    s,
		exec: exec,
	}
}

// ID returns the identifier for the stored recording
func (s *StoredRecordingHandle) ID() string {
	return s.key.ID
}

// Key returns the Key for the stored recording
func (s *StoredRecordingHandle) Key() *Key {
	return s.key
}

// Exec executes any staged operations
func (s *StoredRecordingHandle) Exec() (err error) {
	if !s.executed {
		s.executed = true
		if s.exec != nil {
			err = s.exec(s)
			s.exec = nil
		}
	}

	return
}

// Data gets the data for the stored recording
func (s *StoredRecordingHandle) Data() (*StoredRecordingData, error) {
	return s.s.Data(s.key)
}

// Download retrieves binary reader of the stored recording
func (s *StoredRecordingHandle) Download() (*StoredRecordingBinaryData, error) {
	return s.s.Download(s.key)
}

// Copy copies the stored recording.
//
// NOTE: because ARI offers no forced-copy, this should always return the
// StoredRecordingHandle of the destination, even if the Copy fails.  Doing so
// allows the user to Delete the existing StoredRecording before retrying.
func (s *StoredRecordingHandle) Copy(dest string) (*StoredRecordingHandle, error) {
	return s.s.Copy(s.key, dest)
}

// Delete deletes the recording
func (s *StoredRecordingHandle) Delete() error {
	return s.s.Delete(s.key)
}
