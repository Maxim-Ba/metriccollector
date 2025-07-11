package templates

import (
	"fmt"

	"github.com/Maxim-Ba/metriccollector/internal/constants"
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
)

func GetAllMetricsHTMLPage(m *[]metrics.Metrics) string {
	const titlepageStart = `
	<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Document</title>
</head>
<body>
<h1>Метрики</h1>`

	const titlepageEnd = `</body>
</html>
	`
	var body string
	for _, metric := range *m {
		if metric.MType == constants.Gauge {

			body += fmt.Sprintf("Тип: %s, Метрика: %s Значение %f <br/>", metric.MType, metric.ID, *metric.Value)
		} else {
			body += fmt.Sprintf("Тип: %s, Метрика: %s Значение %d <br/>", metric.MType, metric.ID, int64(*metric.Delta))
		}
	}
	return titlepageStart + body + titlepageEnd
}
