package msgpack

type Codec[T any] struct {
	Decode func(*Decoder) (T, error)
	Encode func(*Encoder, T) error
}
