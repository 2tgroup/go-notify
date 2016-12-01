package twilio

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/appscode/go-notify"
)

type Options struct {
	AccountSid string
	AuthToken  string
}

type client struct {
	opt Options
	v   url.Values
}

var _ notify.BySMS = &client{}

func New(opt Options) *client {
	return &client{
		opt: opt,
		v:   url.Values{},
	}
}

func (c *client) From(from string) {
	c.v.Set("From", from)
}

func (c *client) WithBody(body string) {
	c.v.Set("Body", body)
}

func (c *client) To(to string) {
	c.v.Set("To", to)
}

func (c *client) Send() error {
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%v/Messages.json", c.opt.AccountSid)

	rb := *strings.NewReader(c.v.Encode())
	h := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest("POST", urlStr, &rb)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.opt.AccountSid, c.opt.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	_, err = h.Do(req)
	return err
}
