// mockup SMTP server.
//
// Only for testing purpose.  Do not work for concurrent access.

package mail

import (
	"io"
	"net"
	"net/textproto"
	"regexp"
	"strings"
)

var (
	pathPattern = regexp.MustCompile(`<(?:[^>:]+:)?([^>]+@[^>]+)>`)
)

type maildata struct {
	from string
	to   []string
	data string
}

type server struct {
	ch chan<- *maildata
}

func newServer(capacity int) (*server, <-chan *maildata) {
	ch := make(chan *maildata, capacity)
	return &server{
		ch: ch,
	}, ch
}

func (s *server) listenAndServe(network, addr string) error {
	l, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	return s.serve(l)
}

func (s *server) serve(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		// don't go.
		s.process(conn)
	}
}

func (s *server) process(c net.Conn) {
	tc := textproto.NewConn(c)
	defer tc.Close()

	doneHello := false
	data := new(maildata) // len(data.from) > 0 means in-transaction.

	reply := func(code int, msg string, cnt bool) error {
		delim := " "
		if cnt {
			delim = "-"
		}
		return tc.Writer.PrintfLine("%d%s%s", code, delim, msg)
	}

	if reply(220, "localhost goma test mail server", false) != nil {
		return
	}

	for {
		l, err := tc.Reader.ReadLine()
		if err != nil {
			return
		}
		ul := strings.ToUpper(l)

		switch {
		case ul == "NOOP" || strings.HasPrefix(ul, "NOOP "):
			if reply(250, "OK", false) != nil {
				return
			}
		case ul == "HELP" || strings.HasPrefix(ul, "HELP "):
			if reply(250, "Supported commands: EHLO HELO MAIL RCPT DATA RSET NOOP QUIT VRFY", false) != nil {
				return
			}
		case ul == "QUIT":
			reply(221, "OK", false)
			return
		case strings.HasPrefix(ul, "VRFY "):
			if reply(252, "cannot verify, but accept anyway", false) != nil {
				return
			}
		case ul == "RSET":
			data = new(maildata)
			if reply(250, "OK", false) != nil {
				return
			}
		case strings.HasPrefix(ul, "EHLO "):
			if doneHello {
				if reply(503, "Duplicate HELO/EHLO", false) != nil {
					return
				}
				continue
			}
			if reply(250, "localhost greets you", true) != nil {
				return
			}
			if reply(250, "8BITMIME", true) != nil {
				return
			}
			if reply(250, "HELP", false) != nil {
				return
			}
			doneHello = true
		case strings.HasPrefix(ul, "HELO "):
			if doneHello {
				if reply(503, "Duplicate HELO/EHLO", false) != nil {
					return
				}
				continue
			}
			if reply(250, "localhost", false) != nil {
				return
			}
			doneHello = true
		case strings.HasPrefix(ul, "MAIL FROM:"):
			if len(data.from) > 0 {
				if reply(503, "nested MAIL command", false) != nil {
					return
				}
				continue
			}
			m := pathPattern.FindStringSubmatch(l)
			if len(m) != 2 {
				if reply(501, "Syntax: MAIL FROM: <address>", false) != nil {
					return
				}
				continue
			}
			data.from = m[1]
			if reply(250, "OK", false) != nil {
				return
			}
		case strings.HasPrefix(ul, "RCPT TO:"):
			if len(data.from) == 0 {
				if reply(503, "need MAIL first", false) != nil {
					return
				}
				continue
			}
			m := pathPattern.FindStringSubmatch(l)
			if len(m) != 2 {
				if reply(501, "Syntax: RCPT TO: <address>", false) != nil {
					return
				}
				continue
			}
			data.to = append(data.to, m[1])
			if reply(250, "OK", false) != nil {
				return
			}
		case ul == "DATA":
			if len(data.from) == 0 {
				if reply(503, "need MAIL first", false) != nil {
					return
				}
				continue
			}
			if len(data.to) == 0 {
				if reply(503, "need RCPT first", false) != nil {
					return
				}
				continue
			}
			if reply(354, "End data with <CR><LF>.<CR><LF>", false) != nil {
				return
			}
			t, err := io.ReadAll(tc.Reader.DotReader())
			if err != nil {
				return
			}
			data.data = string(t)
			s.ch <- data
			data = new(maildata)
			if reply(250, "OK", false) != nil {
				return
			}
		default:
			if reply(500, "unknown command", false) != nil {
				return
			}
		}
	}
}
