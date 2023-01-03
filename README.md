# go-nut-snmpagent

This very rudimentary package attempts to create a single executable that exposes the NUT variables from a connected UPS as SNMP; essentially, what nut-snmpagent does, without Ruby or any other interpreted language (e.g. scripts et. al.).

## Prerequisites

-   [Network UPS Tools](https://networkupstools.org/) (NUT) installed, configured and running
-   NUT-supported UPS connected to system and fully operational
-   SNMP daemon (e.g. [NET-SNMP 5.9](http://www.net-snmp.org/running) running on the same server
-   A fairly recent Go compiler (requires at least modules to work)

## Installing go-nut-snmpagent

```bash
$ git clone github.com/GwynethLlewelyn/go-nut-snmpagent
$ cd go-nut-snmpagent
$ go build
```

This should leave a binary named `go-nut-snmpagent` on the same directory. By default, Go compiles binaries with debug data, so you might wish to run `strip go-nut-snmpagent` to make the binary smaller.

Then create a `config.ini` to override the defaults on `config.main.ini`, if you wish (e.g. authentication). You can see what options exist on that file, and set those that you wish or need.

## NUT configuration

You don't need to change anything on NUT. `go-nut-snmpagent` will simply connect to it and retrieve the data from the first UPS it can find there.

Of course, if you need _authentication_, then you _might_ need to configure NUT to allow the agent to connect to it (optional, especially if both are running on the same system, in the same local network).

## SNMP configuration

On your `/etc/snmp/snmp.conf` just put the following line:

```
pass_persist .1.3.6.1.4.1.26376.99 /path/where/you/compiled/go-nut-snmpagent
```

Restart the SNMP server daemon (e.g. `sudo systemctl restart snmpd` on most `systemd`-based Un\*xes).

Then try to do a `snmpwalk -v 2c public .1.3.6.1.4.1.263` to see if you can connect to the UPS and extract its data/parameters!

## Direct dependencies

See the top of the `go-nut-snmpagent.go` file (or `go.mod`).

Direct compilation tested so far on macOS Bug Sur (Intel amd64), RaspberryOS (Debian) on a Rasperry Pi Zero W 2 (ARM aarch64), and on a Synology NAS DS218play (ARM aarch64).

## Contributions

All are most welcome! Just fork the project on Github and submit a PR. If I get more than one contributor, I'll add an automated task to credit you here :-)

Special thanks to [@luizluca](https://github.com/luizluca/) for his outstanding work of creating [`nut-snmpagent`](https://github.com/luizluca/nut-snmpagent) (using a lot of tools, some Ruby and some shell scripting), on which this bit of code is based on.

And, of course, thanks to all those nice people out there writing Go libraries to read INI files (@unknwon), access NUT, and talk SNMP, all with native Go libraries!

## Future development...?

Not likely, but I'd love to do a standalone SNMP version, possibly even bypassing NUT. NUT is great, don't get me wrong, and has been well-tested for as long as I can remember, but on embedded systems, the less processes spawned, the merrier!

[![Keybase](https://img.shields.io/keybase/pgp/gwynethllewelyn)](https://keybase.io/gwynethllewelyn)
