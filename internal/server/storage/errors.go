package storage

import "errors"

var ErrUnknownMetricName = errors.New("unknown metrics name")

var ErrUnknownMetricType = errors.New("unknown metrics type")

var ErrDatabaseConnection = errors.New("database connection is not initialized")
