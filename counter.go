package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tochemey/goakt/v3/actor"
)

type Counter struct {
	Path string
	Value int
}

var _ actor.Grain = (*Counter)(nil)

func NewCounter(path string) actor.GrainFactory {
	return func(ctx context.Context) (actor.Grain, error) {
		return &Counter{Path: path}, nil
	}
}

func (c *Counter) OnActivate(ctx context.Context, props *actor.GrainProps) error {
	fmt.Println("Activating counter", c.Path)
	c.Value = 0
	return nil
}

func (c *Counter) OnDeactivate(ctx context.Context, props *actor.GrainProps) error {
	fmt.Println("Deactivating counter", c.Path)
	return nil
}

func (c *Counter) OnReceive(ctx *actor.GrainContext) {
	// fmt.Println("Received message", c.Path)
	now := time.Now()
	defer func() {
		finish := time.Now()
		duration := finish.Sub(now)
		if duration > time.Second*1 {
			fmt.Println("Time taken", duration, c.Path)
		}
	}()
	switch ctx.Message().(type) {
	case *IncrementRequest:
		c.Value++
		ctx.Response(&IncrementResponse{Value: int32(c.Value)})
	default:
		ctx.Unhandled()
	}
}
