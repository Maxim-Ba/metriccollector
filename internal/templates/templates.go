package templates

import (
	"fmt"

	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
)

func GetAllMetricsHTMLPage(m *[]metrics.MetricDTO) string {
	const titlepageStart = `
	<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Document</title>
</head>
<body>
<h1>Метрики</h1>`
   

	const titlepageEnd =`</body>
</html>
	`
var body string 
	for _, metric := range *m {
		if metric.MetricType == "gauge"{

			body += fmt.Sprintf("Тип: %s, Метрика: %s Значение %f <br/>", metric.MetricType, metric.MetricName, metric.Value)
		} else {
			body += fmt.Sprintf("Тип: %s, Метрика: %s Значение %d <br/>", metric.MetricType, metric.MetricName, int64(metric.Value))
		}
}
	return titlepageStart + body + titlepageEnd
}
