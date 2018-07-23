package main

import (
	"net/http"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"errors"
	"strconv"
)

type Client struct {
	AccessToken string
}

func (c *Client) GetStatus(path string) (int, error) {
	expectedStatus := make(map[int]bool)
	return c.execute("GET", path, nil, expectedStatus, nil)
}

func (c *Client) Get(path string, expectedStatus map[int]bool, response interface{}) error {
	_, err := c.execute("GET", path, nil, expectedStatus, response)
	return err
}

func (c *Client) Post(path string, body interface{}, expectedStatus map[int]bool, response interface{}) error {
	_, err := c.execute("POST", path, body, expectedStatus, response)
	return err
}

func (c *Client) Patch(path string, body interface{}, expectedStatus map[int]bool, response interface{}) error {
	_, err := c.execute("PATCH", path, body, expectedStatus, response)
	return err
}

func (c *Client) Delete(path string, expectedStatus map[int]bool) error {
	_, err := c.execute("DELETE", path, nil, expectedStatus, nil)
	return err
}

func (c *Client) execute(method string, path string, body interface{}, expectedStatus map[int]bool, response interface{}) (int, error) {
	requestBody, err := json.Marshal(body)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest(method, "https://app.launchdarkly.com/api/v2"+path, bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", c.AccessToken)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}

	println(method + " " + path + " returned HTTP status " + strconv.Itoa(resp.StatusCode))

	if len(expectedStatus) > 0 && !expectedStatus[resp.StatusCode] {
		return resp.StatusCode, errors.New(method + " " + path + " did not return one of the expected HTTP status codes. Got HTTP " + strconv.Itoa(resp.StatusCode) + "\n" + string(responseBody))
	}

	if response != nil {
		err = json.Unmarshal(responseBody, response)
		if err != nil {
			return resp.StatusCode, err
		}
	}

	return resp.StatusCode, nil
}
