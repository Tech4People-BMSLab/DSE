package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

// ------------------------------------------------------------
// : Types
// ------------------------------------------------------------
type Logger struct {
	*zerolog.Logger
}

// ------------------------------------------------------------
// : Logger (General)
// ------------------------------------------------------------
func NewLogger() *Logger {
	var logger  zerolog.Logger
	var multi   zerolog.LevelWriter
	writer_console := &zerolog.ConsoleWriter{}
 	writer_file    := &lumberjack.Logger{}

	location, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		panic(err)
	}

	writer_console.Out             = os.Stdout
	writer_console.TimeLocation    = location
	writer_console.FormatTimestamp = func(i interface{}) string {
		return time.Now().Format(time.RFC3339)
	}
	writer_console.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s=", i)
	}
	writer_console.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%v", i)
	}
	writer_console.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	writer_file.Filename   = "./logs/api.log"
	writer_file.MaxSize    = 100
	writer_file.MaxBackups = 100
	writer_file.MaxAge     = 30
	writer_file.Compress   = true

	multi  = zerolog.MultiLevelWriter(writer_console, writer_file)
	logger = zerolog.New(multi).With().
		Caller().
		Timestamp().
		Logger()

	return &Logger{&logger}
}
