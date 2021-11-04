package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var build = "develop"

func main() {

	// construct the application logger
	log, err := initLogger("SALES-API")
	if err != nil {
		fmt.Println(err)
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

func initLogger(service string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]interface{}{
		"service": "SALES_API",
	}

	log, err := config.Build()
	if err != nil {
		fmt.Println("Error constrcuting logger:", err)
		os.Exit(1)
	}

	return log.Sugar(), nil
}

// if _, err := maxprocs.Set(); err != nil {
// 	fmt.Println("maxprocs: %w", err)
// 	os.Exit(1)
// }

// g := runtime.GOMAXPROCS(0)

// log.Printf("starting service build[%s] CPU[%d]", build, g)
// defer log.Println("service ended")

// shutdown := make(chan os.Signal, 1)
// signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
// <-shutdown

// log.Println("stopping service")
