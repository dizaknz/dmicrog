package handler

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/micro/go-micro/v2/broker"
	log "github.com/micro/go-micro/v2/logger"

	dmicrog "dmicrog/proto"
)

type Handler struct {
	broker broker.Broker
}

func NewHandler(broker broker.Broker) *Handler {
	return &Handler{
		broker: broker,
	}
}

func (h *Handler) Call(ctx context.Context, req *dmicrog.Request, rsp *dmicrog.Response) error {
	log.Debug("Received event to publish", req.Event)
	topic := req.Event.Typ
	body, err := proto.Marshal(req.Event)
	if err != nil {
		return err
	}
	if err := h.broker.Publish(
		topic,
		&broker.Message{
			Body: body,
		},
	); err != nil {
		log.Error("Failed to publish event", err)
		return err
	}

	rsp.Success = true
	rsp.Message = "200 - OK"

	log.Debug("Published event ", req.Event.Id)

	return nil
}

func (h *Handler) Stream(ctx context.Context, req *dmicrog.StreamingRequest, stream dmicrog.Dmicrog_StreamStream) error {
	log.Debugf("Received stream request for events of type: %s", req.Typ)

	topic := req.Typ
	sub, err := h.broker.Subscribe(
		topic,
		func(ev broker.Event) error {
			msg := &dmicrog.StreamingResponse{
				Event: &dmicrog.EventMessage{},
			}
			if err := proto.Unmarshal(ev.Message().Body, msg.Event); err != nil {
				log.Error("Failed to decode message", err)
				return err
			}
			if err := stream.Send(msg); err != nil {
				log.Error("Failed to send message", err)
				return err
			}
			if err := ev.Ack(); err != nil {
				log.Error("Failed to ack message", err)
				return err
			}
			return nil
		},
	)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-stream.Context().Done():
			log.Debug("Client closed stream")
			return nil
		}
	}
}
