package log

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// New creates a Logger instance with provided options.
func New(opts ...Opt) Logger {
	options := defaultOpts()
	for _, opt := range opts {
		opt(options)
	}
	labels := buildLabels(parseLabelsFormat(options.LabelsFormat), copyLabels(options.Labels), options.LabelsSeparator)
	return &logger{
		level:      options.MinLevel,
		format:     buildFormat(options.Format, labels.notEmpty()),
		levelNames: buildLevelNames(*options),
		logger:     buildLogger(options),
		labels:     labels,
	}
}

// ByOptions creates a Logger using provided options.
func ByOptions(opts Opts) Logger {
	options := mergeOpts(defaultOpts(), &opts)
	labels := buildLabels(parseLabelsFormat(options.LabelsFormat), copyLabels(options.Labels), options.LabelsSeparator)
	return &logger{
		level:      options.MinLevel,
		format:     buildFormat(options.Format, labels.notEmpty()),
		levelNames: buildLevelNames(*options),
		logger:     buildLogger(options),
		labels:     labels,
	}
}

// ByLevelName creates a Logger with log level name provided and options.
//
// It can be useful when log is created from some config.
//
// If level name is not recognised, then opts.MinLevel is used as a backup.
func ByLevelName(levelName string, opts ...Opt) Logger {
	options := defaultOpts()
	for _, opt := range opts {
		opt(options)
	}

	levelNames := buildLevelNames(*options)
	if options.UpperCase {
		levelName = strings.ToUpper(levelName)
	}

	level := options.MinLevel
	for lvl, name := range levelNames {
		if name == levelName {
			level = lvl
			break
		}
	}

	labels := buildLabels(parseLabelsFormat(options.LabelsFormat), copyLabels(options.Labels), options.LabelsSeparator)
	return &logger{
		level:      level,
		format:     buildFormat(options.Format, labels.notEmpty()),
		levelNames: levelNames,
		logger:     buildLogger(options),
		labels:     labels,
	}
}

func buildLogger(opts *Opts) *log.Logger {
	if opts.Logger != nil {
		// ignoring flags, as logger might be already pre-configured,
		// so just returning logger as is.
		return opts.Logger
	}
	if opts.Writer != nil {
		return log.New(opts.Writer, "", opts.Flags)
	}
	return log.New(os.Stderr, "", opts.Flags)
}

type logger struct {
	level      Level
	format     format
	levelNames map[Level]string
	labels     labels
	logger     *log.Logger
}

func (l *logger) Log(lvl Level, v ...any)            { l.log(normalizeLevel(lvl), v...) }
func (l *logger) Logf(lvl Level, f string, v ...any) { l.logf(normalizeLevel(lvl), f, v...) }
func (l *logger) Trace(v ...any)                     { l.log(LevelTrace, v...) }
func (l *logger) Tracef(f string, v ...any)          { l.logf(LevelTrace, f, v...) }
func (l *logger) Debug(v ...any)                     { l.log(LevelDebug, v...) }
func (l *logger) Debugf(f string, v ...any)          { l.logf(LevelDebug, f, v...) }
func (l *logger) Info(v ...any)                      { l.log(LevelInfo, v...) }
func (l *logger) Infof(f string, v ...any)           { l.logf(LevelInfo, f, v...) }
func (l *logger) Notice(v ...any)                    { l.log(LevelInfo, v...) }
func (l *logger) Noticef(f string, v ...any)         { l.logf(LevelInfo, f, v...) }
func (l *logger) Warn(v ...any)                      { l.log(LevelWarn, v...) }
func (l *logger) Warnf(f string, v ...any)           { l.logf(LevelWarn, f, v...) }
func (l *logger) Error(v ...any)                     { l.log(LevelError, v...) }
func (l *logger) Errorf(f string, v ...any)          { l.logf(LevelError, f, v...) }
func (l *logger) Panic(v ...any)                     { l.panic(v...) }
func (l *logger) Panicf(f string, v ...any)          { l.panicf(f, v...) }
func (l *logger) Fatal(v ...interface{})             { l.fatal(v...) }
func (l *logger) Fatalf(f string, v ...any)          { l.fatalf(f, v...) }

func (l *logger) log(level Level, v ...any) {
	if l.level < level || len(v) == 0 {
		return
	}

	l.logger.Output(3, fmt.Sprintf(l.format.value, l.levelNames[level], l.labels.formatted, fmt.Sprint(v...)))
}

func (l *logger) logf(level Level, f string, v ...any) {
	if l.level < level {
		return
	}

	l.logger.Output(3, fmt.Sprintf(l.format.value, l.levelNames[level], l.labels.formatted, fmt.Sprintf(f, v...)))
}

func (l *logger) panic(v ...any) {
	if l.level < LevelPanic {
		return
	}

	msg := fmt.Sprintf(l.format.value, l.levelNames[LevelFatal], l.labels.formatted, fmt.Sprint(v...))
	l.logger.Output(3, msg)
	panic(msg)
}

func (l *logger) panicf(f string, v ...any) {
	if l.level < LevelPanic {
		return
	}

	msg := fmt.Sprintf(l.format.value, l.levelNames[LevelFatal], l.labels.formatted, fmt.Sprintf(f, v...))
	l.logger.Output(3, msg)
	panic(msg)
}

func (l *logger) fatal(v ...any) {
	if l.level < LevelFatal {
		return
	}

	l.logger.Output(3, fmt.Sprintf(l.format.value, l.levelNames[LevelFatal], l.labels.formatted, fmt.Sprint(v...)))
	os.Exit(1)
}

func (l *logger) fatalf(f string, v ...any) {
	if l.level < LevelFatal {
		return
	}

	l.logger.Output(3, fmt.Sprintf(l.format.value, l.levelNames[LevelFatal], l.labels.formatted, fmt.Sprintf(f, v...)))
	os.Exit(1)
}
