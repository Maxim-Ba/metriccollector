package protoserver

import (
	// ...
	"context"
	"net"

	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	metricsService "github.com/Maxim-Ba/metriccollector/internal/server/services/metric"
	"github.com/Maxim-Ba/metriccollector/internal/server/storage"
	"google.golang.org/grpc"

	// импортируем пакет со сгенерированными protobuf-файлами
	"github.com/Maxim-Ba/metriccollector/internal/logger"
	pb "github.com/Maxim-Ba/metriccollector/proto"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer
}

func Start(addr string) error {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterMetricsServer(s, &MetricsServer{})

	if err := s.Serve(listen); err != nil {
		return err
	}
	return nil
}

func (s *MetricsServer) SendMetrics(ctx context.Context, in *pb.SendMetricsRequest) (*pb.SendMetricsResponse, error) {
	var response pb.SendMetricsResponse
	metricsSlice := make([]metrics.Metrics, len(in.Metrics))
	for i, v := range in.Metrics {
		metricsSlice[i] = metrics.Metrics{
			ID:    v.Id,
			MType: v.Type,
			Delta: &v.Delta,
			Value: &v.Value,
		}
	}
	err := metricsService.UpdateMany(storage.StorageInstance, &metricsSlice)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	return &response, nil
}
