package storage

import "errors"


var ErrUnknownMetricName = errors.New("unknown metrics name")

var ErrDatabaseConnection = errors.New("database connection is not initialized")
