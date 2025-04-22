package models

import "time"

type Configs struct {
	LogParams      LogParams      `json:"log_params"`
	AppParams      AppParams      `json:"app_params"`
	PostgresParams PostgresParams `json:"postgres_params"`
	RedisParams    RedisParams    `json:"redis_params"`
	RabbitParams   RabbitParams   `json:"rabbit_params"`
	Auth           Auth           `json:"auth"`
}

type LogParams struct {
	LogDirectory     string `json:"log_directory"`
	LogInfo          string `json:"log_info"`
	LogError         string `json:"log_error"`
	LogWarn          string `json:"log_warn"`
	LogDebug         string `json:"log_debug"`
	MaxSizeMegabytes int    `json:"max_size_megabytes"`
	MaxBackups       int    `json:"max_backups"`
	MaxAge           int    `json:"max_age"`
	Compress         bool   `json:"compress"`
	LocalTime        bool   `json:"local_time"`
}

type AppParams struct {
	GinMode    string `json:"gin_mode"`
	ServerURL  string `json:"server_url"`
	ServerName string `json:"server_name"`
	PortRun    string `json:"port_run"`
}

type PostgresParams struct {
	User     string `json:"user"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	SSLMode  string `json:"sslmode"`
}

type RedisParams struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Password   string `json:"password"`
	DB         int    `json:"db"`
	TTLMinutes int    `json:"ttl_minutes"`
}

type RabbitParams struct {
	URLConn string `json:"url_conn"`
}

type Auth struct {
	JwtSecretKey  string        `json:"jwt_secret_key"`
	JwtTtlMinutes time.Duration `json:"jwt_ttl_minutes"`
}
