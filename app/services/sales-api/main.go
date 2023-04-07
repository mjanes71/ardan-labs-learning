package main

import (
	"fmt"
	"os"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// var build = "develop"

func main() {
	log, err := initLogger("SALES-API")
	if err != nil {
		fmt.Println("Error constructing logger:", err)
		os.Exit(1)
	}

	defer log.Sync()

	// Perform the startup and shutdown sequence
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	return nil
}

// initLogger constructs a Sugared Logger that writes to stdout and
// provides human-readable timestamps.
func initLogger(service string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]any{
		"service": service,
	}

	config.OutputPaths = []string{"stdout"}

	log, err := config.Build(zap.WithCaller(true))
	if err != nil {
		return nil, err
	}

	return log.Sugar(), nil
}