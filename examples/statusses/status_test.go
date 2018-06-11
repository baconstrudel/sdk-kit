package status_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type JsonError struct {
	*string
}

func (je *JsonError) UnmarshalJSON(data []byte) error {
	var msg string
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return err
	}
	je.string = &msg
	return nil
}

func (je JsonError) MarshalJSON() ([]byte, error) {
	return json.Marshal(je.string)
}

func (je JsonError) Error() string {
	return *je.string
}

func TestJsonError(t *testing.T) {
	type TestType struct {
		Err JsonError
	}
	foo := "foo"
	want := TestType{JsonError{&foo}}
	b, err := json.Marshal(want)
	assert.NoError(t, err)
	got := TestType{}
	json.Unmarshal(b, &got)
	assert.Equal(t, want, got)
}

func TestSeperateStatusAndValue(t *testing.T) {
	type Status struct {
		Err    string `json:"error"`
		Status int    `json:"status"`
	}
	type Value struct {
		Foo int `json:"foo"`
	}
	type ApiResponse struct {
		Status
		Value
	}
	jsonString := []byte(`{"error":"no error", "status":32, "foo":5}`)
	data := ApiResponse{}
	err := json.Unmarshal(jsonString, &data)
	assert.NoError(t, err)
	assert.Equal(t, 32, data.Status.Status)
	assert.Equal(t, 5, data.Value.Foo)
	assert.Equal(t, "no error", data.Status.Err)

}
