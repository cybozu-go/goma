![Goma image](goma.png)
[![Build Status](https://travis-ci.org/cybozu-go/goma.svg?branch=master)](https://travis-ci.org/cybozu-go/goma)

Goma is:

* Japanese name of sesame seeds, and
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

Use Go 1.6 or better.

```
go get github.com/cybozu-go/goma/cmd/goma
```

License
-------

[MIT][]

[DevOps]: https://en.wikipedia.org/wiki/DevOps
[MIT]: https://opensource.org/licenses/MIT
