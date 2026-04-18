package core_http_server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

/*
Тут прописан конфиг нашего сервера

Addr - это порт
ShutdownTimeout -  это время которое мы даем серверу на закрытие
*/

type Config struct {
	// порт и время задаем через переменные окружения
	// поэтому используем envconfig
	Addr            string        `envconfig:"ADDR" required:"true"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" deafult:"30s"`
}

// Конструктор конфига, а именно подгружаем переменные с тегом HTTP
func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("HTTP", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

// Конфиг который паникует при ошибке
func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get HTTP server config: %w", err)
		panic(err)
	}
	return config
}
