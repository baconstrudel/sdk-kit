package slides

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSDKGet(t *testing.T) {
	c := Client{
		http.DefaultClient,
		AuthenticateWithStaticToken("secret", "snape kills dumbledorre"),
		DoNotValidate,
	}

	err := c.Do("GET", "http://localhost:9999/foo", nil, nil)
	assert.NoError(t, err)
}

func TestSDKPost(t *testing.T) {
	c := Client{
		http.DefaultClient,
		AuthenticateWithStaticToken("foo", "bar"),
		DoNotValidate,
	}

	var in string
	var out string
	in = "random string"
	err := c.Do("GET", "http://localhost:9999/echo", in, &out)
	assert.NoError(t, err)
	assert.Equal(t, in, out)
}

func TestRawResponse(t *testing.T) {
	c := Client{
		http.DefaultClient,
		AuthenticateWithStaticToken("foo", "bar"),
		DoNotValidate,
	}

	var in string
	out := bytes.Buffer{}
	in = "random string"
	err := c.Do("GET", "http://localhost:9999/echo", in, &out)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%q\n", in), out.String())
}

func TestValidateStatus(t *testing.T) {
	validation := func(i int) error {
		if i == 404 {

			return errors.New("invalid user")
		}
		return nil
	}

	c := Client{
		http.DefaultClient,
		AuthenticateWithStaticToken("secret", "snape kills dumbledorre"),
		validation,
	}

	err := c.Do("GET", "http://localhost:9999/foo", nil, nil)
	assert.EqualError(t, err, "invalid user")
}
