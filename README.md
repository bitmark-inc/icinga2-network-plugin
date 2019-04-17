[![GoDoc](https://godoc.org/github.com/araujobsd/icinga2-network-plugin/plugins?status.svg)](https://godoc.org/github.com/araujobsd/icinga2-network-plugin/)
[![GitHub issues](https://img.shields.io/github/issues/araujobsd/icinga2-network-plugin.svg)](https://github.com/araujobsd/icinga2-network-plugin/issues)
[![GitHub forks](https://img.shields.io/github/forks/araujobsd/icinga2-network-plugin.svg)](https://github.com/araujobsd/icinga2-network-plugin/network)

icinga2-network-plugin
================
This plugin has been created especially for Icinga2, but it is compatible with Nagios 4 too, it was developed using [Go](http://golang.org/).

This plugin uses nvstat software to check Incoming and Outgoing speed and create alerts. 

## Build instructions
1) `make build`

## How to use
Run the plugin.
```
root@status:/usr/home/freebsd # /usr/local/libexec/nagios/bsd_network
[ bsd_network - Version: 0.1 (Marcelo Araujo <araujo__FreeBSD.org>) ]
Options:
	-rw: Incoming Speed Warning (KiB/s)
	-rc: Incoming Speed Critical (KiB/s)
	-tw: Outgoing Speed Warning (KiB/s)
	-tc: Outgoing Speed Critical (KiB/s)
	-i: Interface this should monitor

Note:
The units you specify must be the same units as configured for vnstat(1)

Usage:
./bsd_network -rw=<incomingwarning> -tw=<outgoingwarning> -rc=<incomingcritical> -tc=<outgoingcritical> -i=<interface>
```

## Output example
```
root@node-d2:/usr/home/freebsd # /usr/local/libexec/nagios/bsd_network -rw 200 -rc 220 -tw 150 -tc 200 -i vtnet0
OK -  The current RX is 9.07 kbit/s and TX is 16.93 kbit/s|rx=9.070000;200;220;; tx=16.930000;150;200;;
```

## Copyright and licensing
Distributed under [2-Clause BSD License](https://github.com/araujobsd/icinga2-network-plugin/blob/master/LICENSE).
