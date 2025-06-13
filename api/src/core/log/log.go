package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"

	"dse/src/utils/directory"

	"github.com/phuslu/log"
	l "github.com/phuslu/log"
)

// ------------------------------------------------------------
// : Constants
// ------------------------------------------------------------
const (
	Reset   = "\x1b[0m"
	Black   = "\x1b[30m"
	Red     = "\x1b[31m"
	Green   = "\x1b[32m"
	Yellow  = "\x1b[33m"
	Blue    = "\x1b[34m"
	Magenta = "\x1b[35m"
	Cyan    = "\x1b[36m"
	White   = "\x1b[37m"
	Gray    = "\x1b[90m"
)
// ------------------------------------------------------------
// : Types
// ------------------------------------------------------------
type Logger struct {
	console l.Logger
	file	l.Logger
}

func Init() {
	l.DefaultLogger.Level      = l.TraceLevel
	l.DefaultLogger.Caller     = 2
	l.DefaultLogger.TimeField  = "ts"
	l.DefaultLogger.TimeFormat = "2006.01.02 15:04:05.000"
	l.DefaultLogger.Writer = &l.MultiLevelWriter{
		ConsoleLevel:  l.TraceLevel,
		InfoWriter:    &l.FileWriter{Filename: "logs/app.log", MaxSize: 1 << 20, EnsureFolder: true, MaxBackups: 1},
		WarnWriter:    &l.FileWriter{Filename: "logs/app.log", MaxSize: 1 << 20, EnsureFolder: true, MaxBackups: 1},
		ErrorWriter:   &l.FileWriter{Filename: "logs/app.log", MaxSize: 1 << 20, EnsureFolder: true, MaxBackups: 1},
		ConsoleWriter: &l.ConsoleWriter{
			ColorOutput: true,
			Formatter: func(out io.Writer, args *l.FormatterArgs) (n int, err error) {
				b := &bytes.Buffer{}

				var color, three string
				switch args.Level {
				case "trace":
					color, three = Magenta, "TRC"
				case "debug":
					color, three = Yellow, "DBG"
				case "info":
					color, three = Green, "INF"
				case "warn":
					color, three = Yellow, "WRN"
				case "error":
					color, three = Red, "ERR"
				case "fatal":
					color, three = Red, "FTL"
				case "panic":
					color, three = Red, "PNC"
				default:
					color, three = Gray, "???"
				}

				fn := args.Get("callerfunc")

				// Time and Level
				fmt.Fprintf(b, "%s%s%s %s%s%s ", Gray, args.Time, Reset, color, three, Reset) 

				// Caller and Function
				fmt.Fprintf(b, "%s (%s%s%s) %s>%s", args.Caller, Cyan, fn, Reset, Cyan, Reset) 

				// Message
				if args.Message != "" {
					fmt.Fprintf(b, " %s", args.Message) // Message
				}

				// Properties
				for _, kv := range args.KeyValues {
					if kv.Key == "callerfunc" {
						continue
					}

					if kv.ValueType == 's' {
						kv.Value = strconv.Quote(kv.Value)
					}
					if kv.Key == "error" && kv.Value != "null" {
						fmt.Fprintf(b, " %s%s=%s%s", Red, kv.Key, kv.Value, Reset)
					} else {
						fmt.Fprintf(b, " %s%s=%s%s%s", Yellow, kv.Key, Gray, kv.Value, Reset)
					}
				}

				// End
				fmt.Fprintf(b, "\n")

				return out.Write(b.Bytes())
			},
		},
	}
}

func Info() *l.Entry {
	return l.DefaultLogger.Info().Str("app", os.Getenv("APP"))
}

func Debug() *l.Entry {
	return l.DefaultLogger.Debug().Str("app", os.Getenv("APP"))
}

func Warn() *l.Entry {
	return l.DefaultLogger.Warn().Str("app", os.Getenv("APP"))
}

func Error() *l.Entry {
	return l.DefaultLogger.Error().Str("app", os.Getenv("APP"))
}

func Fatal() *l.Entry {
	return l.DefaultLogger.Fatal().Str("app", os.Getenv("APP"))
}

func Panic() *l.Entry {
	return l.DefaultLogger.Panic().Str("app", os.Getenv("APP"))
}

func Trace() *l.Entry {
	return l.DefaultLogger.Trace().Str("app", os.Getenv("APP"))
}

func Println(v ...interface{}) {
	log.DefaultLogger.Trace().Any("message", fmt.Sprint(v...)).Msg("")
}

func PrintStack() {
	stack  := make([]uintptr, 10)
	length := runtime.Callers(2, stack[:])
	frames := runtime.CallersFrames(stack[:length])

	cd, err := directory.GetCurrentDirectory()
	if err != nil {
		fmt.Printf("error getting current directory: %v\n", err)
		return
	}

	cd = strings.ReplaceAll(cd, "\\", "/")

	for {
		frame, more := frames.Next()
		// framefn     := frame.Function
		frameln     := frame.Line
		framefile   := frame.File
		framefile   = strings.ReplaceAll(framefile, "\\", "/")

		if strings.HasPrefix(framefile, cd) {
			framefile = strings.TrimPrefix(framefile, cd)
			fmt.Printf("fn: %s%-30s%s ln: %s%d%s\n", Yellow, framefile, Reset, Blue, frameln, Reset)
		}

		if !more {
			fmt.Printf("\n")
			break
		}
	}
}
