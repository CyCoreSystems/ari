# ari - Golang Asterisk Rest Interface (ARI) library
[![Build Status](https://travis-ci.org/CyCoreSystems/ari.png)](https://travis-ci.org/CyCoreSystems/ari) [![](https://godoc.org/github.com/CyCoreSystems/ari?status.svg)](http://godoc.org/github.com/CyCoreSystems/ari)

This is a go-based ARI library.  It also includes some common convenience wrappers for various tasks, which can be found in /ext.

This library maintains semver, and APIs between major releases **do** change.
Therefore, always use a vendoring tool which supports semver, such as `glide` or
`dep` or use the `gopkg.in` aliasing, such as `gopkg.in/CyCoreSystems/ari.v3`.

The `v3` branch is the most well-tested branch, while `v4` fixes a number of
shortcomings of `v3`, particularly for interoperating with proxies clients.

There is also a NATS-based `ari-proxy` which is designed to work with this
client library.  It can be found at
[CyCoreSystems/ari-proxy](https://github.com/CyCoreSystems/ari-proxy).


# Resource Keys

In order to facilitate the construction of ARI systems across many Asterisk
instances, in version 4, we introduce the concept of Resource Keys.  Previous
versions expected a simple ID (string) field for the identification of a
resource to ARI.  This reflects how ARI itself operates.  However, for systems
with multiple Asterisk instances, more metadata is necessary in order to
properly address a resource.  Specifically, we need to know the Asterisk node.
There is also the concept of a Dialog, which is, generically, a named
transaction with logically-bound resources.  This Key includes all of these.

```go
package ari

// Key identifies a unique resource in the system
type Key struct {
   // Kind indicates the type of resource the Key points to.  e.g., "channel",
   // "bridge", etc.
   Kind string

   // ID indicates the unique identifier of the resource
   ID string

   // Node indicates the unique identifier of the Asterisk node on which the
   // resource exists or will be created
   Node string

   // Dialog indicates a named scope of the resource, for receiving events
   Dialog string
}
```
At a basic level, when the specific Asterisk ID is not needed, a key can consist
of a simple ID string:

```go
  key := ari.NewKey(ari.KeyChannel, "myID")
```

For more interesting systems, however, we can declare the Node ID:

```go
  key := ari.NewKey(ari.KeyBridge, "myID", ari.WithNode("00:01:02:30:40:50"))
```

We can also bind a dialog:

```go
  key := ari.NewKey(ari.KeyChannel, "myID",
   ari.WithNode("00:01:02:30:40:50"),
   ari.WithDialog("privateConversation"))
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

