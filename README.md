# ari - Golang Asterisk Rest Interface (ARI) library
[![Build Status](https://travis-ci.org/CyCoreSystems/ari.png)](https://travis-ci.org/CyCoreSystems/ari) [![](https://godoc.org/github.com/CyCoreSystems/ari?status.svg)](http://godoc.org/github.com/CyCoreSystems/ari)

This is a go-based ARI library.  It also includes some common convenience wrappers for various tasks, which can be found in /ext.

# Getting started

This library maintains semver, and APIs between major releases **do** change.
Therefore, always use a vendoring tool which supports semver, such as [glide](http://glide.sh/) or
[dep](https://github.com/golang/dep).

Version `4.x.x` is the current version.  It offers a number of
new features focused on facilitating ARI across large clusters and simplifies
the API.

There is also a NATS-based `ari-proxy` which is designed to work with this
client library.  It can be found at
[CyCoreSystems/ari-proxy](https://github.com/CyCoreSystems/ari-proxy).

# Cloud-ready

All configuration options for the client are able to be sourced by environment
variable, making it easy to build applications without configuration files.
Moreover, the default connection to Asterisk is set to `localhost` on port 8088,
which should run on Kubernetes deployments without configuration.

The available environment variables (and defaults) are:

  - `ARI_APPLICATION` (*randomly-generated UUID*)
  - `ARI_URL` (`http://localhost:8088/ari`)
  - `ARI_WSURL` (`ws://localhost:8088/ari/events`)
  - `ARI_WSORIGIN` (`http://localhost/`)
  - `ARI_USERNAME` (*none*)
  - `ARI_PASSWORD` (*none*)

When using the `ari-proxy`, the process is even easier.

# Resource Keys

In order to facilitate the construction of ARI systems across many Asterisk
instances, in version 4, we introduce the concept of Resource Keys.  Previous
versions expected a simple ID (string) field for the identification of a
resource to ARI.  This reflects how ARI itself operates.  However, for systems
with multiple Asterisk instances, more metadata is necessary in order to
properly address a resource.  Specifically, we need to know the Asterisk node.
There is also the concept of a Dialog, which offers an orthogonal logical
grouping of events which transcends nodes and applications.  This is not
meaningful in the native client, but other transports, such as the ARI proxy,
may make use of this for alternative routing of events.

This Key includes all of these data.

```go
package ari

// Key identifies a unique resource in the system
type Key struct {
   // Kind indicates the type of resource the Key points to.  e.g., "channel",
   // "bridge", etc.
   Kind string   `json:"kind"`

   // ID indicates the unique identifier of the resource
   ID string `json:"id"`

   // Node indicates the unique identifier of the Asterisk node on which the
   // resource exists or will be created
   Node string `json:"node,omitempty"`

   // Dialog indicates a named scope of the resource, for receiving events
   Dialog string `json:"dialog,omitempty"`
}
```
At a basic level, when the specific Asterisk ID is not needed, a key can consist
of a simple ID string:

```go
  key := ari.NewKey(ari.ChannelKey, "myID")
```

For more interesting systems, however, we can declare the Node ID:

```go
  key := ari.NewKey(ari.BridgeKey, "myID", ari.WithNode("00:01:02:30:40:50"))
```

We can also bind a dialog:

```go
  key := ari.NewKey(ari.ChannelKey, "myID",
   ari.WithNode("00:01:02:30:40:50"),
   ari.WithDialog("privateConversation"))
```

We can also create a new key from an existing key.  This allows us to easily
copy the location information from the original key to a new key of a different
resource.  The location information is everything (including the Dialog) except
for the key Kind and ID.

```go
  brKey := key.New(ari.BridgeKey, "myBridgeID")

```

All ARI operations which accepted an ID for an operator now expect an `*ari.Key`
instead.  In many cases, this can be easily back-ported by wrapping IDs with
`ari.NewKey("channel", id)`.

# Staging resources

A common issue for ARI resources is making sure a subscription exists before
events for that resource are sent.  Otherwise, important events which occur too
quickly can become lost.  This results in a chicken-and-egg problem for
subscriptions.

In order to address this common issue, resource handles creation operations now
offer a `StageXXXX` variant, which returns the handle for the resource without
actually creating the resource.  Once all of the subscriptions are bound to this
handle, the caller may call `resource.Exec()` in order to create the resource in
Asterisk.

```go
   h := NewChannelHandle(key, c, nil)

   // Stage a playback
   pb, err := h.StagePlay("myPlaybackID", "sound:tt-monkeys")
   if err != nil {
      return err
   }
   
   // Add a subscription to the staged playback
   startSub := pb.Subscribe(EventTypes.PlaybackStarted)
   defer startSub.Cancel()

   // Execute the staged playback
   pb.Exec()

   // Wait for something to happen
   select {
      case <-time.After(time.Second):
        fmt.Println("timeout waiting for playback to start")
        return errors.New("timeout")
      case <-startSub.Events():
        fmt.Println("playback started")
   }
```

# Play

Playback of media and waiting for (DTMF) responses therefrom is an incredibly
common task in telephony.  ARI provides many tools to perform these types of
actions, but the asynchronous nature of the interface makes it fairly tedious to
build these kinds of things.

In `ext/play`, there resides a tool for executing many common tasks surrounding
media playback and response sequences.  The core function, `play.Play()`
plays, in sequence, a series of audio media URIs.  It can be extended to expect
and (optionally) wait for a DTMF response by supplying it with a Match function.
There is a small convenience wrapper `play.Prompt()` which sets some common
defaults for playbacks which expect a response.

The execution of a `Play` is configured by any number of option functions, which
supply structured modifiers for the behaviour of the playback.  You can even
supply your own Match function for highly-customized matching.


