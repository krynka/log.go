package log

import (
	"log"
	"sort"
	"testing"
)

func Test_mergeOpts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		base     Opts
		update   Opts
		expected Opts
	}{
		{
			name: "flags",
			base: Opts{
				Flags: log.LstdFlags,
			},
			update: Opts{
				Flags: log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile,
			},
			expected: Opts{
				Flags: log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile,
			},
		},
		{
			name: "update-level-names",
			base: Opts{
				LevelNames: map[Level]string{
					LevelFatal: LevelNameFatal,
					LevelPanic: LevelNamePanic,
					LevelError: LevelNameError,
					LevelWarn:  LevelNameWarn,
					LevelInfo:  LevelNameInfo,
					LevelDebug: LevelNameDebug,
					LevelTrace: LevelNameTrace,
				},
			},
			update: Opts{
				LevelNames: map[Level]string{
					LevelWarn: "warning",
					LevelInfo: "information",
				},
			},
			expected: Opts{
				LevelNames: map[Level]string{
					LevelFatal: LevelNameFatal,
					LevelPanic: LevelNamePanic,
					LevelError: LevelNameError,
					LevelWarn:  "warning",
					LevelInfo:  "information",
					LevelDebug: LevelNameDebug,
					LevelTrace: LevelNameTrace,
				},
			},
		},
		{
			name: "upper-case",
			base: Opts{
				UpperCase: false,
			},
			update: Opts{
				UpperCase: true,
			},
			expected: Opts{
				UpperCase: true,
			},
		},
		{
			name: "min-level",
			base: Opts{
				MinLevel: LevelError,
			},
			update: Opts{
				MinLevel: LevelDebug,
			},
			expected: Opts{
				MinLevel: LevelDebug,
			},
		},
		{
			name: "min-level-no-update",
			base: Opts{
				MinLevel: LevelError,
			},
			update: Opts{},
			expected: Opts{
				MinLevel: LevelError,
			},
		},
		{
			name: "min-level-too-small",
			base: Opts{
				MinLevel: LevelError,
			},
			update: Opts{
				MinLevel: -10,
			},
			expected: Opts{
				MinLevel: LevelFatal,
			},
		},
		{
			name: "format",
			base: Opts{
				Format: "${level}: ${msg}",
			},
			update: Opts{
				Format: "[${level}]: ${labels} ${msg}",
			},
			expected: Opts{
				Format: "[${level}]: ${labels} ${msg}",
			},
		},
		{
			name: "labels",
			base: Opts{
				Labels: []string{"100", "15"},
			},
			update: Opts{
				Labels: []string{"101", "100"},
			},
			expected: Opts{
				Labels: []string{"101", "100"},
			},
		},
		{
			name: "writer",
			base: Opts{
				Logger: nil,
			},
			update: Opts{
				Writer: log.Default().Writer(),
			},
			expected: Opts{
				Writer: log.Default().Writer(),
			},
		},
		{
			name: "logger",
			base: Opts{
				Logger: nil,
			},
			update: Opts{
				Logger: log.Default(),
			},
			expected: Opts{
				Logger: log.Default(),
			},
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := mergeOpts(&test.base, &test.update)

			if test.expected.Flags != result.Flags {
				t.Errorf("flags expected %d, but got %d", test.expected.Flags, result.Flags)
			}

			if test.expected.Format != result.Format {
				t.Errorf("format expected %d, but got %d", test.expected.Flags, result.Flags)
			}

			for level, name := range test.expected.LevelNames {
				if result.LevelNames[level] != name {
					t.Errorf("levelName of %d level expected %q, but got %q", level, name, result.LevelNames[level])
				}
			}

			sort.Slice(result.Labels, func(i, j int) bool { return result.Labels[i] < result.Labels[j] })
			sort.Slice(test.expected.Labels, func(i, j int) bool { return test.expected.Labels[i] < test.expected.Labels[j] })
			if len(test.expected.Labels) == len(result.Labels) {
				for i := range result.Labels {
					if test.expected.Labels[i] != result.Labels[i] {
						t.Errorf("labels expected %#v, but got %#v", test.expected.Labels, result.Labels)
					}
				}
			} else {
				t.Errorf("labels expected %#v, but got %#v", test.expected.Labels, result.Labels)
			}

			if test.expected.LabelsFormat != result.LabelsFormat {
				t.Errorf("labels format expected %q, but got %q", test.expected.LabelsFormat, result.LabelsFormat)
			}

			if test.expected.LabelsSeparator != result.LabelsSeparator {
				t.Errorf("labels separator expected %q, but got %q", test.expected.LabelsSeparator, result.LabelsSeparator)
			}

			if test.expected.MinLevel != result.MinLevel {
				t.Errorf("min level expected %d, but got %d", test.expected.MinLevel, result.MinLevel)
			}

			if test.expected.Writer != result.Writer {
				t.Errorf("writer expected %#v, but got %#v", test.expected.Writer, result.Writer)
			}

			if test.expected.Logger != result.Logger {
				t.Errorf("logger expected %#v, but got %#v", test.expected.Logger, result.Logger)
			}

			if test.expected.UpperCase != result.UpperCase {
				t.Errorf("logger expected %t, but got %t", test.expected.UpperCase, result.UpperCase)
			}
		})
	}
}
