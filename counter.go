package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/tochemey/goakt/v3/actor"
)

type Counter struct {
	Id    string
	Value int
}

const statePath = "./state/%s.txt"

var _ actor.Grain = (*Counter)(nil)

func NewCounter(id string) actor.GrainFactory {
	return func(ctx context.Context) (actor.Grain, error) {
		return &Counter{Id: id}, nil
	}
}

func (c *Counter) OnActivate(ctx context.Context, props *actor.GrainProps) error {
	fmt.Println("Activating counter", c.Id)
	value, err := readIntFromFile(c.Id)
	if err != nil {
		c.Value = 0
	}
	c.Value = value
	return nil
}

func (c *Counter) OnDeactivate(ctx context.Context, props *actor.GrainProps) error {
	fmt.Printf("Deactivating counter ID: %s, value: %d. Persisting state.\n", c.Id, c.Value)
	writeIntToFile(c.Id, c.Value)
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

// writeIntToFile writes an integer to a file as a string
func writeIntToFile(filename string, value int) error {
	path := fmt.Sprintf(statePath, filename)
	fmt.Println("Writing to file", path)
	content := strconv.Itoa(value)
	return os.WriteFile(path, []byte(content), 0644)
}

// readIntFromFile reads an integer from a file
func readIntFromFile(filename string) (int, error) {
	path := fmt.Sprintf(statePath, filename)
	fmt.Println("Reading from file", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data))
}
