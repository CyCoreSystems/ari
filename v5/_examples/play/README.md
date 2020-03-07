# Examples - Play

This example ARI application listens for calls coming into the Stasis app
"test" and then answers the line, plays a sound to the caller, and hangs up.

## Asterisk dialplan

An example dialplan for extension `100` would be something like this:

```asterisk
exten = 100,1,Stasis("test")
```

## Asterisk ARI configuration

In order for the example application to connect with Asterisk, a few settings
must be enabled.

`http.conf` settings:

```ini
[general]
enabled=yes
bindaddr=127.0.0.1
bindport=8088
```

`ari.conf` settings:

```ini
[general]
enabled = yes
allowed_origins = * ; tighten this down later

[admin]
type = user
read_only = no
password_format = crypt
password = $6$/ejLut/kmjN6E5.g$tXEeth2SQoVYSs0AG0wWIoB3XRJEqK9vm0JGxQHU7Q/IIR/Ln5Zho40fcPUv1n8jvOJWYMJg0/4fLdJpSB2du1
```

**NOTE**: to obtain an encrypted password, you can use the `ari mkpassword`
command from Asterisk.  In this case, the following was done:

```
# asterisk -rx "ari mkpasswd admin"
```

## Runtime

Now, execute the example application, and it should connect to Asterisk and
register the "test" ARI application.

You may verify that it is registered by running the `ari show apps` Asterisk
command:

```
# asterisk -rx "ari show apps"
```

## Call

Now, make a call into your Asterisk box to extension 100, and you should hear
the playback (assuming you have the Asterisk extra sounds installed).  You
should also see the call come in on your application's log.


