package mail

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/mail"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cybozu-go/goma"
)

const (
	testAddress = "localhost:13840"
)

var (
	chServer <-chan *maildata

	envFrom     = os.Getenv("TEST_MAIL_FROM")
	envTo       = os.Getenv("TEST_MAIL_TO")
	envServer   = os.Getenv("TEST_MAIL_HOST")
	envUser     = os.Getenv("TEST_MAIL_USER")
	envPassword = os.Getenv("TEST_MAIL_PASSWORD")
)

func TestMain(m *testing.M) {
	flag.Parse()

	l, err := net.Listen("tcp", testAddress)
	if err != nil {
		log.Fatal(err)
	}
	s, ch := newServer(10)
	chServer = ch
	go func() {
		s.serve(l)
	}()
	os.Exit(m.Run())
}

func TestHeaderPattern(t *testing.T) {
	t.Parallel()

	if headerPattern.MatchString("X-hoge fuga") {
		t.Error("X-hoge fuga")
	}
	if headerPattern.MatchString("hoge-fuga") {
		t.Error("hoge-fuga")
	}
	if headerPattern.MatchString("X-hoge:") {
		t.Error("X-hoge:")
	}
	if !headerPattern.MatchString("X-123-Hoge-fuga") {
		t.Error("X-123-Hoge-fuga")
	}
	if !headerPattern.MatchString("x-123-hoge-fuga") {
		t.Error("x-123-hoge-fuga")
	}
}

