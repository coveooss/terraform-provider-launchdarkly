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

func (c *Client) Get(path string, expectedStatus int, response interface{}) error {
	return c.execute("GET", path, nil, expectedStatus, response)
}

func (c *Client) Post(path string, body interface{}, expectedStatus int, response interface{}) error {
	return c.execute("POST", path, body, expectedStatus, response)
}

func (c *Client) Patch(path string, body interface{}, expectedStatus int, response interface{}) error {
	return c.execute("PATCH", path, body, expectedStatus, response)
}

func (c *Client) Delete(path string, expectedStatus int) error {
	return c.execute("DELETE", path, nil, expectedStatus, nil)
}

func (c *Client) execute(method string, path string, body interface{}, expectedStatus int, response interface{}) error {
	requestBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, "https://app.launchdarkly.com/api/v2"+path, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", c.AccessToken)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != expectedStatus {
		return errors.New(method + " " + path + " did not return expected HTTP status code " + strconv.Itoa(expectedStatus) + ". Got " + strconv.Itoa(resp.StatusCode))
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if response != nil {
		err = json.Unmarshal(responseBody, response)
		if err != nil {
			return err
		}
	}

	return nil
}
