package internal

import (
	"context"
	"control-panel-bk/internal/database"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Database interface {
	error
}

func ControlPanelServer() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Connect to various database
	database.InitializeAllDbs(database.Dbs)

	server := &http.Server{
		Handler: routes(),
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
	}

	go func() {
		log.Printf("Server started on port %s\n", os.Getenv("PORT"))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the engine
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server exited gracefully")
	os.Exit(0)
}
