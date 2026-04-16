package core_logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// структура самого логгера
type Logger struct {
	*zap.Logger
	//встраиваем все метоы *zap.Logger

	file *os.File
	//файл куда будут писаться файлы
}

// функция для проброса метода через контекст
func FromContext(ctx context.Context) *Logger {
	//получаем из контекса по ключу "log" наш логер
	log, ok := ctx.Value("log").(*Logger)
	if !ok {
		panic("no logger in context")
	}
	//если логера нет, то вызываем панику (потому что нет смысла продолжать без логера)
	// а если есть, то просто возвращаем логер

	return log
}

// конструктор самого логера
func NewLogger(config Config) (*Logger, error) {

	//определяем уровень логированя, берем ее из конфига, откуда мы раннее его передавали
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	//cоздаем директорию, которую также передавали в config
	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("mkdir log folder: %w", err)
	}

	//создаем путь до файла
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(
		config.Folder,
		fmt.Sprintf("%s.log", timestamp),
	)

	//открываем наш файл, если его нет, то создаем через флаг os.O_CREATE
	// и еще выдаем права
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	// настраиваем визуальную часть а именно, чтобы текст был цветным и формат времени
	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	// чтобы был формат не JSON а консольный
	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	// собираем через NewTee, потому что пишем и в терминал и в файл
	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),
	)

	// создаем сам логгер прикручиваем к нему наше ядро и так же чтобы показывал строку и файл
	zapLogger := zap.New(core, zap.AddCaller())

	// возвращаем логгер
	return &Logger{
		Logger: zapLogger,
		file:   logFile,
	}, nil
}

// метод чтобы добавлять к логеру еще какую-то информацию, где field это формат ключ значение
func (l *Logger) With(field ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(field...),
		file:   l.file,
	}
}

// закрываем файл чтобы не утекали данные
func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close application logger:", err)
	}
}
