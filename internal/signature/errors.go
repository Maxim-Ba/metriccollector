package signature

import "errors"

var ErrInvalidSignature = errors.New("signature is not equal dist signature")
