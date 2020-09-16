package main

import (
	"context"
	"strconv"
	"time"

	dmicrog "dmicrog/proto"

	"github.com/golang/protobuf/ptypes"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-plugins/registry/consul/v2"
)

func main() {
	log.Info("Running demo")

	reg := consul.NewRegistry()
	client := dmicrog.NewDmicrogService(
		"dmicrog",
		grpc.NewClient(
			client.Registry(reg),
		),
	)

	if err := demo(client, "DEMO", 100); err != nil {
		log.Fatal(err)
	}
}

func demo(client dmicrog.DmicrogService, typ string, num int) error {
	done := make(chan struct{})
	// hack send a message of type auto create topic and bide a while
	if err := call(client, strconv.Itoa(0), typ, time.Now().UTC()); err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	go stream(client, typ, done)
	for i := 1; i < num; i++ {
		if err := call(client, strconv.Itoa(i), typ, time.Now().UTC()); err != nil {
			close(done)
			return err
		}
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(1 * time.Second)
	close(done)
	return nil
}

func call(client dmicrog.DmicrogService, id, typ string, t time.Time) error {
	ts, err := ptypes.TimestampProto(time.Now().UTC())
	if err != nil {
		log.Fatal(err)
	}
	request := &dmicrog.Request{
		Event: &dmicrog.EventMessage{
			Id:        id,
			Typ:       typ,
			Timestamp: ts,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := client.Call(ctx, request)
	if err != nil {
		return err
	}
	log.Infof("Call response: %t %s", response.Success, response.Message)

	return nil
}

func stream(client dmicrog.DmicrogService, typ string, done chan struct{}) error {
	s, err := client.Stream(
		context.Background(),
		&dmicrog.StreamingRequest{
			Typ: typ,
		},
	)
	if err != nil {
		return err
	}
	for {
		select {
		case <-done:
			return nil
		default:
			r, err := s.Recv()
			if err != nil {
				return err
			}
			log.Info("Stream event: ", r)
		}
	}
}