func TestConstruct(t *testing.T) {
	t.Parallel()

	_, err := construct(nil)
	if err != goma.ErrNoKey {
		t.Error(`err != goma.ErrNoKey`)
	}

	// Invalid mail address
	_, err = construct(map[string]interface{}{
		"from": "3383829289298289222 28923982398 383892389 292398",
	})
	if err == nil {
		t.Error("from should be invalid")
	}

	// Invalid header
	_, err = construct(map[string]interface{}{
		"from": "Hirotaka Yamamoto <ymmt@example.org>",
		"header": map[string]interface{}{
			"Reply-To": "reply@example.org",
		},
	})
	if err == nil {
		t.Error("header should be invalid")
	}

	a, err := construct(map[string]interface{}{
		"from": "Hirotaka Yamamoto <ymmt@example.org>",
	})
	if err != nil {
		t.Fatal(err)
	}
	if a.(*action).from.Address != "ymmt@example.org" {
		t.Error(`a.from.Address != "ymmt@example.org"`)
	}
	if a.(*action).server != defaultServer {
		t.Error(`a.(*action).server != defaultServer`)
	}

	a, err = construct(map[string]interface{}{
		"from": "Hirotaka Yamamoto <ymmt@example.org>",
		"to": []interface{}{
			"Hirotaka Yamamoto <ymmt@example.org>",
			"ymmt2@example.org",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if a.(*action).to[1].Address != "ymmt2@example.org" {
		t.Error(`a.to[1].Address != "ymmt2@example.org"`)
	}

	a, err = construct(map[string]interface{}{
		"from": "Hirotaka Yamamoto <ymmt@example.org>",
		"init_to": []interface{}{
			"Hirotaka Yamamoto <ymmt@example.org>",
			"ymmt2@example.org",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(a.(*action).to) != 0 {
		t.Error(`len(a.(*action).to) != 0`)
	}
	if a.(*action).initTo[1].Address != "ymmt2@example.org" {
		t.Error(`a.initTo[1].Address != "ymmt2@example.org"`)
	}

	a, err = construct(map[string]interface{}{
		"from": "Hirotaka Yamamoto <ymmt@example.org>",
		"fail_to": []interface{}{
			"Hirotaka Yamamoto <ymmt@example.org>",
			"ymmt2@example.org",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(a.(*action).to) != 0 {
		t.Error(`len(a.(*action).to) != 0`)
	}
	if a.(*action).failTo[1].Address != "ymmt2@example.org" {
		t.Error(`a.failTo[1].Address != "ymmt2@example.org"`)
	}

	a, err = construct(map[string]interface{}{
		"from": "Hirotaka Yamamoto <ymmt@example.org>",
		"recover_to": []interface{}{
			"Hirotaka Yamamoto <ymmt@example.org>",
			"ymmt2@example.org",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(a.(*action).to) != 0 {
		t.Error(`len(a.(*action).to) != 0`)
	}
	if a.(*action).recoverTo[1].Address != "ymmt2@example.org" {
		t.Error(`a.recoverTo[1].Address != "ymmt2@example.org"`)
	}
}

func TestInitMail(t *testing.T) {
	a, err := construct(map[string]interface{}{
		"from": "Hirotaka Yamamoto <ymmt@example.org>",
		"to": []interface{}{
			"Hirotaka Yamamoto <y@example.org>",
			"ymmt2@example.org",
		},
		"init_to": []interface{}{
			"kazu@example.org",
		},
		"server": testAddress,
	})
	if err != nil {
		t.Error(err)
	}

	err = a.Init("monitor1")
	if err != nil {
		t.Fatal(err)
	}

	data := <-chServer
	if data.from != "ymmt@example.org" {
		t.Error(`data.from != "ymmt@example.org"`)
	}
	if len(data.to) != 3 {
		t.Error(`len(data.to) != 3`)
	}
	msg, err := mail.ReadMessage(strings.NewReader(data.data))
	if err != nil {
		t.Error(err)
	}
	if msg.Header.Get("From") != `"Hirotaka Yamamoto" <ymmt@example.org>` {
		t.Error("?", msg.Header.Get("From"))
	}
	al, err := msg.Header.AddressList("To")
	if err != nil {
		t.Error(err)
	}
	if len(al) != 3 {
		t.Error(`len(al) != 3`)
	}
	_, err = msg.Header.Date()
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(msg.Header.Get("Subject"), "monitor1") {
		t.Error(`!strings.Contains(msg.Header.Get("Subject"), "monitor1")`)
	}
	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Contains(body, []byte("Event: init")) {
		t.Error(`!bytes.Contains(body, []byte("Event: init"))`)
	}
}

func TestFailMail(t *testing.T) {
	a, err := construct(map[string]interface{}{
		"from": "Hirotaka Yamamoto <ymmt@example.org>",
		"fail_to": []interface{}{
			"kazu@example.org",
		},
		"server": testAddress,
	})
	if err != nil {
		t.Error(err)
	}

	err = a.Fail("monitor1", 123.45)
	if err != nil {
		t.Fatal(err)
	}

	data := <-chServer
	if data.from != "ymmt@example.org" {
		t.Error(`data.from != "ymmt@example.org"`)
	}
	if len(data.to) != 1 {
		t.Error(`len(data.to) != 1`)
	}
	msg, err := mail.ReadMessage(strings.NewReader(data.data))
	if err != nil {
		t.Error(err)
	}
	al, err := msg.Header.AddressList("To")
	if err != nil {
		t.Error(err)
	}
	if len(al) != 1 {
		t.Error(`len(al) != 1`)
	}
	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Contains(body, []byte("Event: fail")) {
		t.Error(`!bytes.Contains(body, []byte("Event: fail"))`)
	}
	if !bytes.Contains(body, []byte("Value: 123.45")) {
		t.Error(`!bytes.Contains(body, []byte("Value: 123.45"))`)
	}
}

func TestRecoverMail(t *testing.T) {
	a, err := construct(map[string]interface{}{
		"from": "Hirotaka Yamamoto <ymmt@example.org>",
		"to": []interface{}{
			"hogefuga@example.org",
			"abc@example.org",
		},
		"fail_to": []interface{}{
			"kazu@example.org",
		},
		"server": testAddress,
	})
	if err != nil {
		t.Error(err)
	}

	err = a.Recover("monitor1", 39120*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	data := <-chServer
	if data.from != "ymmt@example.org" {
		t.Error(`data.from != "ymmt@example.org"`)
	}
	if len(data.to) != 2 {
		t.Error(`len(data.to) != 2`)
	}
	msg, err := mail.ReadMessage(strings.NewReader(data.data))
	if err != nil {
		t.Error(err)
	}
	al, err := msg.Header.AddressList("To")
	if err != nil {
		t.Error(err)
	}
	if len(al) != 2 {
		t.Error(`len(al) != 2`)
	}
	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Contains(body, []byte("Event: recover")) {
		t.Error(`!bytes.Contains(body, []byte("Event: recover"))`)
	}
	if !bytes.Contains(body, []byte("Duration: 39")) {
		t.Error(`!bytes.Contains(body, []byte("Duration: 39"))`)
	}
}

func TestBcc(t *testing.T) {
	a, err := construct(map[string]interface{}{
		"from": "Hirotaka Yamamoto <ymmt@example.org>",
		"to": []interface{}{
			"hogefuga@example.org",
			"abc@example.org",
		},
		"bcc":    true,
		"server": testAddress,
	})
	if err != nil {
		t.Error(err)
	}

	err = a.Init("monitor1")
	if err != nil {
		t.Fatal(err)
	}

	data := <-chServer
	msg, err := mail.ReadMessage(strings.NewReader(data.data))
	if err != nil {
		t.Error(err)
	}
	if len(msg.Header.Get("To")) != 0 {
		t.Error(`len(msg.Header.Get("To")) != 0`)
	}
	if len(msg.Header.Get("Bcc")) != 0 {
		t.Error(`len(msg.Header.Get("Bcc")) != 0`)
	}
}

func TestSubject(t *testing.T) {
	_, err := construct(map[string]interface{}{
		"from": "Hirotaka Yamamoto <ymmt@example.org>",
		"to": []interface{}{
			"hogefuga@example.org",
			"abc@example.org",
		},
		"subject": `test subject "{{ .NoSuchKey }}"`,
		"server":  testAddress,
	})
	if err == nil {
		t.Error("subject is not a valid template")
	}

	a, err := construct(map[string]interface{}{
		"from": "Hirotaka Yamamoto <ymmt@example.org>",
		"to": []interface{}{
			"hogefuga@example.org",
			"abc@example.org",
		},
		"subject": `test subject "{{ .Event }}"`,
		"server":  testAddress,
	})
	if err != nil {
		t.Error(err)
	}

	err = a.Init("monitor1")
	if err != nil {
		t.Fatal(err)
	}

	data := <-chServer
	msg, err := mail.ReadMessage(strings.NewReader(data.data))
	if err != nil {
		t.Error(err)
	}
	if msg.Header.Get("Subject") != `test subject "init"` {
		t.Error("msg.Header.Get(\"Subject\") != `test subject \"init\"`")
	}
}

func TestBody(t *testing.T) {
	_, err := construct(map[string]interface{}{
		"from": "Hirotaka Yamamoto <ymmt@example.org>",
		"to": []interface{}{
			"hogefuga@example.org",
			"abc@example.org",
		},
		"body":   `test body "{{ .NoSuchKey }}"`,
		"server": testAddress,
	})
	if err == nil {
		t.Error("body is not a valid template")
	}

	a, err := construct(map[string]interface{}{
		"from": "Hirotaka Yamamoto <ymmt@example.org>",
		"to": []interface{}{
			"hogefuga@example.org",
			"abc@example.org",
		},
		"body":   `test body "{{ .Event }}"`,
		"server": testAddress,
	})
	if err != nil {
		t.Error(err)
	}

	err = a.Init("monitor1")
	if err != nil {
		t.Fatal(err)
	}

	data := <-chServer
	msg, err := mail.ReadMessage(strings.NewReader(data.data))
	if err != nil {
		t.Error(err)
	}
	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Contains(body, []byte(`test body "init"`)) {
		t.Error("!bytes.Contains(body, []byte(`test body \"init\"`))")
	}
}

func TestHeader(t *testing.T) {
	_, err := construct(map[string]interface{}{
		"from": "Hirotaka Yamamoto <ymmt@example.org>",
		"to": []interface{}{
			"hogefuga@example.org",
			"abc@example.org",
		},
		"header": map[string]interface{}{
			"Hoge": "Fuga",
		},
		"server": testAddress,
	})
	if err == nil {
		t.Error("header should be invalid")
	}

	a, err := construct(map[string]interface{}{
		"from": "Hirotaka Yamamoto <ymmt@example.org>",
		"to": []interface{}{
			"hogefuga@example.org",
			"abc@example.org",
		},
		"header": map[string]interface{}{
			"X-Hoge": "Fuga",
		},
		"server": testAddress,
	})
	if err != nil {
		t.Error(err)
	}

	err = a.Init("monitor1")
	if err != nil {
		t.Fatal(err)
	}

	data := <-chServer
	msg, err := mail.ReadMessage(strings.NewReader(data.data))
	if err != nil {
		t.Error(err)
	}
	if msg.Header.Get("X-Hoge") != "Fuga" {
		t.Error(`msg.Header.Get("X-Hoge") != "Fuga"`)
	}
}

func TestExternalServer(t *testing.T) {
	t.Parallel()

	if len(envFrom) == 0 || len(envTo) == 0 || len(envServer) == 0 {
		t.Skip()
	}

	a, err := construct(map[string]interface{}{
		"from":     envFrom,
		"to":       []interface{}{envTo},
		"server":   envServer,
		"user":     envUser,
		"password": envPassword,
	})
	if err != nil {
		t.Error(err)
	}

	err = a.Init("monitor1")
	if err != nil {
		t.Fatal(err)
	}
}
