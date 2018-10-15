package pixela

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type PostPixelParam struct {
	Date     string `json:"date"`
	Quantity string `json:"quantity"`
}

func (c *Client) PostPixel(graphID string, param *PostPixelParam) error {
	if !c.IsAuthenticated() {
		return errors.New("this client is not authenticated")
	}
	url := fmt.Sprintf("https://pixe.la/v1/users/%s/graphs/%s", c.username, graphID)
	c.Logger.Printf("%#v", param)
	b, err := json.Marshal(param)
	if err != nil {
		return errors.Wrap(err, "cannot marshal post pixel param")
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return errors.Wrap(err, "failed to create a new request")
	}
	req.Header.Set("X-USER-TOKEN", c.token)

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
		return fmt.Errorf("cannot post a pixel: %s", ret.Message)
	}
	return nil
}
