package log

type (
	Logger interface {
		TraceLogger
		DebugLogger
		InfoLogger
		NoticeLogger
		WarnLogger
		ErrorLogger
		PanicLogger
		FatalLogger
		LevelLogger
	}

	TraceLogger interface {
		Trace(...any)
		Tracef(string, ...any)
	}

	DebugLogger interface {
		Debug(...any)
		Debugf(string, ...any)
	}

	InfoLogger interface {
		Info(...any)
		Infof(string, ...any)
	}

	NoticeLogger interface {
		Notice(...any)
		Noticef(string, ...any)
	}

	WarnLogger interface {
		Warn(...any)
		Warnf(string, ...any)
	}

	ErrorLogger interface {
		Error(...any)
		Errorf(string, ...any)
	}

	PanicLogger interface {
		Panic(...any)
		Panicf(string, ...any)
	}

	FatalLogger interface {
		Fatal(...any)
		Fatalf(string, ...any)
	}

	LevelLogger interface {
		Log(Level, ...any)
		Logf(Level, string, ...any)
	}
)
