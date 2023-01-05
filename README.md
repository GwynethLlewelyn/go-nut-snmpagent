# go-nut-snmpagent

This very rudimentary package attempts to create a single executable that exposes the NUT variables from a connected UPS as SNMP; essentially, what [`nut-snmpagent`](https://github.com/luizluca/nut-snmpagent) does, without Ruby or any other interpreted language (e.g. scripts et. al.) — but using the AgentX protocol instead.

## Prerequisites

-   [Network UPS Tools](https://networkupstools.org/) (NUT) installed, configured and running
-   NUT-supported UPS connected to system and fully operational
-   SNMP daemon (e.g. [NET-SNMP 5.9](http://www.net-snmp.org/running) running on the same server, supporting the [AgentX protocol](https://www.rfc-editor.org/rfc/rfc2741)
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

1. On your `/etc/snmp/snmp.conf` just make sure you have something like this:

```
com2sec         notConfigUser  default       public
group           notConfigGroup v1            notConfigUser
group           notConfigGroup v2c           notConfigUser
access          notConfigGroup ""            any  noauth exact systemview none none

view            systemview     included      .1
master  agentx
agentxperms 770 770 daemon users

```

2. Restart the SNMP server daemon (e.g. `sudo systemctl restart snmpd` on most `systemd`-based Un\*xes).

3. Launch `$ ./go-nut-snmpagent`. If all goes well, it should connect to the SNMP server _and_ the NUT server, and start transcoding the information it gets from NUT into SNMP, automagically.

4. Then try to do a `snmpwalk -v 2c -c public localhost .1.3.6.1.4.1.318` to see if you can get all the data via SNMP!

## Final steps

You might wish to add `go-nut-snmpagent` to `systemd` (on Linux, including the Raspberry Pi, Synology NAS, and many other embedded systems) or `launchd` (on macOS), or, if you're so bold, launch it as a service under Windows. Appropriate defaults are found on the `/scripts` directory (Windows instructions not included, since I have no clue how _that_ works).

**TODO:** add more scripts for older methods (e.g. `/etc/init.d`/`/etc/rc.d` or similar Jurassic setups).

## Direct dependencies

See the top of the `go-nut-snmpagent.go` file (or `go.mod`).

Direct compilation tested so far on macOS Bug Sur (Intel amd64), RaspberryOS (Debian) on a Rasperry Pi Zero W 2 (ARM aarch64), and on a Synology NAS DS218play (ARM aarch64).

## Contributions

All contributions/submare most welcome! Just fork the project on Github and submit a PR. If I get more than one contributor, I'll add an automated task to credit you here :-)

Special thanks to [@luizluca](https://github.com/luizluca/) for his outstanding work of creating [`nut-snmpagent`](https://github.com/luizluca/nut-snmpagent) (using a lot of tools, some Ruby and some shell scripting), on which some of this code is based on. @luizluca actually uses an older mechanism to communicate with the SNMP server daemon; we use AgentX. 

And, of course, thanks to all those nice people out there writing Go libraries to read INI files (@unknwon), access NUT, and talk SNMP, all with native Go libraries!

## Future development...?

Not likely, but I'd love to do a standalone SNMP version, possibly even bypassing NUT. NUT is great, don't get me wrong, and has been well-tested for as long as I can remember, but on embedded systems, the less processes spawned, the merrier!

[![Keybase](https://img.shields.io/keybase/pgp/gwynethllewelyn)](https://keybase.io/gwynethllewelyn)
