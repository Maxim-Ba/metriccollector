package handlers

import "errors"


var ErrNoMetricName = errors.New("no name metrics")
var ErrNoMetricsType = errors.New("not allowed metric type")
var ErrWrongValue = errors.New("wrong value")
