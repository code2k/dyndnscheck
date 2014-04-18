# DynDNSCheck

DynDNSCheck is a simple tool for monitoring your DynDNS host.

## How does it work?

When started, DynDNSCheck will get the IP of your DynDNS host and compares it
with your current external IP. If the IPs are different it will notify you by
email.

DynDNSCheck is designed to be triggered periodically.

## Installation

You need a working installation of [Go 1.2.1](http://golang.org) and your
[GOPATH](http://golang.org/doc/code.html#GOPATH) environment variable must be
set.

Install and build DynDNSCheck into your workspace at `$GOPATH`:

```
$ go get github.com/code2k/dyndnscheck
```
You will find the executable at `$GOPATH/bin/dyndnscheck`.

## Configuration

DynDNSCheck is trying to load a configuration named `config.json` in the current
directory. Alternatively you can specify the location of the configuration by
using the command line argument `--config path`.

Use [`config.json.template`](https://raw.githubusercontent.com/code2k/dyndnscheck/master/config.json.template) as a starting point.


