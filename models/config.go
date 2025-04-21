package models

// ServerConfig 服务器配置
type ServerConfig struct {
	Port    string `json:"port"`
	Version string `json:"version"`
	Name    string `json:"name"`
}
