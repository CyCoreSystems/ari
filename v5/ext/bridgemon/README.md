# bridgemon

[![](https://godoc.org/github.com/CyCoreSystems/ari?status.svg)](http://godoc.org/github.com/CyCoreSystems/ari)

Bridge Monitor provides a simple tool to monitor and cache a bridge's data for
easy, efficient access by other routines.  It is safe for multi-threaded use and
can be closed manually or whenever the bridge is destroyed.

It is created by passing a bridge handle in.  The bridge should already exist
for this to be operational, and initial data is loaded when the monitor is
created.

There are two method for consuming the data.  `Data()` provides arbitrary access
to the cached bridge data while `Watch()` provides a channel over which the
bridge data will be sent whenever updates are made.

