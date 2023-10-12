package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/indig0fox/a3go/a3interface"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var ll *lumberjack.Logger
var armaWriter *armaIoWriter

var Log, FileOnly, ArmaOnly zerolog.Logger
var ActiveOptions *LoggerOptionsType = &LoggerOptionsType{}

type LoggerOptionsType struct {
	// LogPath is the path to the log file
	Path string
	// LogAddonName is the name of the addon that will be used to send log messages to arma
	AddonName string
	// LogExtensionName is the name of the extension that will be used to send log messages to arma
	ExtensionName string
	// ExtensionVersion is the version of this extension
	ExtensionVersion string
	// LogDebug determines if we should send Debug level messages to file & arma
	Debug bool
	// LogTrace is used to determine if file should receive trace level, regardless of debug
	Trace bool
}

func RotateLogs() {
	ll.Rotate()
}

// ArmaIoWriter is a custom type that implements the io.Writer interface and sends the output to Arma with the "log" callback
type armaIoWriter struct{}

func (w *armaIoWriter) Write(p []byte) (n int, err error) {
	// write to arma log
	a3interface.WriteArmaCallback(ActiveOptions.ExtensionName, ":LOG:", string(p))
	return len(p), nil
}

// console writer
func InitLoggers(o *LoggerOptionsType) {
	ActiveOptions = o

	// create a new lumberjack file logger (adds log rotation and compression)
	ll = &lumberjack.Logger{
		Filename:   ActiveOptions.Path,
		MaxSize:    5,
		MaxBackups: 8,
		MaxAge:     14,
		Compress:   false,
		LocalTime:  true,
	}

	// create a new io writer using the a3go callback function
	// this will be used to write to the arma log
	armaWriter = new(armaIoWriter)

	// create format functions for RPT log messages
	armaLogFormatLevel := func(i interface{}) string {
		return strings.ToUpper(
			fmt.Sprintf(
				"%s:",
				i,
			))
	}
	armaLogFormatTimestamp := func(i interface{}) string {
		return ""
	}

	FileOnly = zerolog.New(zerolog.ConsoleWriter{
		Out:        ll,
		TimeFormat: time.RFC3339,
		NoColor:    true,
	}).With().Timestamp().Caller().Logger()

	if ActiveOptions.Trace {
		FileOnly = FileOnly.Level(zerolog.TraceLevel)
	} else if ActiveOptions.Debug {
		FileOnly = FileOnly.Level(zerolog.DebugLevel)
	} else {
		FileOnly = FileOnly.Level(zerolog.InfoLevel)
	}

	ArmaOnly = zerolog.New(zerolog.ConsoleWriter{
		Out:             armaWriter,
		TimeFormat:      "",
		NoColor:         true,
		FormatLevel:     armaLogFormatLevel,
		FormatTimestamp: armaLogFormatTimestamp,
	}).With().Str("extension_version", ActiveOptions.ExtensionVersion).Logger()

	if ActiveOptions.Debug {
		ArmaOnly = ArmaOnly.Level(zerolog.DebugLevel)
	} else {
		ArmaOnly = ArmaOnly.Level(zerolog.InfoLevel)
	}

	// create something that can send the same message to both loggers
	// this is used to send messages to the arma log
	// and the file log
	Log = zerolog.New(zerolog.MultiLevelWriter(
		zerolog.ConsoleWriter{
			Out:        ll,
			TimeFormat: time.RFC3339,
			NoColor:    true,
		},
		zerolog.ConsoleWriter{
			Out:             armaWriter,
			TimeFormat:      "",
			NoColor:         true,
			FormatTimestamp: armaLogFormatTimestamp,
			FormatLevel:     armaLogFormatLevel,
			FieldsExclude:   []string{zerolog.CallerFieldName, "ctx"},
		},
	)).With().Timestamp().Caller().Logger()

	if ActiveOptions.Debug {
		Log = Log.Level(zerolog.DebugLevel)
	} else {
		Log = Log.Level(zerolog.InfoLevel)
	}
	if ActiveOptions.Trace {
		Log = Log.Level(zerolog.TraceLevel)
	}

}
