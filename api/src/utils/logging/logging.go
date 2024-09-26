package logging

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
// : Logger
// ------------------------------------------------------------
func NewLogger(name string) *Logger {
	var logger  zerolog.Logger
	var multi   zerolog.LevelWriter
	writerConsole := &zerolog.ConsoleWriter{}
 	writerFile    := &lumberjack.Logger{}

	location, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		panic(err)
	}

	writerConsole.Out             = os.Stdout
	writerConsole.TimeLocation    = location
	writerConsole.FormatTimestamp = func(i interface{}) string {
		return time.Now().Format(time.RFC3339)
	}
	writerConsole.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s=", i)
	}
	writerConsole.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%v", i)
	}
	writerConsole.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	writerFile.Filename   = "./logs/api.log"
	writerFile.MaxSize    = 10
	writerFile.MaxBackups = 100
	writerFile.MaxAge     = 28
	writerFile.Compress   = true

	hostname, _ := os.Hostname()

	multi  = zerolog.MultiLevelWriter(writerConsole, writerFile)
	logger = zerolog.New(multi).With().Caller().Timestamp().Str("hostname", hostname).Logger()

	return &Logger{&logger}
}
