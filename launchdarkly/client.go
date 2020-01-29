package launchdarkly

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
	"strconv"
)

type Client struct {
	AccessToken string
}

func (c *Client) GetStatus(url string) (int, error) {
	status, _, err := c.execute("GET", url, nil, []int{}, 0)
	return status, err
}

func (c *Client) Get(url string, expectedStatus []int) (interface{}, error) {
	_, response, err := c.execute("GET", url, nil, expectedStatus, 0)

	var parsedResponse interface{}
	json.Unmarshal(response, &parsedResponse)

	return parsedResponse, err
}

func (c *Client) GetInto(url string, expectedStatus []int, target interface{}) error {
	_, response, err := c.execute("GET", url, nil, expectedStatus, 0)

	json.Unmarshal(response, target)

	return err
}

func (c *Client) Post(url string, body interface{}, expectedStatus []int, target interface{}) error {
	_, response, err := c.execute("POST", url, body, expectedStatus, 0)

	json.Unmarshal(response, target)

	return err
}

func (c *Client) Patch(url string, body interface{}, expectedStatus []int, numberOfRetry int) ([]byte, error) {
	_, response, err := c.execute("PATCH", url, body, expectedStatus, numberOfRetry)
	return response, err
}

func (c *Client) Delete(url string, expectedStatus []int) error {
	_, _, err := c.execute("DELETE", url, nil, expectedStatus, 0)
	return err
}

func (c *Client) execute(method string, url string, body interface{}, expectedStatus []int, numberOfRetry int) (int, []byte, error) {
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

	defer resp.Body.Close()
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
			if numberOfRetry > 0 {
				toRetry := false
				for _, status := range []int{429} {
					if status == resp.StatusCode {
						toRetry = true
						break
					}
				}
				if toRetry {
					println("Will retry " + method + " " + url + " after one minute")
					time.Sleep(time.Minute)
					return c.execute(method, url, body, expectedStatus, numberOfRetry - 1)
				}
			} 
			return resp.StatusCode, nil, errors.New(method + " " + url + " did not return one of the expected HTTP status codes. Got HTTP " + strconv.Itoa(resp.StatusCode) + "\n" + string(responseBody))
		}
	} 

	return resp.StatusCode, responseBody, nil
}
