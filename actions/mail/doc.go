/*
Package mail implements "mail" action type that send mails.

The mail body and subject can be customized by text/template:
https://golang.org/pkg/text/template/

The template is rendered with this struct:

    struct {
        Monitor   string    // The monitor name.
        Host      string    // The hostname where goma server is running.
        Time      time.Time // The time of the event.
        Event     string    // One of "init", "fail", or "recover".
        Value     float64   // The probe(filter) value.  Set on failure.
        Duration  int       // Failure duration in seconds.  Set on recovery.
        Version   string    // Goma version such as "0.1".
    }

The constructor takes these parameters:

    Name        Type               Default       Description
    from        string                           Sender mail address.  Required.
    to          []string           nil           Destination mail addresses.
    init_to     []string           nil           Addresses for "init".
    fail_to     []string           nil           Addresses for "fail".
    recover_to  []string           nil           Addresses for "recover".
    subject     string             (See source)  Subject template.
    body        string             (See source)  Mail body template.
    server      string             localhost:25  SMTP server address.
    user        string                           SMTP auth user.  Optional.
    password    string                           SMTP auth password.  Optional.
    header      map[string]string  nil           Extra headers.
    bcc         bool               false         If true, suppress To header.

If no destination address is given for an event, mail is not sent.
For example, mail is not sent on "init" event if both to and init_to are nil.

Extra headers must begin with "X-" for security reasons.
*/
package mail
