package log

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	FmtConsole = "console"
	FmtJson    = "json"
)

var (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	PanicLevel = zapcore.PanicLevel
	FatalLevel = zapcore.FatalLevel
)

var (
	Any         = zap.Any
	Array       = zap.Array
	Object      = zap.Object
	Binary      = zap.Binary
	Bool        = zap.Bool
	Bools       = zap.Bools
	ByteString  = zap.ByteString
	ByteStrings = zap.ByteStrings
	Complex64   = zap.Complex64
	Complex64s  = zap.Complex64s
	Complex128  = zap.Complex128
	Complex128s = zap.Complex128s
	Duration    = zap.Duration
	Durations   = zap.Durations
	Err         = zap.Error
	Errors      = zap.Errors
	Float32     = zap.Float32
	Float32s    = zap.Float32s
	Float64     = zap.Float64
	Float64s    = zap.Float64s
	Int         = zap.Int
	Ints        = zap.Ints
	Int8        = zap.Int8
	Int8s       = zap.Int8s
	Int16       = zap.Int16
	Int16s      = zap.Int16s
	Int32       = zap.Int32
	Int32s      = zap.Int32s
	Int64       = zap.Int64
	Int64s      = zap.Int64s
	Namespace   = zap.Namespace
	Reflect     = zap.Reflect
	Stack       = zap.Stack
	StackSkip   = zap.StackSkip
	String      = zap.String
	Stringer    = zap.Stringer
	Strings     = zap.Strings
	Time        = zap.Time
	Times       = zap.Times
	Uint        = zap.Uint
	Uints       = zap.Uints
	Uint8       = zap.Uint8
	Uint8s      = zap.Uint8s
	Uint16      = zap.Uint16
	Uint16s     = zap.Uint16s
	Uint32      = zap.Uint32
	Uint32s     = zap.Uint32s
	Uint64      = zap.Uint64
	Uint64s     = zap.Uint64s
	Uintptr     = zap.Uintptr
	Uintptrs    = zap.Uintptrs
)

var (
	std = New()
	mu  sync.Mutex
)

var _ Logger = &logger{}

func Init(options ...Option) {
	mu.Lock()
	defer mu.Unlock()
	std = New(options...)
}

type Field = zapcore.Field

type Level = zapcore.Level

type Logger interface {
	Log(level int, msg string, fields ...Field)
	Logf(level int, format string, v ...interface{})
	Logw(level int, msg string, keysAndValues ...interface{})

	Info(msg string, fields ...Field)
	Infof(format string, v ...interface{})
	Infow(msg string, keysAndValues ...interface{})

	Debug(msg string, fields ...Field)
	Debugf(format string, v ...interface{})
	Debugw(msg string, keysAndValues ...interface{})

	Warn(msg string, fields ...Field)
	Warnf(format string, v ...interface{})
	Warnw(msg string, keysAndValues ...interface{})

	Error(msg string, fields ...Field)
	Errorf(format string, v ...interface{})
	Errorw(msg string, keysAndValues ...interface{})

	Panic(msg string, fields ...Field)
	Panicf(format string, v ...interface{})
	Panicw(msg string, keysAndValues ...interface{})

	Fatal(msg string, fields ...Field)
	Fatalf(format string, v ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	WithValues(keysAndValues ...interface{}) *logger

	Flush()
}
type Option func(*logger)

type logger struct {
	zapLogger         *zap.Logger
	OutputPaths       []string `json:"output-paths"       mapstructure:"output-paths"`
	ErrorOutputPaths  []string `json:"error-output-paths" mapstructure:"error-output-paths"`
	Level             Level    `json:"level"              mapstructure:"level"`
	Format            string   `json:"format"             mapstructure:"format"`
	DisableCaller     bool     `json:"disable-caller"     mapstructure:"disable-caller"`
	DisableStacktrace bool     `json:"disable-stacktrace" mapstructure:"disable-stacktrace"`
	Development       bool     `json:"development"        mapstructure:"development"`
	Name              string
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func milliSecondsDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Millisecond))
}

