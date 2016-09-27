package mock

//go:generate mockgen -package mock -destination application.go github.com/CyCoreSystems/ari Application
//go:generate mockgen -package mock -destination asterisk.go github.com/CyCoreSystems/ari Asterisk
//go:generate mockgen -package mock -destination variables.go github.com/CyCoreSystems/ari Variables
//go:generate mockgen -package mock -destination bridge.go github.com/CyCoreSystems/ari Bridge
//go:generate mockgen -package mock -destination channel.go github.com/CyCoreSystems/ari Channel
//go:generate mockgen -package mock -destination device.go github.com/CyCoreSystems/ari DeviceState
//go:generate mockgen -package mock -destination playback.go github.com/CyCoreSystems/ari Playback
//go:generate mockgen -package mock -destination mailbox.go github.com/CyCoreSystems/ari Mailbox
//go:generate mockgen -package mock -destination sound.go github.com/CyCoreSystems/ari Sound
//go:generate mockgen -package mock -destination liveRecording.go github.com/CyCoreSystems/ari LiveRecording
//go:generate mockgen -package mock -destination storedRecording.go github.com/CyCoreSystems/ari StoredRecording
//go:generate mockgen -package mock -destination logging.go github.com/CyCoreSystems/ari Logging
//go:generate mockgen -package mock -destination subscription.go github.com/CyCoreSystems/ari Subscription
//go:generate mockgen -package mock -destination bus.go github.com/CyCoreSystems/ari Bus
//go:generate mockgen -package mock -destination audio_player.go github.com/CyCoreSystems/ari/ext/audio Player
