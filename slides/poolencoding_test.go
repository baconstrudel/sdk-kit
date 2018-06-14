package slides

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkPoolEncoding(b *testing.B) {
	testcases := []interface{}{"foo", "afsödjlköalkjfsdöfadsljk", struct{ I int }{42}, "afdöjlkfasdöjlkfsadöjlkfdsaljökfadsöljkfasdljökafsdölkjfasdöjlkfsdaöjlkfdasjlökfadsöljkfdsalökjfdsajölkadsfjlköfdsaölkjarewpiouüögdfxlkhj.uiwtaoesrpöbgöhjdlvkxcesatruiopüs"}
	for i, v := range testcases {
		b.Run(fmt.Sprintf("pooled %d", i), func(b *testing.B) {
			pool := New(func(w io.Writer) encoder {
				return json.NewEncoder(w)
			}, func(r io.Reader) decoder {
				return json.NewDecoder(r)
			})
			for i := 0; i < b.N; i++ {
				encoder := pool.Encoder()
				err := encoder.Encode(v)
				assert.NoError(b, err)
				_, err = io.Copy(ioutil.Discard, encoder.Buf)
				assert.NoError(b, err)
				pool.PutBackEncoder(encoder)
			}
		})
		b.Run(fmt.Sprintf("not pooled %d", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				buf := bytes.Buffer{}
				encoder := json.NewEncoder(&buf)
				err := encoder.Encode(v)
				assert.NoError(b, err)
				_, err = io.Copy(ioutil.Discard, &buf)
				assert.NoError(b, err)
			}
		})
	}

}
