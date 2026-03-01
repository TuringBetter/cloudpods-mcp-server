package main

import (
	"cloudpods-mcp-server/pkg/api"
	"cloudpods-mcp-server/pkg/config"
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"
)

func main() {
	configPath := flag.String("config",
		"D:\\workspace\\study\\goproject\\src\\cloudpods-mcp-server\\cloudpods-mcp-server.yaml",
		"config file path")
	flag.Parse()
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("fail to load config: %v\n\n", err)
		os.Exit(1)
	}
	/*
		utils.InitLogger(cfg.Debug, "D:\\workspace\\study\\goproject\\src\\cloudpods-mcp-server\\server_log.log")
		defer utils.Logger.Sync()
	*/
	zap.L().Info("cloudpods mcp server starting...",
		zap.String("address", cfg.Address),
		zap.Int("port", cfg.Port),
		zap.Bool("debug", cfg.Debug),
	)
	server, err := api.NewMCPServer(cfg)
	if err != nil {
		zap.L().Error("fail to create a cloudpods mcp server", zap.Error(err))
	}
	err = server.Start()
	if err != nil {
		zap.L().Error("fail to start mcp server")
	}
}
