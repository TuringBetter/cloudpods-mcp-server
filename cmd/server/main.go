package main

import (
	"cloudpods-mcp-server/pkg/api"
	"cloudpods-mcp-server/pkg/config"
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger(debug bool) {
	var cfg zap.Config
	if debug {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
	}
	logger, err := cfg.Build()
	if err != nil {
		fmt.Printf("fail to initialize logger: %v\n", err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(logger)
}

func main() {
	configPath := flag.String("config",
		"cloudpods-mcp-server.yaml",
		"path to the config file (default: cloudpods-mcp-server.yaml in the current directory)")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("fail to load config: %v\n\n", err)
		os.Exit(1)
	}

	initLogger(cfg.Debug)
	defer zap.L().Sync() //nolint:errcheck

	zap.L().Info("cloudpods mcp server starting...",
		zap.String("address", cfg.Address),
		zap.Int("port", cfg.Port),
		zap.Bool("debug", cfg.Debug),
	)

	server, err := api.NewMCPServer(cfg)
	if err != nil {
		zap.L().Error("fail to create a cloudpods mcp server", zap.Error(err))
		os.Exit(1)
	}

	if err = server.Start(); err != nil {
		zap.L().Error("fail to start mcp server", zap.Error(err))
		os.Exit(1)
	}
}
