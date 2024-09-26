package logutil

import (
	"os"
	"sync"
	"time"

	"github.com/golang-module/carbon"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ------------------------------------------------------------
// : Types
// ------------------------------------------------------------
type Logger struct {
	*zerolog.Logger
}

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var dt carbon.Carbon = carbon.NewCarbon()

var f = &lumberjack.Logger{
	Filename:   "logs/main.log",
	MaxBackups: 100,         // Max number of backup log files
	MaxAge:     1,           // Max age in days
	Compress:   false,       // Compress the backup log files
}

var once sync.Once

// ------------------------------------------------------------
// : Logger
// ------------------------------------------------------------
func NewLogger(name string) *Logger {
	once.Do(func() {
		f.Rotate()
	})

	// Create writer for console
	writer := zerolog.ConsoleWriter{
		Out: os.Stderr,
		FormatTimestamp: func(i interface{}) string {
			return time.Now().Format(time.RFC3339)
		},
	}

    writer_multi := zerolog.MultiLevelWriter(writer, f)
    logger       := zerolog.New(writer_multi).With().
		Caller().
		Timestamp().
		Logger()

    return &Logger{&logger}
}
