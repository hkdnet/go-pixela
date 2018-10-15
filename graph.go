package pixela

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type CreateGraphParam struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Unit  string `json:"unit"`
	Type  string `json:"type"`
	Color string `json:"color"`
}

type ShowGraphsResponse struct {
	Graphs []Graph `json:"graphs"`
}

type Graph struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Unit  string `json:"unit"`
	Type  string `json:"type"`
	Color string `json:"color"`
}

func (c *Client) CreateGraph(param *CreateGraphParam) error {
	if !c.IsAuthenticated() {
		return errors.New("this client is not authenticated")
	}
	url := fmt.Sprintf("https://pixe.la/v1/users/%s/graphs", c.username)
	c.Logger.Printf("%#v", param)
	b, err := json.Marshal(param)
	if err != nil {
		return errors.Wrap(err, "cannot marshal create graph param")
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
		return fmt.Errorf("cannot create a graph: %s", ret.Message)
	}
	return nil
}

func (c *Client) ShowGraphs() ([]Graph, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("this client is not authenticated")
	}
	url := fmt.Sprintf("https://pixe.la/v1/users/%s/graphs", c.username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a new request")
	}

	req.Header.Set("X-USER-TOKEN", c.token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()

	var r io.Reader = resp.Body
	// r = io.TeeReader(resp.Body, os.Stderr) // for debug
	var ret ShowGraphsResponse
	err = json.NewDecoder(r).Decode(&ret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response body")
	}

	return ret.Graphs, nil
}
