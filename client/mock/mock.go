package mock

//go:generate mockgen -package mock -destination application.go github.com/CyCoreSystems/ari Application
//go:generate mockgen -package mock -destination asterisk.go github.com/CyCoreSystems/ari Asterisk
//go:generate mockgen -package mock -destination variables.go github.com/CyCoreSystems/ari Variables
//go:generate mockgen -package mock -destination bridge.go github.com/CyCoreSystems/ari Bridge
