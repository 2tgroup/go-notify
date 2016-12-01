package smtp

import (
	"github.com/appscode/go-notify"
	h2t "github.com/jaytaylor/html2text"
	gomail "gopkg.in/gomail.v2"
)

type Options struct {
	Host               string
	Port               int
	InsecureSkipVerify bool
	Username, Password string
}

type client struct {
	opt  Options
	mail *gomail.Message
	body string
	html bool
}

var _ notify.ByEmail = &client{}

func New(opt Options) *client {
	return &client{
		opt:  opt,
		mail: gomail.NewMessage(),
	}
}

func (c *client) From(from string) {
	c.mail.SetHeader("From", from)
}

func (c *client) WithSubject(subject string) {
	c.mail.SetHeader("Subject", subject)
}
func (c *client) WithBody(body string) {
	c.body = body
}

func (c *client) To(to string, cc ...string) {
	tos := append([]string{to}, cc...)
	c.mail.SetHeader("To", tos...)
}

func (c *client) Send() error {
	if c.html {
		c.mail.SetBody("text/html", c.body)
		if t, err := h2t.FromString(c.body); err == nil {
			c.mail.AddAlternative("text/plain", t)
		}
	} else {
		c.mail.SetBody("text/plain", c.body)
	}

	if c.opt.Username == "" && c.opt.Password == "" {
		d := gomail.NewDialer(c.opt.Host, c.opt.Port, c.opt.Username, c.opt.Password)
		return d.DialAndSend(c.mail)
	} else {
		d := gomail.Dialer{Host: c.opt.Host, Port: c.opt.Port}
		return d.DialAndSend(c.mail)
	}
	return nil
}

func (c *client) SendHtml() error {
	c.html = true
	return c.Send()
}
