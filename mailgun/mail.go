package mailgun

import (
	notify "github.com/appscode/go-notify"
	h2t "github.com/jaytaylor/html2text"
	"github.com/kelseyhightower/envconfig"
	mailgun "github.com/mailgun/mailgun-go"
)

const Uid = "mailgun"

type Options struct {
	Domain       string // MAILGUN_DOMAIN
	ApiKey       string // MAILGUN_API_KEY
	PublicApiKey string // MAILGUN_PUBLIC_API_KEY
}

type client struct {
	to      []string
	from    string
	subject string
	body    string
	html    bool
	tag     string

	mg mailgun.Mailgun
}

var _ notify.ByEmail = &client{}

func New(opt Options) *client {
	return &client{
		mg: mailgun.NewMailgun(opt.Domain, opt.ApiKey, opt.PublicApiKey),
	}
}

func Default() (*client, error) {
	var opt Options
	err := envconfig.Process(Uid, &opt)
	if err != nil {
		return nil, err
	}
	return New(opt), nil
}

func (c *client) From(from string) {
	c.from = from
}

func (c *client) WithSubject(subject string) {
	c.subject = subject
}

func (c *client) WithBody(body string) {
	c.body = body
}

func (c *client) To(to string, cc ...string) {
	c.to = append([]string{to}, cc...)
}

func (c *client) Send() error {
	text := c.body
	if c.html {
		if t, err := h2t.FromString(c.body); err == nil {
			text = t
		}
	}
	msg := c.mg.NewMessage(c.from, c.subject, text, c.to...)
	if c.html {
		msg.SetHtml(c.body)
	}
	if c.tag != "" {
		msg.AddTag(c.tag)
	}
	msg.SetTracking(true)
	msg.SetTrackingClicks(true)
	msg.SetTrackingOpens(true)
	response, id, err := c.mg.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) SendHtml() error {
	c.html = true
	return c.Send()
}
