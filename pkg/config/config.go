package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Address string
	Port    int
	Debug   bool

	CloudpodsAPI string
	Username     string
	Password     string
	Domain       string
	Project      string

	EnableMonitoring bool
	EnableCreation   bool
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{
		Address: "0,0,0,0",
		Port:    9990,
		Debug:   false,
	}

	// 从配置文件中读取
	if configPath == "" {
		return nil, fmt.Errorf("configpath cannot be empty")
	}

	viper.SetConfigFile(configPath)
	if err := viper.ReadConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}

	// 将配置文件映射到对象中
	config.Address = viper.GetString("server.address")
	config.Port = viper.GetInt("server.port")
	config.Debug = viper.GetBool("server.debug")

	config.CloudpodsAPI = viper.GetString("cloudpods.api_url")
	config.Username = viper.GetString("cloudpods.username")
	config.Password = viper.GetString("cloudpods.password")
	config.Domain = viper.GetString("cloudpods.domain")
	config.Project = viper.GetString("cloudpods.project")

	config.EnableMonitoring = viper.GetBool("cloudpods.enable_monitoring")
	config.EnableCreation = viper.GetBool("cloudpods.enable_creation")

	// 从环境变量覆盖
	if addr := os.Getenv("MCP_SERVER_ADDRESS"); addr != "" {
		config.Address = addr
	}

	if portStr := os.Getenv("MCP_SERVER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			config.Port = port
		}
	}

	if debug := os.Getenv("MCP_SERVER_DEBUG"); debug == "true" {
		config.Debug = true
	}

	if apiURL := os.Getenv("CP_API_URL"); apiURL != "" {
		config.CloudpodsAPI = apiURL
	}

	if username := os.Getenv("CP_USERNAME"); username != "" {
		config.Username = username
	}

	if password := os.Getenv("CP_PASSWORD"); password != "" {
		config.Password = password
	}

	// 验证必要配置
	if config.CloudpodsAPI == "" {
		return nil, fmt.Errorf("cloudpods API URL未配置")
	}

	if config.Username == "" || config.Password == "" {
		return nil, fmt.Errorf("cloudpods凭据未配置")
	}

	return config, nil
}

/*
func main() {
	configPath := flag.String("config",
		"D:\\workspace\\study\\goproject\\src\\cloudpods-mcp-server\\cloudpods-mcp-server.yaml",
		"config file path")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("fail to load config:%v\n\n", err)
		os.Exit(1)
	}

}
*/
