/*
Package mysql implements "mysql" probe type that test MySQL servers.

The underlying driver is https://github.com/go-sql-driver/mysql .

The value returned from a SELECT query will be the value of the probe.
The SELECT statement should return a floating point value.

The constructor takes these parameters:

    Name       Type     Default   Description
    dsn        string             DSN for MySQL server.  Required.
    query      string             SELECT statement.  Required.
    errval     float64  0         Return value upon an error.

This probe utilizes max_execution_time system variable if available
(for MySQL 5.7.8+).  If not, the probe will kill the running thread
when the deadline expires.
*/
package mysql
