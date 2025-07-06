package main

import (
	"fmt"
	"net/http"
	"context"
	"time"

	"github.com/tochemey/goakt/v3/actor"

)

func main() {
	actorSystem, _ := actor.NewActorSystem("PageCounterActorSystem",
		actor.WithActorInitMaxRetries(3),
	)

	actorSystem.Start(context.Background())

	http.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		path := r.PathValue("id")
		grain, _ := actorSystem.GrainIdentity(ctx, path, NewCounter(path))
		message := &IncrementRequest{}

		res, err := actorSystem.AskGrain(ctx, grain, message, time.Millisecond*10)
		if err != nil {
			fmt.Println("error", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		count, ok := res.(*IncrementResponse)
		if !ok {
			fmt.Println("invalid response, failed to cast to *IncrementResponse")
			http.Error(w, "invalid response", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "page count: %d", count.Value)
	})
	http.ListenAndServe(":8089", nil)
}
