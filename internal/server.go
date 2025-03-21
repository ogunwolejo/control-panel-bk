package internal

import (
	"context"
	"control-panel-bk/config"
	"control-panel-bk/internal/aws"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)


func ControlPanelServer() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Load AWS Configuration
	go func() {
		if err := config.LoadAwsConfiguration(); err != nil {
			log.Fatalln(err)
		}
	}()

	if _, err := aws.ConnectMongoDB(); err != nil {
		log.Print("Error MongoDB: ", err)
	}

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
