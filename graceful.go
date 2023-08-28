package goro

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GracefulShutdown(server *http.Server, timeout time.Duration, logger func(message string, arguments ...any)) {
	wait := make(chan os.Signal, 1)

	signal.Notify(wait, os.Interrupt, syscall.SIGTERM)

	<-wait

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	logger("Shutting down server with %s timeout", timeout)

	if err := server.Shutdown(ctx); err != nil {
		logger("err while shutting down server: %v", err)
	} else {
		logger("Server was shut gracefully")
	}
}
