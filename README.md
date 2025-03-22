# log.go

Simple golang logger with levels based on standard library log.Logger.

## Usage
```go
import (
    "context"

    "github.com/krynka/log.go"
)

func main() {
    logger := log.New(log.TraceLevel())
    logger.Debug("log message")

    logger = log.WithLabels(logger, "user=1000")
    logger.Debug("log message")

    ctx = log.ToContext(context.Background(), logger)
    log.FromContext(ctx).Debug("log message from context")
}
```

Prints:
```
2025/03/22 15:07:50.348957 [debug] log message
2025/03/22 15:07:50.348957 [debug] user=1000 log message
2025/03/22 15:07:50.348957 [debug] user=1000 log message from context
```

## Interface

```go
type Logger interface {
    Trace(v ...any)
    Tracef(format string, v ...any)
    Debug(v ...any)
    Debugf(format string, v ...any)
    Info(v ...any)
    Infof(format string, v ...any)
    Notice(v ...any)
    Noticef(format string, v ...any)
    Warn(v ...any)
    Warnf(format string, v ...any)
    Error(v ...any)
    Errorf(format string, v ...any)
    Panic(v ...any)
    Panicf(format string, v ...any)
    Fatal(v ...any)
    Fatalf(format string, v ...any)
    Log(level Level, v ...any)
    Logf(level Level, format string, v ...any)
}
```

The same log level Info and Notice were added just in case some library expects Notice or Info in the logger interface.

## Logger configuration

Logger can be modified:
```go
log.New(
    // allows to name each log level
    log.LevelNames(map[Level]string{
        LevelTrace: "Trace",
    }),
    // updates all level names to upper case
    log.UpperCaseNames(),
    // allows to set the log writer, e.g. file, network etc.
    log.Writer(fileWriter),
    // sets flags to standard log.Logger
    log.Flags(stdlog.Ldate | stdlog.Ltime | stdlog.Lmicroseconds),
    // sets debug log level etc.
    log.DebugLevel(),
    // sets the base log.Logger
    log.CustomLogger(stdlog.Default()),
    // etc
)
```

## Logger message

```go
logger := log.New()
logger.Info("log message")
```

By default the log looks like
```
2025/03/22 15:07:50.348957 [info] log message
```

Message could be configured via options in constructor:
```go
logger := log.New(log.FileAndLine(), log.UpperCaseNames(), log.Format("${level}: ${msg}"))
logger.Info("log message")
```
Prints:
```
2025/03/22 15:07:50.348957 main.go:10: INFO: log message
```

## Labels

Message can include additional labels, that prints in each log.

It can be useful to include user id or any other details to the logger, that can help to understand where does the log come from.

```go
logger := log.New()
logger.Info("message without labels")
labeledLogger := log.WithLabels(logger, "userId=1000", "traceId=108")
labeledLogger.Info("message with labels")
```

```
2025/03/22 15:07:50.348957 [info] message without labels
2025/03/22 15:07:50.348957 [info] userId=1000 traceId=108 message with labels
```

## Message format

By default message format looks like the following:
```
[${level}] ${labels} ${msg}
```

Which prints the following without labels:
```
[info] log message
```

And prints the following with labels:
```
[info] user=1000 log message
```

This format can be overwritten either during logger creation or after it:
```go
logger := log.New(log.Format("${labels} ${level}: ${msg}"))
newLogger := log.WithFormat(logger, "${level} - ${msg}")
```

## Context

Logger can be used with context:
```go
ctx = log.ToContext(ctx, log.New())
log.FromContext(ctx).Info("log message")
```
