/*
Package exec implements "exec" probe type that runs an arbitrary command.

The value of the probe will be 0 if command exits successfully,
or 1.0 if command timed out or does not exit normally.

If parse is true, the command output (stdout) will be interpreted
as a floating point number, and will be used as the probe value.

The constructor takes these parameters:

	Name       Type      Default   Description
	command    string              The command to run.
	args       []string      nil   Command arguments.
	parse      bool        false   See the above description.
	errval     float64         0   When parse is true and command failed,
	                               this value is returned as the probe value.
	env        []string      nil   Environment variables.  See os.Environ.
*/
package exec
