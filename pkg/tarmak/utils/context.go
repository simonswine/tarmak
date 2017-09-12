package utils

import (
	"context"
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
)

func GetContext() (context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		select {
		case <-ch:
			log.Info("caught signal")
		case <-ctx.Done():
			log.Info("context done")
		}
		signal.Stop(ch)
		log.Info("cancelling")
		cancel()
	}()
	return ctx, cancel
}
