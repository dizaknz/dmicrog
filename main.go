package main

import (
	"time"

	"dmicrog/handler"
	dmicrog "dmicrog/proto"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	log "github.com/micro/go-micro/v2/logger"
	kafka "github.com/micro/go-plugins/broker/segmentio/v2"
	"github.com/micro/go-plugins/registry/consul/v2"
	kafkalib "github.com/segmentio/kafka-go"
)

const (
	app     = "dmicrog"
	version = "v0.0.1"
)

func main() {
	reg := consul.NewRegistry()
	// NOTE: go-plugin broker sets batch size to 1, this cannot be changed
	b := kafka.NewBroker(
		broker.Addrs("kafka:9092"),
		kafka.WriterConfig(
			kafkalib.WriterConfig{
				Balancer:      &kafkalib.Hash{},
				MaxAttempts:   3,
				QueueCapacity: 1000,
				BatchTimeout:  100 * time.Millisecond,
				ReadTimeout:   3 * time.Second,
				WriteTimeout:  3 * time.Second,
				Async:         false,
				RequiredAcks:  -1,
			},
		),
		kafka.ReaderConfig(
			kafkalib.ReaderConfig{
				QueueCapacity:   1000,
				MinBytes:        1024,
				MaxBytes:        1024 * 1024,
				MaxWait:         100 * time.Millisecond,
				ReadLagInterval: 5 * time.Second,
				GroupBalancers: []kafkalib.GroupBalancer{
					kafkalib.RangeGroupBalancer{},
					kafkalib.RoundRobinGroupBalancer{},
				},
				HeartbeatInterval:     3 * time.Second,
				WatchPartitionChanges: true,
				RetentionTime:         1 * time.Hour,
				MaxAttempts:           3,
			},
		),
	)
	srv := micro.NewService(
		micro.Name(app),
		micro.Version(version),
		micro.Registry(reg),
		micro.Broker(b),
	)
	dmicrog.RegisterDmicrogHandler(
		srv.Server(),
		handler.NewHandler(b),
	)
	srv.Init()
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
