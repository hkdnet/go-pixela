package pixela

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

const createUserURL = "https://pixe.la/v1/users"

type CreateUserParam struct {
	Token               string `json:"token"`
	Username            string `json:"username"`
	AgreeTermsOfService string `json:"agreeTermsOfService"`
	NotMinor            string `json:"notMinor"`
}

func (c *Client) CreateUser(param *CreateUserParam) error {
	c.Logger.Printf("%#v", param)
	b, err := json.Marshal(param)
	if err != nil {
		return errors.Wrap(err, "cannot marshal create user param")
	}
	req, err := http.NewRequest("POST", createUserURL, bytes.NewBuffer(b))
	if err != nil {
		return errors.Wrap(err, "failed to create a new request")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send post request")
	}
	defer resp.Body.Close()

	var r io.Reader = resp.Body
	// r = io.TeeReader(resp.Body, os.Stderr) // for debug
	dec := json.NewDecoder(r)
	var ret Response
	err = dec.Decode(&ret)
	if err != nil {
		return errors.Wrap(err, "failed to decode response body")
	}
	if !ret.IsSuccess {
		return fmt.Errorf("cannot create user: %s", ret.Message)
	}
	return nil
}