func New(options ...Option) *logger {
	logger := &logger{
		Level:             InfoLevel,
		DisableCaller:     false,
		DisableStacktrace: false,
		Format:            FmtConsole,
		Development:       true,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		Name:              "",
	}

	for _, option := range options {
		option(logger)
	}
	encodeLevel := zapcore.CapitalLevelEncoder

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     timeEncoder,
		EncodeDuration: milliSecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	zapConfig := &zap.Config{
		Level:             zap.NewAtomicLevelAt(logger.Level),
		Development:       logger.Development,
		DisableCaller:     logger.DisableCaller,
		DisableStacktrace: logger.DisableStacktrace,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         logger.Format,
		EncoderConfig:    encoderConfig,
		OutputPaths:      logger.OutputPaths,
		ErrorOutputPaths: logger.ErrorOutputPaths,
	}

	zapLogger, err := zapConfig.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	if logger.Name != "" {
		zapLogger = zapLogger.Named(logger.Name)
	}

	logger.zapLogger = zapLogger

	return logger
}

func OptLevel(level Level) Option {
	return func(l *logger) {
		l.Level = level
	}
}

func OptDisableCaller(disableCaller bool) Option {
	return func(l *logger) {
		l.DisableCaller = disableCaller
	}
}

func OptDisableStacktrace(disableStacktrace bool) Option {
	return func(l *logger) {
		l.DisableStacktrace = disableStacktrace
	}
}

func OptFormat(format string) Option {
	return func(l *logger) {
		l.Format = format
	}
}

func OptOutputPaths(outputPaths []string) Option {
	return func(l *logger) {
		l.OutputPaths = outputPaths
	}
}

func OptErrorOutputPaths(errorOutputPaths []string) Option {
	return func(l *logger) {
		l.ErrorOutputPaths = errorOutputPaths
	}
}

func OptName(name string) Option {
	return func(l *logger) {
		l.Name = name
	}
}

func WithValues(keysAndValues ...interface{}) *logger {
	return std.WithValues(keysAndValues...)
}

func (l *logger) WithValues(keysAndValues ...interface{}) *logger {
	newZapLogger := l.zapLogger.With(handleFields(l.zapLogger, keysAndValues)...)
	l.zapLogger = newZapLogger

	return l
}

func Flush() {
	std.Flush()
}

func (l *logger) Flush() {
	_ = l.zapLogger.Sync()
}

func (l *logger) Enabled() bool { return true }

func Log(level int, msg string, fields ...Field) {
	std.Log(level, msg, fields...)
}

func (l *logger) Log(level int, msg string, fields ...Field) {
	if checkedEntry := l.zapLogger.Check(zapcore.Level(-1*level), msg); checkedEntry != nil {
		checkedEntry.Write(fields...)
	}
}

func Logf(level int, format string, args ...interface{}) {
	std.Logf(level, format, args...)
}

func (l *logger) Logf(level int, format string, args ...interface{}) {
	if checkedEntry := l.zapLogger.Check(zapcore.Level(-1*level), fmt.Sprintf(format, args...)); checkedEntry != nil {
		checkedEntry.Write()
	}
}

func Logw(level int, msg string, keysAndValues ...interface{}) {
	std.Logw(level, msg, keysAndValues...)
}

func (l *logger) Logw(level int, msg string, keysAndValues ...interface{}) {
	if checkedEntry := l.zapLogger.Check(zapcore.Level(-1*level), msg); checkedEntry != nil {
		checkedEntry.Write(handleFields(l.zapLogger, keysAndValues)...)
	}
}

func ErrorLog(err error, fields ...Field) {
	if !std.DisableStacktrace {
		fields = append(fields, zap.StackSkip("stack", 1))
	}
	std.zapLogger.Error(err.Error(), fields...)
}

func (l *logger) ErrorLog(err error, fields ...Field) {
	l.zapLogger.Error(err.Error(), fields...)
}

func ErrorLogw(err error, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Errorw(err.Error(), keysAndValues...)
}

func (l *logger) ErrorLogw(err error, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Errorw(err.Error(), keysAndValues...)
}

// Debug method output debug level log.
func Debug(msg string, fields ...Field) {
	std.zapLogger.Debug(msg, fields...)
}

func (l *logger) Debug(msg string, fields ...Field) {
	l.zapLogger.Debug(msg, fields...)
}

// Debugf method output debug level log.
func Debugf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Debugf(format, v...)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Debugf(format, v...)
}

// Debugw method output debug level log.
func Debugw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Debugw(msg, keysAndValues...)
}

