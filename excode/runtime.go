package excode

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Runner interface {
	Serve(ctx context.Context) error
	Clean(ctx context.Context) error
}

type LivenessChecker interface {
	Liveness(ctx context.Context) error
}

type ReadinessChecker interface {
	Readiness(ctx context.Context) error
}

type RunConfig struct {
	HealthCheckAddress string
}

func alwaysOkHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func livenessHandler(r interface{}) http.HandlerFunc {
	checker, ok := r.(LivenessChecker)
	if !ok {
		return alwaysOkHandler
	}

	return func(w http.ResponseWriter, r *http.Request) {
		err := checker.Liveness(r.Context())
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))

			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func readinessHandler(r interface{}) http.HandlerFunc {
	checker, ok := r.(ReadinessChecker)
	if !ok {
		return alwaysOkHandler
	}

	return func(w http.ResponseWriter, r *http.Request) {
		err := checker.Readiness(r.Context())
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))

			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func runHealthChecker(ctx context.Context, runner Runner, config RunConfig) {
	mux := http.NewServeMux()

	mux.HandleFunc("/live", livenessHandler(runner))
	mux.HandleFunc("/ready", readinessHandler(runner))

	log.Println("Run health checker...")

	server := http.Server{
		Addr:    config.HealthCheckAddress,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancelShutdownCtx := context.WithTimeout(context.Background(), 30*time.Second)

		go forceExitIfContextTimeout(shutdownCtx, "Health checker shutdown timed out... Forcing exit now...")

		log.Println("Shutdown health checker...")
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(fmt.Errorf("health checker shutdown error: %w", err))
		}

		cancelShutdownCtx()
	}()

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func runServer(ctx context.Context, runner Runner) {
	log.Println("Run server...")

	err := runner.Serve(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func runCleaner(ctx context.Context, runner Runner) {
	log.Println("Clean server...")

	err := runner.Clean(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func forceExitIfContextTimeout(ctx context.Context, msg string) {
	<-ctx.Done()
	if ctx.Err() == context.DeadlineExceeded {
		log.Fatal(msg)
	}
}

// Run is a helper function to run the server.
// It will handle graceful shutdown and log errors.
// It will also handle liveness & readiness probe.
// Keep in note that any error will result in os.Exit(1).
func Run(ctx context.Context, runner Runner, config RunConfig) {
	ctx, cancel := context.WithCancel(ctx)

	if config.HealthCheckAddress == "" {
		config.HealthCheckAddress = ":8086"
	}

	// Run health checker in separate goroutine
	go runHealthChecker(ctx, runner, config)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		cleanCtx, cancelCleanCtx := context.WithTimeout(ctx, 30*time.Second)

		go forceExitIfContextTimeout(cleanCtx, "Graceful shutdown timed out... Forcing exit now...")

		runCleaner(cleanCtx, runner) // Run cleaner
		cancelCleanCtx()             // Cancel clean context
		cancel()                     // Cancel context to stop health checker & server
	}()

	// Run server in main goroutine
	runServer(ctx, runner)

	// Wait for context to be done
	<-ctx.Done()
}
