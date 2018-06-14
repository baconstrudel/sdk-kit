package slides

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
}

func BenchmarkPipe(b *testing.B) {
	longStringBuf := bytes.Buffer{}
	io.Copy(&longStringBuf, io.LimitReader(rand.Reader, 1000))
	testdata := []struct {
		Name string
		ts   TestStruct
	}{
		{"tiny", TestStruct{"f", "b"}},
		{"short", TestStruct{"foo", "bar"}},
		{"mid", TestStruct{"123467890123467890123467890", "123467890123467890123467890"}},
		{"long", TestStruct{"f", longStringBuf.String()}},
	}

	for _, benchcase := range testdata {
		b.Run("using payloadReader "+benchcase.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				pr, pw := io.Pipe()
				wg := sync.WaitGroup{}
				wg.Add(1)

				go func() {
					defer pw.Close()
					json.NewEncoder(pw).Encode(benchcase.ts)
					wg.Done()
				}()
				http.Post("http://localhost:9999", "text", pr)
				wg.Wait()

			}
		})
		b.Run("using buffer "+benchcase.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				buf := bytes.Buffer{}
				err := json.NewEncoder(&buf).Encode(benchcase.ts)
				assert.NoError(b, err)
			}
		})
		b.Run("using pooled buffer "+benchcase.Name, func(b *testing.B) {
			pool := sync.Pool{
				New: func() interface{} {
					return new(bytes.Buffer)
				},
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				buf := pool.Get().(*bytes.Buffer)
				err := json.NewEncoder(buf).Encode(benchcase.ts)
				assert.NoError(b, err)
				buf.Reset()
				pool.Put(buf)
			}
		})
	}
}
