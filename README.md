Goma
====

[![GitHub release](https://img.shields.io/github/release/cybozu-go/goma.svg?maxAge=60)][releases]
[![GoDoc](https://godoc.org/github.com/cybozu-go/goma?status.svg)][godoc]
[![CircleCI](https://circleci.com/gh/cybozu-go/goma.svg?style=svg)](https://circleci.com/gh/cybozu-go/goma)
[![Go Report Card](https://goreportcard.com/badge/github.com/cybozu-go/goma)](https://goreportcard.com/report/github.com/cybozu-go/goma)

Goma is:

* Japanese name of sesame seeds, ![Goma image](goma.png) and
* an extensible monitoring agent written in Go, described here.

Abstract
--------

Goma is a general purpose monitoring server/client.  It can run
multiple monitoring processes concurrently in a single server process.

Basically, goma does active (not passive) monitoring to objects like
web sites or local OS, and kicks actions on failure and/or recovery.

Monitor processes are loaded from configuration files from a directory
at start up, and can be added/started/stopped/removed dynamically via
command-line and REST API.

Goma is designed with [DevOps][] in mind.  Developers can define
and add monitors for their programs easily by putting a rule file
to the configuration directory or by REST API.  Monitoring rules
and actions can be configured flexibly as goma can run arbitrary
commands for them.

### What goma is not

goma is *not* designed for metrics collection.
Use other tools such as Zabbix for that purpose.

Architecture
------------

goma can run multiple independent **monitors** in a single process.

A monitor consists of a **probe**, one or more **actions**, and optionally
a **filter**.  A monitor probes something periodically, and kick actions
for failure when the probe, or filtered result of the probe, reports
failures.  The monitor kicks actions for recovery when the probe or
filtered result of the probe reports recovery from failures.

A probe checks a system and report its status as a floating point number.
*All probes have timeouts*; if a probe cannot return a value before
the timeout, goma will cancel the probe.

A filter manipulates probe outputs; for example, a filter can produce
moving average of probe outputs.

An action implements actions on failures and recoveries.

goma comes with a set of probes, actions, and filters.  New probes,
actions, and filters can be added as compiled-in plugins.

**Pull requests to add new plugins are welcome!**

Usage
-----

Read [USAGE.md](USAGE.md) for details.

Install
-------

The latest officially supported Go version is recommended.

```
go install github.com/cybozu-go/goma/cmd/goma@latest
```

License
-------

[MIT][]

[releases]: https://github.com/cybozu-go/goma/releases
[godoc]: https://godoc.org/github.com/cybozu-go/goma
[DevOps]: https://en.wikipedia.org/wiki/DevOps
[MIT]: https://opensource.org/licenses/MIT
