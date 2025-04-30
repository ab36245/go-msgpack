package msgpack

import "github.com/aivoicesystems/aivoice/common/errors"

var Error = errors.Make("msgpack", nil)

var DecodeError = Error.Make("decode")
