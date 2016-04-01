User Guide
==========

Table of contents:

* [Running goma agent](#agent)
* [Client commands](#client)
* [Defining monitors](#define)
* [Probes](#probes)
* [Filters](#filters)
* [Actions](#actions)
* [REST API](#api)

<a name="agent" />
Running goma agent
------------------

`goma` command works as a monitoring agent (server) if no command
is given at command-line.

`goma` does not provide so-called *daemon* mode.  Please use [systemd][]
or [upstart][] to run it in the background.

At startup, `goma` will load [TOML][] configuration files from a
directory (default is `/usr/local/etc/goma`).  Each file can define
multiple monitors as described below.

<a name="client" />
Client commands
---------------

`goma` works as clients if a command is given at command-line.

By default, `goma` connects to the agent running on "localhost:3838".
Use `-s` option to specify other addresses.

### list

`goma list` lists all registered monitors.

### show

`goma show ID` show the status of a monitor for ID.
The ID can be identified by list command.

### start

`goma start ID` starts the monitor for ID.

### stop

`goma stop ID` stops the monitor for ID.

### register

`goma register FILE` loads monitor definitions from a TOML file,
register them into the agent, then starts the new monitors.

If FILE is "-", definitions are read from stdin.

### unregister

`goma unregister ID` stops and unregister the monitor for ID.

### verbosity

`goma verbosity LEVEL` changes the logging threshold.  
Available levels are: `debug`, `info`, `warn`, `error`, `critical`

`goma verbosity` queries the current logging threshold.

<a name="define" />
Defining monitors
-----------------

Monitors can be defined in a TOML file like this:

```
[[monitor]]
name = "monitor1"
interval = 10
timeout = 1
min = 0.0
max = 0.3

  [monitor.probe]
  type = "exec",
  command = ["/some/probe/cmd"]

  [monitor.filter]
  type = "average"

  [[monitor.actions]]
  type = "exec"
  command = ["/some/action/cmd"]
```

| Key | Type | Default | Required | Description |
| --- | ---- | ------: | -------- | ----------- |
| `name` | string | | Yes | Descriptive monitor name. |
| `interval` | int | 60 | No | Interval seconds between probes. |
| `timeout` | int | 59 | No | Timeout seconds for a probe. |
| `min` | float | 0.0 | No | The minimum of the normal probe output. |
| `max` | float | 0.0 | No | The maximum of the normal probe output. |
| `probe` | table | | Yes | Probe properties.  See below. |
| `filter` | table | | No | Filter properties.  See below. |
| `actions` | list of table | | Yes | List of action properties.  See below. |

[Annotated sample file](sample.toml).

<a name="probes" />
Probes
------

<a name="filters" />
Filters
-------

<a name="actions" />
Actions
-------

<a name="api" />
REST API
--------

[systemd]: https://www.freedesktop.org/wiki/Software/systemd/
[upstart]: http://upstart.ubuntu.com/
[TOML]: https://github.com/toml-lang/toml
