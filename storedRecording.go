package ari

// StoredRecording represents a communication path interacting with an Asterisk
// server for stored recording resources
type StoredRecording interface {

	// List lists the recordings
	List() ([]StoredRecordingHandle, error)

	// Get gets the Recording by type
	Get(name string) StoredRecordingHandle

	// data gets the data for the stored recording
	Data(name string) (*StoredRecordingData, error)

	// Copy copies the recording to the destination name
	Copy(name string, dest string) (StoredRecordingHandle, error)

	// StageCopy creates a `StoredRecordingHandle` with a `Copy` operation staged.
	StageCopy(name string, dest string) StoredRecordingHandle

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

// A StoredRecordingHandle is a reference to a stored recording that can be operated on
type StoredRecordingHandle interface {
	// ID returns the identifier for the stored recording
	ID() string

	// Data gets the data for the stored recording
	Data() (d *StoredRecordingData, err error)

	// Copy copies the stored recording
	Copy(dest string) (h StoredRecordingHandle, err error)

	// StageCopy creates a `StoredRecordingHandle` with a `Copy` operation staged.
	StageCopy(dest string) (h StoredRecordingHandle)

	// Delete deletes the recording
	Delete() (err error)

	// Exec executes any staged operations attached to the handle.
	Exec() (err error)
}
