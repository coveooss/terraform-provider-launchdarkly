package launchdarkly

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

func (c *Client) GetStatus(url string) (int, error) {
	status, _, err := c.execute("GET", url, nil, make([]int, 0))
	return status, err
}

func (c *Client) Get(url string, expectedStatus []int) (interface{}, error) {
	_, response, err := c.execute("GET", url, nil, expectedStatus)
	return response, err
}

func (c *Client) Post(url string, body interface{}, expectedStatus []int) (interface{}, error) {
	_, response, err := c.execute("POST", url, body, expectedStatus)
	return response, err
}

func (c *Client) Patch(url string, body interface{}, expectedStatus []int) (interface{}, error) {
	_, response, err := c.execute("PATCH", url, body, expectedStatus)
	return response, err
}

func (c *Client) Delete(url string, expectedStatus []int) error {
	_, _, err := c.execute("DELETE", url, nil, expectedStatus)
	return err
}

func (c *Client) execute(method string, url string, body interface{}, expectedStatus []int) (int, interface{}, error) {
	requestBody, err := json.Marshal(body)
	if err != nil {
		return 0, nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Authorization", c.AccessToken)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}

	println(method + " " + url + " returned HTTP status " + strconv.Itoa(resp.StatusCode))

	if len(expectedStatus) > 0 {
		found := false

		for _, status := range expectedStatus {
			if status == resp.StatusCode {
				found = true
				break
			}
		}

		if !found {
			return resp.StatusCode, nil, errors.New(method + " " + url + " did not return one of the expected HTTP status codes. Got HTTP " + strconv.Itoa(resp.StatusCode) + "\n" + string(responseBody))
		}
	}

	var response interface{}
	if len(responseBody) > 0 {
		err = json.Unmarshal(responseBody, &response)
		if err != nil {
			return resp.StatusCode, nil, err
		}
	}

	return resp.StatusCode, response, nil
}
