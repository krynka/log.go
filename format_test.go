package log

import "testing"

func Test_buildFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		format    string
		hasLabels bool
		expected  string
	}{
		{
			name:      "empty-no-labels",
			format:    "",
			hasLabels: false,
			expected:  "%[3]s\n",
		},
		{
			name:      "empty-with-labels",
			format:    "",
			hasLabels: false,
			expected:  "%[3]s\n",
		},
		{
			name:      "all-set-no-labels",
			format:    "${level}: ${labels} ${msg}",
			hasLabels: false,
			expected:  "%[1]s: %[3]s\n",
		},
		{
			name:      "all-set-with-labels",
			format:    "${level}: ${labels} ${msg}",
			hasLabels: true,
			expected:  "%[1]s: %[2]s %[3]s\n",
		},
		{
			name:      "with-new-line-character",
			format:    "${level}: ${labels} ${msg}\n",
			hasLabels: true,
			expected:  "%[1]s: %[2]s %[3]s\n",
		},
		{
			name:      "no-placeholders-with-no-space-in-the-end",
			format:    "LOG:",
			hasLabels: true,
			expected:  "LOG: %[3]s\n",
		},
		{
			name:      "no-placeholders-with-space-in-the-end",
			format:    "LOG: ",
			hasLabels: true,
			expected:  "LOG: %[3]s\n",
		},
		{
			name:      "no-placeholders-with-new-line-in-the-end",
			format:    "LOG:\n",
			hasLabels: true,
			expected:  "LOG:\n%[3]s\n",
		},
		{
			name:      "multiple-placeholders",
			format:    "${level}: ${labels} [${level}] ${msg}",
			hasLabels: true,
			expected:  "%[1]s: %[2]s [%[1]s] %[3]s\n",
		},
		{
			name:      "no-msg-placeholders",
			format:    "${level}: ${labels}",
			hasLabels: true,
			expected:  "%[1]s: %[2]s %[3]s\n",
		},
		{
			name:      "no-labels",
			format:    "${labels} ${level}: ${labels} ${msg} ${labels}",
			hasLabels: false,
			expected:  "%[1]s: %[3]s\n",
		},
		{
			name:      "additional placeholders",
			format:    "[${level}] %f ${labels} %d ${msg} %s",
			hasLabels: true,
			expected:  "[%[1]s] %%f %[2]s %%d %[3]s %%s\n",
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			format := buildFormat(test.format, test.hasLabels)
			if format.value != test.expected {
				t.Errorf("expected %q, got %q", test.expected, format.value)
			}
		})
	}
}

func Test_format(t *testing.T) {
	t.Parallel()

	defaultFormat := buildFormat(defaultFormat, true)
	if defaultFormat.value != "[%[1]s] %[2]s %[3]s\n" {
		t.Fatalf("expected %q, got %q", "[%[1]s] %[2]s %[3]s\n", defaultFormat.value)
	}
	noLabelsFormat := defaultFormat.clearLabels()
	if noLabelsFormat.value != "[%[1]s] %[3]s\n" {
		t.Fatalf("expected %q, got %q", "[%[1]s] %[3]s\n", noLabelsFormat.value)
	}
	if noLabelsFormat.value == defaultFormat.value {
		t.Fatalf("format was built from a copy after clearLabels(), but original format was edited")
	}
	if defaultFormat.value != "[%[1]s] %[2]s %[3]s\n" {
		t.Fatalf("original format was updated after clearLabels()")
	}
	labeledFormat := noLabelsFormat.withLabels(true)
	if labeledFormat.value != "[%[1]s] %[2]s %[3]s\n" {
		t.Fatalf("expected %q, got %q", "[%[1]s] %[3]s\n", labeledFormat.value)
	}
	if noLabelsFormat.value != "[%[1]s] %[3]s\n" {
		t.Fatalf("original format was updated after withLabels(true)")
	}
}
