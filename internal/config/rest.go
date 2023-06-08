package config

import "time"

type REST struct {
	Listen     string `env:"REST_LISTEN" envDefault:":8080"`
	APIVersion string `env:"REST_API_VERSION" envDefault:"v1"`

	ReadTimeout   time.Duration `env:"REST_READ_TIMEOUT" envDefault:"300s"`
	WriteTimeout  time.Duration `env:"REST_WRITE_TIMEOUT" envDefault:"300s"`
	HandleTimeout time.Duration `env:"REST_HANDLE_TIMEOUT" envDefault:"300s"`

	PingDelay time.Duration `json:"REST_PING_DELAY" envDefault:"30s"`
}
