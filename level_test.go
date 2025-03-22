package log

import "testing"

func Test_normalizeLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		level    Level
		expected Level
	}{
		{
			name:     "smaller-than-fatal",
			level:    0,
			expected: LevelFatal,
		},
		{
			name:     "bigger-than-trace",
			level:    100,
			expected: LevelTrace,
		},
		{
			name:     "trace",
			level:    LevelTrace,
			expected: LevelTrace,
		},
		{
			name:     "debug",
			level:    LevelDebug,
			expected: LevelDebug,
		},
		{
			name:     "info",
			level:    LevelInfo,
			expected: LevelInfo,
		},
		{
			name:     "warn",
			level:    LevelWarn,
			expected: LevelWarn,
		},
		{
			name:     "error",
			level:    LevelError,
			expected: LevelError,
		},
		{
			name:     "panic",
			level:    LevelPanic,
			expected: LevelPanic,
		},
		{
			name:     "fatal",
			level:    LevelFatal,
			expected: LevelFatal,
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := normalizeLevel(test.level)
			if result != test.expected {
				t.Errorf("expected %q, but got %q", test.expected, result)
			}
		})
	}
}
