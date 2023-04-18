package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"net/http"
	"context"

	"github.com/ardanlabs/conf/v3"
	"github.com/ardanlabs/service/foundation/logger"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"github.com/ardanlabs/service/business/web/v1/debug"
	"github.com/mjanes71/ardan-labs-learning/app/services/sales-api/handlers"
	
)
/*
	This is where we put project todos
	Need to figure out timeouts for http service
*/
var build = "develop"
func main() {
	log, err := logger.New("SALES-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}
func run(log *zap.SugaredLogger) error {
	// =========================================================================
	// GOMAXPROCS
	opt := maxprocs.Logger(log.Infof)
	if _, err := maxprocs.Set(opt); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))
	defer log.Infow("shutdown")

	// =========================================================================
	// Configuration

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s,mask"` // or noprint to not print at all
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "copyright information here",
		},
	}

	const prefix = "SALES"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// =========================================================================
	// App Starting

	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)

	// =========================================================================
	// Start Debug Service

	log.Infow("startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)
	// this is an example of an orphan goroutine
	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.StandardLibraryMux()); err != nil {
			log.Errorw("shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "ERROR", err)
		}
	}()

	// =========================================================================
	// Start API Service

	log.Infow("startup", "status", "initializing V1 API support")
	shutdown := make(chan os.Signal, 1)
	// realy incoming signals (sigint and sigterm) to the shutdown channel
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	
	apiMux := handlers.APIMux(handlers.APIMuxConfig{
		Shutdown: shutdown,
		Log: log,
	})
	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	// open a channel called serverErrors. its expecting 1 error to be sent to it
	serverErrors := make(chan error, 1)

	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		// send the error from api.ListenandServe to the serverErrors channel
		// this would catch a server error
		serverErrors <- api.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	select {
		// this is not the ideal case
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

		//this is the ideal case
	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		// ctx will be passed to the 'blocking call' (shutdown below) to give instruction
		// on how long to wait to timeout
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel() // bill bets that if you don't have a defer cancel, you're probably messing up

		// Shutdown gracefully shuts down the server without interrupting any active connections. 
		// Shutdown works by first closing all open listeners, then closing all idle connections, 
		// and then waiting indefinitely for connections to return to idle and then shut down. 
		// If the provided context expires before the shutdown is complete, Shutdown returns the 
		// context's error, otherwise it returns any error returned from closing the Server's underlying Listener(s).
		if err := api.Shutdown(ctx); err != nil {
			// this part only gets executed if you reach the timeout limit defined in ctx
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}
	return nil
}
