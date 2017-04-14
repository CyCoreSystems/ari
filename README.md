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
