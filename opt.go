package log

import (
	"io"
	"log"
)

type Opt func(*Opts)

// Flags that set to log.Logger message header.
func Flags(flags int) Opt {
	return func(opts *Opts) {
		opts.Flags = flags
	}
}

// FileAndLine turns on log.Lshotfile on log.Logger,
// which displays the file and a line where the log was called on.
func FileAndLine() Opt {
	return func(opts *Opts) {
		opts.Flags |= log.Lshortfile
	}
}

// UTC sets timzone UTC to the log message header.
func UTC() Opt {
	return func(opts *Opts) {
		opts.Flags |= log.LUTC
	}
}

// Format replaces default format of a logger, there are some predefined placeholders:
//
//	`${level}`: is a logger level name, e.g. DEBUG, INFO etc.;
//
//	`${labels}`: is a placeholder where additional labels supposed to appear, when they are added;
//
//	`${msg}`: is a logger message;
//
//	Default: `[${level}] ${labels} ${msg}\n`
//	Example: `[debug] userId:1000 successfully logged in`
func Format(newFormat string) Opt {
	return func(opts *Opts) {
		opts.Format = newFormat
	}
}

// LevelName changes the level name for a particular log level.
func LevelName(level Level, newName string) Opt {
	return func(opts *Opts) {
		if _, ok := opts.LevelNames[level]; ok {
			opts.LevelNames[level] = newName
		}
	}
}

// LevelNames updates multiple level names of log levels.
func LevelNames(update map[Level]string) Opt {
	return func(opts *Opts) {
		for k, v := range update {
			if _, ok := opts.LevelNames[k]; ok {
				opts.LevelNames[k] = v
			}
		}
	}
}

// UpperCaseNames updates all the level names to ana upper case.
func UpperCaseNames() Opt {
	return func(opts *Opts) {
		opts.UpperCase = true
	}
}

// MinLevel sets the log level.
//
//	Default: info
func MinLevel(level Level) Opt {
	return func(opts *Opts) {
		opts.MinLevel = normalizeLevel(level)
	}
}

// TraceLevel sets a trace log level.
func TraceLevel() Opt {
	return func(opts *Opts) {
		opts.MinLevel = LevelTrace
	}
}

// DebugLevel sets a debug log level.
func DebugLevel() Opt {
	return func(opts *Opts) {
		opts.MinLevel = LevelDebug
	}
}

// InfoLevel sets a info log level.
func InfoLevel() Opt {
	return func(opts *Opts) {
		opts.MinLevel = LevelInfo
	}
}

// WarnLevel sets a warning log level.
func WarnLevel() Opt {
	return func(opts *Opts) {
		opts.MinLevel = LevelWarn
	}
}

// ErrorLevel sets an error log level.
func ErrorLevel() Opt {
	return func(opts *Opts) {
		opts.MinLevel = LevelError
	}
}

// PanicLevel sets a panic log level.
func PanicLevel() Opt {
	return func(opts *Opts) {
		opts.MinLevel = LevelPanic
	}
}

// FatalLevel sets a fatal log level.
func FatalLevel() Opt {
	return func(opts *Opts) {
		opts.MinLevel = LevelError
	}
}

// Writer sets a writer where log will be written to.
func Writer(writer io.Writer) Opt {
	return func(opts *Opts) {
		opts.Writer = writer
	}
}

// CustomerLogger sets a custom log.Logger instance as a base for the logging.
func CustomLogger(logger *log.Logger) Opt {
	return func(opts *Opts) {
		opts.Logger = logger
	}
}

// Labels sets default labels on logs.
func Labels(labels ...string) Opt {
	return func(opts *Opts) {
		opts.Labels = labels
	}
}

// LabelsFormat sets up how the labels should be printed in log inside ${labels} placeholder.
//
//	When there are no labels, ${labels} is replace by empty string.
//	But when there is a label it might be needed to add brackets around, or a space on the left or right.
//	E.g. "(${labels})" or "${labels}:"
//	Default: "${labels}"
func LabelsFormat(newFormat string) Opt {
	return func(opts *Opts) {
		opts.LabelsFormat = newFormat
	}
}

// LabelsSeparator sets a separator between labels in the logs.
//
//	Default: " " (space)
func LabelsSeparator(sep string) Opt {
	return func(opts *Opts) {
		opts.LabelsSeparator = sep
	}
}

type Opts struct {
	Flags           int
	Format          string
	LevelNames      map[Level]string
	Labels          []string
	LabelsFormat    string
	LabelsSeparator string
	MinLevel        Level
	Writer          io.Writer
	Logger          *log.Logger
	UpperCase       bool
}

func defaultOpts() *Opts {
	return &Opts{
		Flags:  log.Ldate | log.Ltime | log.Lmicroseconds,
		Format: defaultFormat,
		LevelNames: map[Level]string{
			LevelFatal: LevelNameFatal,
			LevelPanic: LevelNamePanic,
			LevelError: LevelNameError,
			LevelWarn:  LevelNameWarn,
			LevelInfo:  LevelNameInfo,
			LevelDebug: LevelNameDebug,
			LevelTrace: LevelNameTrace,
		},
		LabelsFormat:    LabelsPlaceholder,
		LabelsSeparator: " ",
		MinLevel:        LevelInfo,
		Logger:          nil,
	}
}

func mergeOpts(base, update *Opts) *Opts {
	if update.Flags != 0 {
		base.Flags = update.Flags
	}
	if update.Format != "" {
		base.Format = update.Format
	}
	if len(update.Labels) > 0 {
		base.Labels = update.Labels
	}
	if update.LabelsFormat != "" {
		base.LabelsFormat = update.LabelsFormat
	}
	if update.LabelsSeparator != "" {
		base.LabelsSeparator = update.LabelsSeparator
	}
	if len(update.LevelNames) > 0 {
		for k, v := range update.LevelNames {
			if _, ok := base.LevelNames[k]; ok {
				base.LevelNames[k] = v
			}
		}
	}
	if update.MinLevel != 0 {
		base.MinLevel = normalizeLevel(update.MinLevel)
	}
	if update.Writer != nil {
		base.Writer = update.Writer
	}
	if update.Logger != nil {
		base.Logger = update.Logger
	}
	base.UpperCase = update.UpperCase
	return base
}
