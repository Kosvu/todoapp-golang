package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Создаем структуру конфига для логгера
type Config struct {
	Level string `envconfig:"LEVEL" default:"DEBUG"`
	// задаем уровень который берем из переменной окружения
	//required:"true" значит что это ]2 обязательно
	Folder string `envconfig:"FOLDER" required:"true"`
	// задаем директорию которюу тоже берем из переменной окружения
	//required:"true" значит что это ]2 обязательно
}

// конструктор создания конфига
func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("LOGGER", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
		//возвращаем в случае ошибки пустой конфиг и обернутую ошибку
	}
	//грузим переменные окружения. задаем префикс LOGGER
	//значит при создании нужно будет писать переменную не как LEVEL, а как LOGGER_LEVEL

	return config, nil
	//возвращаем если нет ошибки
}

// конфиг который должен быть, иначе паника
func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get Logger config: %w", err)
		panic(err)
	}
	// если нет логера, нет смысла прожолжать, паникуем

	return config
	// если все хорошо, то продолжаем
}
