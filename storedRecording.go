package ari

// StoredRecording represents a communication path interacting with an Asterisk
// server for stored recording resources
type StoredRecording interface {

	// List lists the recordings
	List() ([]*StoredRecordingHandle, error)

	// Get gets the Recording by type
	Get(name string) *StoredRecordingHandle

	// data gets the data for the stored recording
	Data(name string) (StoredRecordingData, error)

	// Copy copies the recording to the destination name
	Copy(name string, dest string) (*StoredRecordingHandle, error)

	// Delete deletes the recording
	Delete(name string) error
}

// StoredRecordingData is the data for a stored recording
type StoredRecordingData struct {
	Format string `json:"format"`
	Name   string `json:"name"`
}

// ID returns the identifier for the stored recording.
func (d StoredRecordingData) ID() string {
	return d.Name //TODO: does the identifier include the Format and Name?
}

// NewStoredRecordingHandle creates a new stored recording handle
func NewStoredRecordingHandle(name string, s StoredRecording) *StoredRecordingHandle {
	return &StoredRecordingHandle{
		name: name,
		s:    s,
	}
}

// A StoredRecordingHandle is a reference to a stored recording that can be operated on
type StoredRecordingHandle struct {
	name string
	s    StoredRecording
}

// ID returns the identifier for the stored recording
func (s *StoredRecordingHandle) ID() string {
	return s.name
}

// Data gets the data for the stored recording
func (s *StoredRecordingHandle) Data() (d StoredRecordingData, err error) {
	d, err = s.s.Data(s.name)
	return
}

// Copy copies the stored recording
func (s *StoredRecordingHandle) Copy(dest string) (h *StoredRecordingHandle, err error) {
	h, err = s.s.Copy(s.name, dest)
	return
}

// Delete deletes the recording
func (s *StoredRecordingHandle) Delete() (err error) {
	err = s.s.Delete(s.name)
	return
}
