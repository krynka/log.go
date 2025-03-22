package log

import "strings"

type Level int

// Fatal is the smallest log level, the narrowest visibility of the logs.
// Trace is the highest log level, the widest visibility of the logs.
const (
	LevelFatal Level = iota + 1
	LevelPanic
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

const (
	LevelNameFatal = "fatal"
	LevelNamePanic = "panic"
	LevelNameError = "error"
	LevelNameWarn  = "warn"
	LevelNameInfo  = "info"
	LevelNameDebug = "debug"
	LevelNameTrace = "trace"
)

func defaultLevelName(level Level) string {
	switch level {
	case LevelTrace:
		return LevelNameTrace
	case LevelDebug:
		return LevelNameDebug
	case LevelInfo:
		return LevelNameInfo
	case LevelWarn:
		return LevelNameWarn
	case LevelError:
		return LevelNameError
	case LevelPanic:
		return LevelNamePanic
	case LevelFatal:
		return LevelNameFatal
	}
	return ""
}

func normalizeLevel(level Level) Level {
	if level < LevelFatal {
		return LevelFatal
	}
	if level > LevelTrace {
		return LevelTrace
	}
	return level
}

func buildLevelNames(opts Opts) map[Level]string {
	levelNames := copyLevelNames(opts.LevelNames)
	if opts.UpperCase {
		for k, v := range levelNames {
			levelNames[k] = strings.ToUpper(v)
		}
	}
	return levelNames
}

func copyLevelNames(m map[Level]string) map[Level]string {
	levelNames := make(map[Level]string, LevelTrace)
	for level := LevelFatal; level <= LevelTrace; level++ {
		levelName, ok := m[level]
		if !ok {
			levelName = defaultLevelName(level)
		}
		levelNames[level] = levelName
	}
	return levelNames
}

// WithLevel returns a new logger from existing with a new log level,
// the original Logger keeps an original log level.
func WithLevel(l Logger, newLevel Level) Logger {
	log, ok := l.(*logger)
	if !ok {
		return l
	}

	newLevel = normalizeLevel(newLevel)
	if log.level == newLevel {
		return log
	}

	return &logger{
		level:      newLevel,
		format:     log.format,
		levelNames: log.levelNames,
		labels:     log.labels,
		logger:     log.logger,
	}
}
