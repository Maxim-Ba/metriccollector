package handlers

import (
	"errors"

	"github.com/Maxim-Ba/metriccollector/internal/server/handlers/middleware"
)

var ErrNoMetricName = errors.New("no name metrics")
var ErrNoMetricsType = errors.New("not allowed metric type")
var ErrWrongValue = errors.New("wrong value")
var ErrWrongBodyEncoding = middleware.ErrWrongBodyEncoding
