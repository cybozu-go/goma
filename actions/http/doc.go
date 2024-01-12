/*
Package http implements "http" action type that send events to HTTP(S) server.

GET or POST form variables are automatically appended as follows:

	Name           Description
	monitor        The monitor name.
	host           Hostname where goma server is running.
	event          One of "init", "fail", or "recover".
	value          The probe(filter) value.  Appended on failure.
	duration       Failure duration in seconds.  Appended on recovery.
	version        Goma version such as "0.1".

The constructor takes these parameters:

	Name         Type               Default  Description
	url_init     string                      URL to access on monitor startup.  Optional.
	url_fail     string                      URL to access on monitor failure.  Optional.
	url_recover  string                      URL to access on monitor recovery.  Optional.
	method       string             GET      HTTP method to use.
	agent        string             goma/0.1 User-Agent string.
	header       map[string]string  nil      HTTP headers.
	params       map[string]string  nil      Additional form parameters.
	timeout      int                30       Timeout seconds for requests.
	                                         Zero means the default timeout.

If URL is not given for an event type, no request is sent for the event.

Proxy can be specified through environment variables.
See net.http.ProxyFromEnvironment for details.

Basic authentication can be used by embedding user:password in URLs.
*/
package http
