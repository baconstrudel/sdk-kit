package slides

import (
	"bytes"
	"crypto/sha512"
	"hash"
	"net/http"
	"net/url"
	"sync"
)

type shared struct {
	client  *http.Client
	baseURL *url.URL
}

type resources struct {
	jsonbuffer *bytes.Buffer
	hasher     *hash.Hash
}

type SDK struct {
	shared
	pool sync.Pool
}

type SDKWorker struct {
}

func NewSDK(shared shared) SDK {
	pool := sync.Pool{
		New: func() interface{} {
			h := sha512.New512_256()
			return resources{
				&bytes.Buffer{},
				&h,
			}
		},
	}
	sdk := SDK{
		shared: shared,
		pool:   pool,
	}
	return sdk
}
