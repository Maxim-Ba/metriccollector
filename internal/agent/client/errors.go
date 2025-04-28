package client

import "errors"

var ErrServerInternalError = errors.New("server internal error")
var ErrRequestTimeout = errors.New("request timeout")