func (l *logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Debugw(msg, keysAndValues...)
}

// Info method output info level log.
func Info(msg string, fields ...Field) {
	std.zapLogger.Info(msg, fields...)
}

func (l *logger) Info(msg string, fields ...Field) {
	l.zapLogger.Info(msg, fields...)
}

// Infof method output info level log.
func Infof(format string, v ...interface{}) {
	std.zapLogger.Sugar().Infof(format, v...)
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.zapLogger.Sugar().Infof(format, v...)
}

// Infow method output info level log.
func Infow(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Infow(msg, keysAndValues...)
}

func (l *logger) Infow(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Infow(msg, keysAndValues...)
}

// Warn method output warning level log.
func Warn(msg string, fields ...Field) {
	std.zapLogger.Warn(msg, fields...)
}

func (l *logger) Warn(msg string, fields ...Field) {
	l.zapLogger.Warn(msg, fields...)
}

// Warnf method output warning level log.
func Warnf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Warnf(format, v...)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Warnf(format, v...)
}

// Warnw method output warning level log.
func Warnw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Warnw(msg, keysAndValues...)
}

func (l *logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Warnw(msg, keysAndValues...)
}

// Error method output error level log.
func Error(msg string, fields ...Field) {
	std.zapLogger.Error(msg, fields...)
}

func (l *logger) Error(msg string, fields ...Field) {
	l.zapLogger.Error(msg, fields...)
}

// Errorf method output error level log.
func Errorf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Errorf(format, v...)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Errorf(format, v...)
}

// Errorw method output error level log.
func Errorw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Errorw(msg, keysAndValues...)
}

func (l *logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Errorw(msg, keysAndValues...)
}

// Panic method output panic level log and shutdown application.
func Panic(msg string, fields ...Field) {
	std.zapLogger.Panic(msg, fields...)
}

func (l *logger) Panic(msg string, fields ...Field) {
	l.zapLogger.Panic(msg, fields...)
}

// Panicf method output panic level log and shutdown application.
func Panicf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Panicf(format, v...)
}

func (l *logger) Panicf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Panicf(format, v...)
}

// Panicw method output panic level log.
func Panicw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Panicw(msg, keysAndValues...)
}

func (l *logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Panicw(msg, keysAndValues...)
}

// Fatal method output fatal level log.
func Fatal(msg string, fields ...Field) {
	std.zapLogger.Fatal(msg, fields...)
}

func (l *logger) Fatal(msg string, fields ...Field) {
	l.zapLogger.Fatal(msg, fields...)
}

// Fatalf method output fatal level log.
func Fatalf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Fatalf(format, v...)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Fatalf(format, v...)
}

// Fatalw method output Fatalw level log.
func Fatalw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Fatalw(msg, keysAndValues...)
}

func (l *logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Fatalw(msg, keysAndValues...)
}

func handleFields(l *zap.Logger, args []interface{}, additional ...zap.Field) []zap.Field {
	// a slightly modified version of zap.SugaredLogger.sweetenFields
	if len(args) == 0 {
		// fast-return if we have no suggared fields.
		return additional
	}

	// unlike Zap, we can be pretty sure users aren't passing structured
	// fields (since logr has no concept of that), so guess that we need a
	// little less space.
	fields := make([]zap.Field, 0, len(args)/2+len(additional))
	for i := 0; i < len(args); {
		// check just in case for strongly-typed Zap fields, which is illegal (since
		// it breaks implementation agnosticism), so we can give a better error message.
		if _, ok := args[i].(zap.Field); ok {
			l.DPanic("strongly-typed Zap Field passed to logr", zap.Any("zap field", args[i]))

			break
		}

		// make sure this isn't a mismatched key
		if i == len(args)-1 {
			l.DPanic("odd number of arguments passed as key-value pairs for logging", zap.Any("ignored key", args[i]))

			break
		}

		// process a key-value pair,
		// ensuring that the key is a string
		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			// if the key isn't a string, DPanic and stop logging
			l.DPanic(
				"non-string key argument passed to logging, ignoring all later arguments",
				zap.Any("invalid key", key),
			)

			break
		}

		fields = append(fields, zap.Any(keyStr, val))
		i += 2
	}

	return append(fields, additional...)
}
