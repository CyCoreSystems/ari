package mock

//go:generate mockgen -package mock -destination application.go github.com/CyCoreSystems/ari Application
//go:generate mockgen -package mock -destination asterisk.go github.com/CyCoreSystems/ari Asterisk
//go:generate mockgen -package mock -destination variables.go github.com/CyCoreSystems/ari Variables
//go:generate mockgen -package mock -destination bridge.go github.com/CyCoreSystems/ari Bridge
//go:generate mockgen -package mock -destination channel.go github.com/CyCoreSystems/ari Channel
//go:generate mockgen -package mock -destination device.go github.com/CyCoreSystems/ari DeviceState
//go:generate mockgen -package mock -destination playback.go github.com/CyCoreSystems/ari Playback
//go:generate mockgen -package mock -destination subscription.go github.com/CyCoreSystems/ari Subscription
