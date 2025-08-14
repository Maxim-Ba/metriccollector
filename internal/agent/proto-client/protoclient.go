package protoclient

import (
	"context"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	pb "github.com/Maxim-Ba/metriccollector/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Connection struct {
	conn   *grpc.ClientConn
	client pb.MetricsClient
}

func Connect(addr string) (*Connection, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewMetricsClient(conn)
	return &Connection{conn: conn, client: client}, nil
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

func (c *Connection) SendMetrics(metrics []*metrics.Metrics) {
	logger.LogInfo("Connection SendMetrics")

	reqData := make([]*pb.Metric, len(metrics))
	for i, v := range metrics {
		reqData[i] = &pb.Metric{
			Id:    v.ID,
			Type:  v.MType,
			Delta: *v.Delta,
			Value: *v.Value,
		}
	}
	resp, err := c.client.SendMetrics(context.Background(), &pb.SendMetricsRequest{Metrics: reqData})
	if err != nil {
		logger.LogError(err)
	}
	if resp.Error != "" {
		logger.LogError(resp.Error)
	}
}
