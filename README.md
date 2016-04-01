![Goma image](goma.png)
[![Build Status](https://travis-ci.org/cybozu-go/goma.png)](https://travis-ci.org/cybozu-go/goma)

goma - an extensible monitoring agent
============================================================

**goma** is a general purpose monitoring agent that can kick actions
on monitor failures and recoveries.  
goma is *not* designed for metrics collection; use other tools such as
Zabbix for that purpose.

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
**Pull requests to add new ones are welcome!**

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

[MIT](https://opensource.org/licenses/MIT)
