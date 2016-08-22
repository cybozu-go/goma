/*
Package exec implements "exec" action type that runs arbitrary commands.

Monitor information are passed by environment variables:

    Name           Description
    GOMA_MOINTOR   The name of the monitor.
    GOMA_EVENT     Event name.  One of "init", "fail" or "recover".
    GOMA_VALUE     The probe(filter) value.  Available on failure.
    GOMA_DURATION  Failure duration in seconds.  Available on recovery.
    GOMA_VERSION   Goma version such as "0.1".

The constructor takes these parameters:

    Name         Type      Default  Description
    command      string             The command to run.  Required.
    args         []string  nil      Arguments for the command.
    env          []string  nil      Environment variables.  See os.Environ.
    timeout      int       0        Timeout seconds for command execution.
                                    Zero disables timeout.
    debug        bool      false    If true, command outputs are logged on failure.
*/
package exec
