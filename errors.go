package msgpack

import "github.com/ab36245/go-errors"

var Error = errors.Make("msgpack", nil)

var DecodeError = Error.Make("decode")
