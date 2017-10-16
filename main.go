package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {

	// We create an errgroup with a common context. If one of the goroutines
	// returns with an error the context will be canceled and all the other
	// goroutines will receive a signal that they should stop working.
	eg, ctx := errgroup.WithContext(context.Background())

	for i := 0; i < 8; i++ {
		i := i
		eg.Go(func() error {
			return startWorker(ctx, i)
		})
	}

	eg.Go(func() error {
		return listenForInterrupt(ctx)
	})

	if err := eg.Wait(); err != nil {
		log.Printf("We stopped because of the following err: %s", err)
	}
	log.Print("All goroutines are stopped")
}

func startWorker(ctx context.Context, i int) error {
	log.Printf("worker %d: started", i)
	for {

		// check if we should stop working
		select {
		case <-ctx.Done():
			log.Printf("worker%d: stopped", i)
			// cleanup can be done here
			return nil
		default:
		}

		// do some work
		time.Sleep(500 * time.Millisecond)
	}
}

func listenForInterrupt(ctx context.Context) error {
	log.Print("listenForInterrupt: started")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// this will block until we receive on one of the channels
	// because we don't have a default case.
	select {
	case <-ctx.Done():
		// this case won't happen in this example but could be triggererd
		// in the future if one of the worker goroutines stops with an error.
		log.Print("listenForInterupt: some other goroutine canceled the context")
		return nil
	case <-signals:
		log.Print("listenForInterrupt: received interrupt")
		return errors.New("Received Interrupt")
	}
}
