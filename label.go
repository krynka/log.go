package log

import (
	"fmt"
	"strings"
)

type labels struct {
	values    []string
	separator string
	format    string
	formatted string
}

func (l labels) notEmpty() bool {
	return len(l.formatted) > 0
}

func (l labels) isEmpty() bool {
	return len(l.formatted) == 0
}

func (l labels) clear() labels {
	if l.isEmpty() {
		return l
	}
	return labels{
		values:    nil,
		separator: l.separator,
		format:    l.format,
		formatted: "",
	}
}

func (l labels) add(newValues ...string) labels {
	values := joinLabels(l.values, newValues)
	if len(values) == len(l.values) {
		// probably empty values were added
		return l
	}
	return buildLabels(l.format, values, l.separator)
}

func (l labels) setSeparator(sep string) labels {
	if l.separator == sep {
		return l
	}
	return buildLabels(l.format, l.values, sep)
}

func (l labels) setFormat(newFormat string) labels {
	if l.format == newFormat {
		return l
	}
	return buildLabels(newFormat, l.values, l.separator)
}

func buildLabels(format string, values []string, sep string) labels {
	formatted := ""
	if len(values) > 0 {
		formatted = format
		if strings.Contains(format, "%[1]s") {
			formatted = fmt.Sprintf(format, strings.Join(values, sep))
		}
	}

	return labels{
		values:    values,
		separator: sep,
		format:    format,
		formatted: formatted,
	}
}

func copyLabels(labels []string) []string {
	if len(labels) == 0 {
		return nil
	}
	newLabels := make([]string, 0, len(labels))
	for _, label := range labels {
		if label != "" {
			newLabels = append(newLabels, escapeFormats(label))
		}
	}
	return newLabels
}

func parseLabelsFormat(format string) string {
	format = escapeFormats(format)
	return strings.ReplaceAll(format, LabelsPlaceholder, "%[1]s")
}

// ClearLabels remove all the labels and returns another instance of Logger, so the original logger is not affected and keeps all the labels.
func ClearLabels(l Logger) Logger {
	log, ok := l.(*logger)
	if !ok {
		return l
	}
	return &logger{
		level:      log.level,
		format:     log.format.clearLabels(),
		levelNames: log.levelNames,
		labels:     log.labels.clear(),
		logger:     log.logger,
	}
}

// WithLabels adds label(s) returns another instance of Logger,
// so the original logger is not affected and keeps previous set of labels.
func WithLabels(l Logger, labels ...string) Logger {
	log, ok := l.(*logger)
	if !ok {
		return l
	}
	newLabels := log.labels.add(labels...)
	return &logger{
		level:      log.level,
		format:     log.format.withLabels(newLabels.notEmpty()),
		levelNames: log.levelNames,
		labels:     newLabels,
		logger:     log.logger,
	}
}

// WithLabelSeparator replaces labels separator and returns another instance of Logger,
// so the original logger is not affected and keeps previous separator.
func WithLabelSeparator(l Logger, sep string) Logger {
	log, ok := l.(*logger)
	if !ok {
		return l
	}
	return &logger{
		level:      log.level,
		format:     log.format,
		levelNames: log.levelNames,
		labels:     log.labels.setSeparator(sep),
		logger:     log.logger,
	}
}

// WithLabelsFormat replaces labels format and returns another instance of Logger,
// so the original logger is not affected and keeps previous format.
func WithLabelsFormat(l Logger, newFormat string) Logger {
	log, ok := l.(*logger)
	if !ok {
		return l
	}
	return &logger{
		level:      log.level,
		format:     log.format,
		levelNames: log.levelNames,
		labels:     log.labels.setFormat(parseLabelsFormat(newFormat)),
		logger:     log.logger,
	}
}

func joinLabels(labels []string, newLabels []string) []string {
	result := make([]string, 0, len(labels)+len(newLabels))
	result = append(result, labels...)
	for _, label := range newLabels {
		if label != "" {
			result = append(result, label)
		}
	}
	return result
}
