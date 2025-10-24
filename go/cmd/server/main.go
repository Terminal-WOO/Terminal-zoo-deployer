package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ClappFormOrg/AI-CO/go/internal/server"
	"github.com/ClappFormOrg/AI-CO/go/pkg/log"

	"github.com/kelseyhightower/envconfig"
)

const (
	Component       string         = "cmd.server"
	Program         string         = "aico"
	SignalTerminate syscall.Signal = syscall.SIGTERM
)

var (
	Version   string = "dev"
	Commit    string = "none"
	BuildDate string = "unknown"

	logger      log.Logger      = nil
	showVersion bool            = false
	verbose     int             = 0
	stop        chan os.Signal  = make(chan os.Signal, 1)
	wg          *sync.WaitGroup = new(sync.WaitGroup)
)

// init initializes command-line flags for the server application.
// It sets up both short and long versions of flags for version display and verbosity control.
func init() {
	// Define flags with default values
	flag.BoolVar(&showVersion, "V", false, "Print version and exit (short)")
	flag.BoolVar(&showVersion, "version", false, "Print version and exit (long)")
	flag.IntVar(&verbose, "v", 0, "Increase verbosity (can be used up to 2 times)")
	flag.IntVar(&verbose, "verbose", 0, "Increase verbosity (long version)")

	signal.Notify(stop, SignalTerminate, os.Interrupt)
}

func terminate() {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		panic(err) // We should have panic if we cant find our own PID.
	}

	if err := p.Signal(SignalTerminate); err != nil {
		panic(err) // NOTE: This will happen on windows.
	}
}

func main() {
	flag.Parse()

	// If the `-v` or `--version` flag is given print the version and exit.
	if showVersion {
		fmt.Printf("Version: %s commit:%s built:%s\n", Version, Commit, BuildDate)
		os.Exit(0)
	}

	// Clamp verbosity level to max of 2
	if verbose > 2 {
		fmt.Fprintln(os.Stderr, "Warning: Maximum verbosity level is 2. Clamping to 2.")
		verbose = 2
	}

	// Set logging level based on verbosity
	switch verbose {
	case 1:
		logger = NewLogrusLogger(log.LevelInfo)
	case 2:
		logger = NewLogrusLogger(log.LevelDebug)
	default:
		logger = NewLogrusLogger(log.LevelWarn)
	}
	logger = log.NewComponentLogger(logger, Component)

	// Process the environment variables oconfiguration.
	spec := &Specification{Logger: logger}
	if err := envconfig.Process(Program, spec); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to process environment variables: %s\n", err)
		os.Exit(1)
	}

	opts := []server.Option{server.WithLogger(logger)}

	if spec.TLSKeyFile != "" {
		keyFileBytes, err := os.ReadFile(spec.TLSKeyFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: read file %q: %s\n", spec.TLSKeyFile, err)
			os.Exit(1)
		}
		opts = append(opts, server.WithSecretTLSKey(keyFileBytes))

		crtFileBytes, err := os.ReadFile(spec.TLSCertFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: read file %q: %s\n", spec.TLSCertFile, err)
			os.Exit(1)
		}
		opts = append(opts, server.WithSecretTLSCert(crtFileBytes))
	}

	handler, err := server.NewHandler(opts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create server handler: %s\n", err)
		os.Exit(1)
	}

	srv := &http.Server{
		Handler:      cors(handler),
		Addr:         spec.HTTPListenAddress,
		ReadTimeout:  spec.HTTPReadTimeout,
		WriteTimeout: spec.HTTPWriteTimeout,
		IdleTimeout:  spec.HTTPIdleTimeout,
	}

	defer func() {
		// Create a context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), spec.TerminationGracePeriod)
		defer cancel()

		if err := handler.Stop(ctx); err != nil {
			logger.ErrorCtx(ctx, "failed to stop handler", "err", err)
		}

		if err := handler.Close(); err != nil {
			logger.ErrorCtx(ctx, "failed to close handler", "err", err)
		}

		wg.Wait()
		logger.DebugCtx(ctx, "going away...")
	}()

	wg.Go(func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("handler listen and serve", "err", err)
			terminate()
		}
	})

	sig := <-stop
	logger = logger.With("signal", sig, "grace_period", spec.TerminationGracePeriod)
	logger.Warn("received signal, shutting down gracefully...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), spec.TerminationGracePeriod)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.ErrorCtx(ctx, "failed to gracefully shutdown handler server", "err", err)
	}
}
