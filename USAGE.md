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

`goma` command works as a monitoring agent (server) if "serve" is given
at command-line.  `goma` does not provide so-called *daemon* mode.
Please use [systemd][] or [upstart][] to run it in the background.

At startup, `goma` will load [TOML][] configuration files from a
directory (default is `/usr/local/etc/goma`).  Each file can define
multiple monitors as described below.

<a name="client" />
Client commands
---------------

`goma` works as clients for commands other than "serve".

By default, `goma` connects to the agent running on "localhost:3838".
Use `-s` option to change the address.

* list

    `goma list` lists all registered monitors.

* show

    `goma show ID` show the status of a monitor for ID.
    The ID can be identified by list command.

* start

    `goma start ID` starts the monitor for ID.

* stop

    `goma stop ID` stops the monitor for ID.

* register

    `goma register FILE` loads monitor definitions from a TOML file,
    register them into the agent, then starts the new monitors.

    If FILE is "-", definitions are read from stdin.

* unregister

   `goma unregister ID` stops and unregister the monitor for ID.

* verbosity

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
  command = "/some/probe/cmd"

  [monitor.filter]
  type = "average"

  [[monitor.actions]]
  type = "exec"
  command = "/some/action/cmd"
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

See [annotated sample file](sample.toml).

<a name="probes" />
Probes
------

See GoDoc for construction parameters:

* [exec](https://godoc.org/github.com/cybozu-go/goma/probes/exec)
* [http](https://godoc.org/github.com/cybozu-go/goma/probes/http)

<a name="filters" />
Filters
-------

See GoDoc for construction parameters:

* [average](https://godoc.org/github.com/cybozu-go/goma/filters/average)

<a name="actions" />
Actions
-------

See GoDoc for construction parameters:

* [exec](https://godoc.org/github.com/cybozu-go/goma/actions/exec)
* [http](https://godoc.org/github.com/cybozu-go/goma/actions/http)

<a name="api" />
REST API
--------

### /list

GET will return a list of monitor status objects in JSON:

```javascript
[
    {"id": "0", "name": "monitor1", "running": true, "failing": false},
    ...
]
```

### /register

POST will create and start a new monitor.
The request content-type must be `application/json`.

The request body is a JSON object just like TOML monitor table:

```javascript
{
    "name": "monitor1",
    "interval": 10,
    "timeout": 1,
    "min": 0,
    "max": 0.3,
    "probe": {
        "type": "exec",
        "command": "/some/probe/cmd"
    },
    "filter": {
        "type": "average"
    },
    "actions": [
        {
            "type": "exec",
            "command": "/some/action/cmd"
        },
        ...
    ]
}
```

### /monitor/ID

GET returns monitor status for the given ID.
The response is a JSON object:

```javascript
{
    "id": "0",
    "name": "monitor1",
    "running": true,
    "failing": false
}
```

DELETE will stop and unregister the monitor.

POST can stop or start the monitor.
The request content-type should be `text/plain`.
The request body shall contain either `start` or `stop`.

### /verbosity

GET will return the current verbosity.  
Possible values are: "critical", "error", "warning", "info", and "debug".

PUT or POST will modify the verbosity as given by the request body.
The request content-type should be `text/plain`.
The request body shall contain only the new verbosity level string.

[systemd]: https://www.freedesktop.org/wiki/Software/systemd/
[upstart]: http://upstart.ubuntu.com/
[TOML]: https://github.com/toml-lang/toml
