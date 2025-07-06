package generator

import (
	"sync"

	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
)

// worker это наш рабочий, который принимает два канала:
// jobs - канал задач, это входные данные для обработки
// results - канал результатов, это результаты работы воркера

func workerPool(callbacks []func() *metrics.Metrics, maxWorkers int) <-chan *metrics.Metrics {
	var wg sync.WaitGroup

	// создаем буферизованный канал для принятия задач в воркер
	jobs := make(chan (func() *metrics.Metrics), len(callbacks))
	// создаем буферизованный канал для отправки результатов
	results := make(chan *metrics.Metrics, len(callbacks))

	// передаем id, это для наглядности, канал задач и канал результатов
	for w := 0; w < maxWorkers; w++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// в канал задач отправляем какие-то данные
	for _, j := range callbacks {
		jobs <- j
	}
	close(jobs)
	// close(results)
	go func() {
		wg.Wait()      // Ждём завершения всех горутин
		close(results) // Закрываем канал
	}()

	return results
}

func worker(jobs <-chan (func() *metrics.Metrics), results chan<- *metrics.Metrics, wg *sync.WaitGroup) {
	for job := range jobs {
		results <- job()
	}
	defer wg.Done()

}
