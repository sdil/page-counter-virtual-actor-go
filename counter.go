package main

import (
	"context"
	"fmt"

	"github.com/tochemey/goakt/v3/actor"
)

type Counter struct {
	Id string
	Value int
}

var _ actor.Grain = (*Counter)(nil)

func NewCounter(id string) actor.GrainFactory {
	return func(ctx context.Context) (actor.Grain, error) {
		return &Counter{Id: id}, nil
	}
}

func (c *Counter) OnActivate(ctx context.Context, props *actor.GrainProps) error {
	fmt.Println("Activating counter", c.Id)
	c.Value = 0
	return nil
}

func (c *Counter) OnDeactivate(ctx context.Context, props *actor.GrainProps) error {
	fmt.Println("Deactivating counter", c.Id)
	return nil
}

func (c *Counter) OnReceive(ctx *actor.GrainContext) {
	fmt.Println("Received message", c.Id)
	switch ctx.Message().(type) {
	case *IncrementRequest:
		c.Value++
		ctx.Response(&IncrementResponse{Value: int32(c.Value)})
	default:
		ctx.Unhandled()
	}
}
