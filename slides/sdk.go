package slides

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func RequestError(baseError error) error {
	return fmt.Errorf("Unable to create request (%s)", baseError.Error())
}

func EncodingError(baseError error) error {
	return fmt.Errorf("Unable to encode object (%s)", baseError.Error())
}

func AuthenticationError(baseError error) error {
	return fmt.Errorf("Unable to add authentication to request (%s)", baseError.Error())
}

func RequestExecutionError(baseError error) error {
	return fmt.Errorf("Unable to execute request (%s)", baseError.Error())
}

func DecodingError(v interface{}, baseError error) error {
	return fmt.Errorf("Unable to decode response into %T %v, (%s)", v, v, baseError)
}

type Client struct {
	*http.Client
	AuthenticateRequest func(req *http.Request) error
	ValidateStatusCode  func(int) error
}

func DoNotValidate(int) error { return nil }

func AuthenticateWithStaticToken(header, token string) func(req *http.Request) error {
	return func(req *http.Request) error {
		req.Header.Add(header, token)
		return nil
	}
}

func (c Client) Do(method string, url string, payload interface{}, response interface{}) error {
	var payloadBuff bytes.Buffer

	if payload != nil {
		err := json.NewEncoder(&payloadBuff).Encode(payload)
		if err != nil {
			return EncodingError(err)
		}
	}

	req, err := http.NewRequest(method, url, &payloadBuff)
	if err != nil {
		return RequestError(err)
	}

	err = c.AuthenticateRequest(req)
	if err != nil {
		return AuthenticationError(err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return RequestExecutionError(err)
	}
	defer resp.Body.Close()

	err = c.ValidateStatusCode(resp.StatusCode)
	if err != nil {
		return err
	}

	if response != nil {
		if w, ok := response.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(response)
			if err != nil {
				err = DecodingError(response, err)
			}
		}
		if err != nil {
			return err
		}
	}

	return nil
}
