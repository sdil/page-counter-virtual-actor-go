package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tochemey/goakt/v3/actor"
)

func main() {
	actorSystem, _ := actor.NewActorSystem("PageCounterActorSystem",
		actor.WithActorInitMaxRetries(3),
	)
	ctx := context.Background()
	actorSystem.Start(ctx)

	http.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.PathValue("id")
		grain, _ := actorSystem.GrainIdentity(ctx, id, NewCounter(id))
		message := &IncrementRequest{}

		res, err := actorSystem.AskGrain(ctx, grain, message, time.Millisecond*10)
		if err != nil {
			fmt.Println("cannot ask grain", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		count, ok := res.(*IncrementResponse)
		if !ok {
			fmt.Println("cannot cast response to IncrementResponse")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "ID %s, page count: %d", id, count.Value)
	})

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Println("Starting HTTP server on port 8089")
		http.ListenAndServe(":8089", nil)
	}()

	<-shutdown
	fmt.Println("Shutting down...")
	actorSystem.Stop(ctx)
}
