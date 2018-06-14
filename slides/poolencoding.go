package slides

import (
	"bytes"
	"io"
	"sync"
)

type encoder interface {
	Encode(v interface{}) error
}

type bufEncoder struct {
	Buf     *bytes.Buffer
	encoder encoder
}

func (be *bufEncoder) Encode(v interface{}) error {
	return be.encoder.Encode(v)
}

type decoder interface {
	Decode(v interface{}) error
}

type bufDecoder struct {
	buf     *bytes.Buffer
	decoder decoder
}

func (bd *bufDecoder) Decode(v interface{}) error {
	return bd.decoder.Decode(v)
}

type Pool struct {
	decoderPool sync.Pool
	encoderPool sync.Pool
}

func New(bindEncoder func(io.Writer) encoder, bindDecoder func(io.Reader) decoder) Pool {
	return Pool{
		decoderPool: sync.Pool{
			New: func() interface{} {
				b := bytes.Buffer{}
				return &bufDecoder{
					buf:     &b,
					decoder: bindDecoder(&b),
				}
			},
		},
		encoderPool: sync.Pool{
			New: func() interface{} {
				b := bytes.Buffer{}
				return &bufEncoder{
					Buf:     &b,
					encoder: bindEncoder(&b),
				}
			},
		},
	}
}

func (p *Pool) PutBackEncoder(be *bufEncoder) {
	be.Buf.Reset()
	p.encoderPool.Put(be)
}

func (p *Pool) PutBackDecoder(bd *bufDecoder) {
	bd.buf.Reset()
	p.decoderPool.Put(bd)
}

func (p *Pool) Decoder() *bufDecoder {
	return p.decoderPool.Get().(*bufDecoder)
}

func (p *Pool) Encoder() *bufEncoder {
	return p.encoderPool.Get().(*bufEncoder)
}
