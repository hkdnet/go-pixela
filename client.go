package pixela

import (
	"io/ioutil"
	"log"
	"net/http"
)

var DefaultClient Client = Client{
	Logger:     log.New(ioutil.Discard, "", 0),
	HTTPClient: http.DefaultClient,
}

type Client struct {
	Logger     *log.Logger
	HTTPClient *http.Client

	// Auth
	username string
	token    string
}

func (c *Client) Auth(username, token string) {
	c.username = username
	c.token = token
}
func (c *Client) IsAuthenticated() bool {
	return c.username != "" && c.token != ""
}
