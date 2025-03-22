package log

import (
	"strings"
)

const (
	defaultFormat = "[${level}] ${labels} ${msg}\n"

	LevelPlaceholder   = "${level}"
	levelFormat        = "%[1]s"
	LabelsPlaceholder  = "${labels}"
	labelsFormat       = "%[2]s"
	MessagePlaceholder = "${msg}"
	messageFormat      = "%[3]s"

	newLine = "\n"
)

type format struct {
	original  string
	value     string
	hasLabels bool
}

func (f format) withLabels(hasLabels bool) format {
	if f.hasLabels == hasLabels {
		return f
	}
	return buildFormat(f.original, hasLabels)
}

func (f format) clearLabels() format {
	return f.withLabels(false)
}

// WithFormat returns a new logger with a new format, the original Logger keeps an original format.
//
// There are some predefined placeholders that could be used:
//
//	`${level}`: is a logger level name, e.g. DEBUG, INFO etc.;
//
//	`${labels}`: is a placeholder where additional labels supposed to appear, when they are added;
//
//	`${msg}`: is a logger message;
//
//	New format: `<worker-1> [${level}] ${labels} ${msg}`
//	Example: `<worker-1> [debug] userId:1000 successfully updated`
func WithFormat(l Logger, newFormat string) Logger {
	log, ok := l.(*logger)
	if !ok {
		return log
	}
	return &logger{
		level:      log.level,
		format:     buildFormat(newFormat, log.labels.notEmpty()),
		levelNames: log.levelNames,
		labels:     log.labels,
		logger:     log.logger,
	}
}

func buildFormat(newFormat string, hasLabels bool) format {
	original := newFormat
	newFormat = escapeFormats(newFormat)

	newFormat = strings.ReplaceAll(newFormat, LevelPlaceholder, levelFormat)

	if hasLabels {
		newFormat = strings.ReplaceAll(newFormat, LabelsPlaceholder, labelsFormat)
	} else {
		newFormat = strings.ReplaceAll(newFormat, " ${labels} ", " ")
		newFormat = strings.ReplaceAll(newFormat, LabelsPlaceholder, "")
	}

	if strings.Contains(newFormat, MessagePlaceholder) {
		newFormat = strings.ReplaceAll(newFormat, MessagePlaceholder, messageFormat)
	} else {
		if len(newFormat) != 0 {
			lastLetter := newFormat[len(newFormat)-1]
			if lastLetter != ' ' && lastLetter != '\n' {
				newFormat += " "
			}
		}
		newFormat += messageFormat + newLine
	}

	newFormat = strings.TrimSpace(newFormat)

	if newFormat[len(newFormat)-1] != '\n' {
		newFormat += newLine
	}
	return format{
		original:  original,
		value:     newFormat,
		hasLabels: hasLabels,
	}
}

func escapeFormats(format string) string {
	return strings.ReplaceAll(format, "%", "%%")
}
