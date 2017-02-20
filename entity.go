package ari

// An Entity is a representation of an event-attached item.
// +gen slice:"SortBy"
type Entity struct {
	// Type is the type of entity
	Type string

	// ID is the unique identifier for the entity
	ID string
}

// EntitiesFromEvent converts an ARI event to a series of entities that are connected
// to the event.
func EntitiesFromEvent(e Event) (ret []Entity) {
	if v, ok := e.(ChannelEvent); ok {
		for _, id := range v.GetChannelIDs() {
			ret = append(ret, Entity{
				Type: "channel",
				ID:   id,
			})
		}
	}

	if ev, ok := e.(BridgeEvent); ok {
		for _, id := range ev.GetBridgeIDs() {
			ret = append(ret, Entity{
				Type: "bridge",
				ID:   id,
			})
		}
	}

	if ev, ok := e.(EndpointEvent); ok {
		for _, id := range ev.GetEndpointIDs() {
			ret = append(ret, Entity{
				Type: "endpoint",
				ID:   id,
			})
		}
	}

	if ev, ok := e.(PlaybackEvent); ok {
		for _, id := range ev.GetPlaybackIDs() {
			ret = append(ret, Entity{
				Type: "playback",
				ID:   id,
			})
		}
	}

	if ev, ok := e.(RecordingEvent); ok {
		for _, id := range ev.GetRecordingIDs() {
			ret = append(ret, Entity{
				Type: "recording",
				ID:   id,
			})
		}
	}

	return
}
