package capnproto2packed

import (
	"errors"
	"io"

	"zombiezen.com/go/capnproto2"

	"github.com/goadesign/goa"
)

// Enforce that codec.Decoder satisfies goa.Decoder at compile time
var (
	_ goa.Decoder = (*ProtoDecoder)(nil)
	_ goa.Encoder = (*ProtoEncoder)(nil)
)

type (
	// Factory uses github.com/ugorji/go/codec to act as an DecoderFactory and EncoderFactory
	Factory struct{}

	// ProtoDecoder stores state between Reset and Decode
	ProtoDecoder struct {
		dec *capnp.Decoder
	}

	// ProtoEncoder stores state between Reset and Encode
	ProtoEncoder struct {
		enc *capnp.Encoder
	}
)

// DecoderFactory is the default factory used by the goa `Consumes` DSL
func DecoderFactory() goa.EncoderFactory {
	return &Factory{}
}

// NewDecoder returns a new capnp.Decoder that satisfies goa.Decoder
func (f *Factory) NewDecoder(r io.Reader) goa.Decoder {
	return &ProtoDecoder{
		dec: capnp.NewPackedDecoder(r),
	}
}

// Decode unmarshals an io.Reader into *capnp.Message v
func (dec *ProtoDecoder) Decode(v interface{}) error {
	msg, ok := v.(*capnp.Message)
	if !ok {
		return errors.New("Cannot decode into struct that doesn't implement *capnp.Message")
	}

	newMsg, err := dec.dec.Decode()
	if err != nil {
		return err
	}

	if newMsg == nil {
		msg = nil
		return nil
	}

	*msg = *newMsg
	return nil
}

// EncoderFactory is the default factory used by the goa `Produces` DSL
func EncoderFactory() goa.EncoderFactory {
	return &Factory{}
}

// NewEncoder returns a new capnp.Encoder that satisfies goa.Encoder
func (f *Factory) NewEncoder(w io.Writer) goa.Encoder {
	return &ProtoEncoder{
		enc: capnp.NewPackedEncoder(w),
	}
}

// Encode marshals a *capnp.Message and writes it to an io.Writer
func (enc *ProtoEncoder) Encode(v interface{}) error {
	msg, ok := v.(*capnp.Message)
	if !ok {
		return errors.New("Cannot encode struct that doesn't implement *capnp.Message")
	}

	var err error

	if err = enc.enc.Encode(msg); err != nil {
		return err
	}

	return nil
}
