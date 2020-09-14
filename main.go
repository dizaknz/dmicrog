package main

import (
	"dmicrog/handler"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	kafka "github.com/micro/go-plugins/broker/segmentio/v2"
	"github.com/micro/go-plugins/registry/consul/v2"
)

const (
	app     = "dmicrog"
	version = "v0.0.1"
)

func main() {
	reg := consul.NewRegistry()
	broker := kafka.NewBroker()
	srv := micro.NewService(
		micro.Name(app),
		micro.Version(version),
		micro.Registry(reg),
		micro.Broker(broker),
	)
	srv.Server().NewHandler(handler.NewHandler(broker))
	srv.Init()
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
