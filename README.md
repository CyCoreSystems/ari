# ari - Golang Asterisk Rest Interface (ARI) library
[![Build Status](https://travis-ci.org/CyCoreSystems/ari.png)](https://travis-ci.org/CyCoreSystems/ari) [![](https://godoc.org/github.com/CyCoreSystems/ari?status.svg)](http://godoc.org/github.com/CyCoreSystems/ari)

This library allows you to easily access ARI in go applications.  The Asterisk Rest Interface (https://wiki.asterisk.org/wiki/pages/viewpage.action?pageId=29395573) is an asynchronous API which allows you to access basic Asterisk objects for custom communications applications.  

This project also includes some convenience wrappers for various tasks, found in /ext.  These include go-idiomatic utilities for playing audio, IVRs, recordings, and other tasks which are tricky to coordinate nicely in ARI alone.  

# Getting started

This library maintains semver, and APIs between major releases **do** change.
We use `GO111MODULE`, so Go version 1.11 or later is required.

Version `5.x.x` is the current version.

There is also a NATS-based `ari-proxy` which is designed to work with this
client library.  It can be found at
[CyCoreSystems/ari-proxy](https://github.com/CyCoreSystems/ari-proxy).

Install with: 
```sh 
go get github.com/CyCoreSystems/ari
```

# Features

## Cloud-ready

All configuration options for the client can be sourced by environment
variable, making it easy to build applications without configuration files.
The default connection to Asterisk is set to `localhost` on port 8088,
which should run on Kubernetes deployments without configuration.

The available environment variables (and defaults) are:

  - `ARI_APPLICATION` (*randomly-generated ID*)
  - `ARI_URL` (`http://localhost:8088/ari`)
  - `ARI_WSURL` (`ws://localhost:8088/ari/events`)
  - `ARI_WSORIGIN` (`http://localhost/`)
  - `ARI_USERNAME` (*none*)
  - `ARI_PASSWORD` (*none*)

If using `ari-proxy`, the process is even easier.

## Resource Keys

In order to facilitate the construction of ARI systems across many Asterisk
instances, in version 4, we introduce the concept of Resource Keys.  Previous
versions expected a simple ID (string) field for the identification of a
resource to ARI.  This reflects how ARI itself operates.  However, for systems
with multiple Asterisk instances, more metadata is necessary in order to
properly address a resource.  Specifically, we need to know the Asterisk node.
There is also the concept of a Dialog, which offers an orthogonal logical
grouping of events which transcends nodes and applications.  This is not
meaningful in the native client, but other transports, such as the ARI proxy,
may make use of this for alternative routing of events. This Key includes all of these data.

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

## Handles

Handles for all of the major entity types are available, which bundle in the
tracking of resources with their manipulations.  Every handle, at a minimum,
internally tracks the resource's cluster-unique Key and the ARI client
connection through which the entity is being interacted.  Using a handle
generally results in less and more readable code.

For instance, instead of calling:

```go
ariClient.Channel().Hangup(channelKey, "normal")
```

you could just call `Hangup()` on the handle:

```go
h.Hangup()
```

While the lower level direct calls have maintained fairly strict semantics to
match the formal ARI APIs, the handles frequently provide higher-level, simpler
operations.  Moreover, most of the extensions (see below) make use of handles.

In general, when operating on longer lifetime entities (such as channels and
bridges), it is easier to use handles wherever you can rather than tracking Keys
and clients discretely.

Obtaining a Handle from a Key is very simple; just call the `Get()` operation on
the resource interface appropriate to the key.  The `Get()` operation is a
local-only operation which does not interact with the Asterisk or ARI proxy at
all, and it is thus quite efficient.

```go
h := ariClient.Channel().Get(channelKey)
```

## Staging resources

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

## Extensions

There are a number of extensions which wrap the lower-level operations in
higher-level ones, making it easier to perform many common tasks.


### AudioURI [![](https://godoc.org/github.com/CyCoreSystems/ari?status.svg)](http://godoc.org/github.com/CyCoreSystems/ari/ext/audiouri)

Constructing Asterisk audio playback URIs can be a bit tedious, particularly for handling
certain edge cases in digits and for constructing dates.

The `audiouri` package provides a number of routines to make the construction of
these URIs simpler.

### Bridgemon [![](https://godoc.org/github.com/CyCoreSystems/ari?status.svg)](http://godoc.org/github.com/CyCoreSystems/ari/ext/bridgemon)

Monitoring a bridge for events and data updates is not difficult, but it
involves a lot of code and often makes several wasteful calls to obtain bridge
data, particularly when accessing it on large bridges.

Bridgemon provides a cache and proxy for the bridge data and bridge events so
that a user can simply `Watch()` for changes in the bridge state and efficiently
retrieve the updated data without multiple requests.

It also shuts itself down automatically when the bridge it is monitoring is
destroyed.

### Play [![](https://godoc.org/github.com/CyCoreSystems/ari?status.svg)](http://godoc.org/github.com/CyCoreSystems/ari/ext/play)

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

### Record [![](https://godoc.org/github.com/CyCoreSystems/ari?status.svg)](http://godoc.org/github.com/CyCoreSystems/ari/ext/record)

Making recordings is another complicated but common task for ARI applications.
The `ext/record`, we provide a simple wrapper which facilitates many common
recording-related operations inside a single recording Session wrapper.

Features include:

  - record with or without a beep at the start
  - listen for various termination types: hangup, dtmf, silence, timeout
  - review, scrap, and save recordings upon completion
  - retrieve the playback URI for the recording

# Documentation and Examples

Go documentation is available at https://godoc.org/github.com/CyCoreSystems/ari

Examples for helloworld, play, script, bridge, and record are available.  Set your environment variables as described above (at minimum, `ARI_USERNAME` and `ARI_PASSWORD`) and run:

```sh
cd /_examples/helloworld
go run ./main.go
```

Other examples:

 - `stasisStart` demonstrates a simple click-to-call announcer system
 - `stasisStart-nats` demonstrates the same click-to-call using the NATS-based
   ARI proxy
 - `bridge` demonstrates a simple conference bridge
 - `play` demonstrates the use of the `ext/play` extension
 - `record` demonstrates the use of the `ext/record` extension

The files in `_ext/infra` demonstrates the minimum necessary changes to the
Asterisk configuration to enable the operation of ARI.


# Tests

Run `go test` to verify 

# Contributing

Contributions welcomed. Changes with tests and descriptive commit messages will get priority handling.  

# License

Licensed under the Apache License, Version 2.0 
