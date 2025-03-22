package main

import (
	"context"
	"errors"

	"github.com/krynka/log.go"
)

// Prints:
// 2025/03/22 17:58:49.512039 [info] default level is info
// 2025/03/22 17:58:49.512180 [debug] log from context: log message
// 2025/03/22 17:58:49.512184 [error] error occured: error log
// 2025/03/22 17:58:49.512186 [info] user:1000 message with label
// 2025/03/22 17:58:49.512187 [info] message without label
// 2025/03/22 17:58:49.512192 example.go:41: [TRC] worker-1 user=2000 custom log
// 2025/03/22 17:58:49.512215 example.go:42: [FATAL] worker-1 user=2000 fatal log
// exit status 1
func main() {
	origLog := log.New()

	origLog.Debug("this won't be logged")
	origLog.Info("default level is info")

	newLog := log.WithLevel(origLog, log.LevelTrace)

	ctx := log.ToContext(context.Background(), newLog)
	if err := logFromContext(ctx, "log message"); err != nil {
		newLog.Errorf("error occured: %s", err)
	}

	newLog = log.WithLabels(newLog, "user:1000")
	newLog.Info("message with label")
	origLog.Info("message without label")

	customLog := log.ByLevelName("TRC",
		log.LevelNames(map[log.Level]string{
			log.LevelTrace: "trc",
			log.LevelDebug: "dbg",
			log.LevelInfo:  "inf",
			log.LevelWarn:  "wrn",
			log.LevelError: "err",
		}),
		log.UpperCaseNames(),
		log.Format("[example] ${labels} ${level}: ${msg}"),
		log.Labels("worker-1", "user=2000"),
		log.FileAndLine(),
	)

	customLog.Trace("custom log")
	customLog.Fatal("fatal log")
}

func logFromContext(ctx context.Context, msg string) error {
	log.FromContext(ctx).Debugf("log from context: %s", msg)
	return errors.New("error log")
}
