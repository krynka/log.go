package log

import "testing"

func Test_parseLabelsFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{
			name:     "no-labels-placeholder",
			format:   "LABELS",
			expected: "LABELS",
		},
		{
			name:     "simple-placeholder",
			format:   LabelsPlaceholder,
			expected: "%[1]s",
		},
		{
			name:     "in-brackets",
			format:   "(${labels})",
			expected: "(%[1]s)",
		},
		{
			name:     "wth-placeholders",
			format:   "(%[1]s ${labels} %s %d)",
			expected: "(%%[1]s %[1]s %%s %%d)",
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := parseLabelsFormat(test.format)
			if result != test.expected {
				t.Errorf("expected %q, but got %q", test.expected, result)
			}
		})
	}
}

func Test_buildLabels(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		format    string
		values    []string
		separator string
		expected  string
	}{
		{
			name:      "empty-labels-with-space-separator",
			format:    "${labels}",
			values:    nil,
			separator: " ",
			expected:  "",
		},
		{
			name:      "empty-labels-with-comma-separator",
			format:    "${labels}",
			values:    []string{},
			separator: ", ",
			expected:  "",
		},
		{
			name:      "empty-labels-within-brackets",
			format:    "(${labels})",
			values:    nil,
			separator: " ",
			expected:  "",
		},
		{
			name:      "labels-with-space-separator",
			format:    "${labels}",
			values:    []string{"A", "B", "C"},
			separator: " ",
			expected:  "A B C",
		},
		{
			name:      "labels-with-comma-separator",
			format:    "${labels}",
			values:    []string{"A", "B", "C"},
			separator: ", ",
			expected:  "A, B, C",
		},
		{
			name:      "labels-within-brackets",
			format:    "(${labels})",
			values:    []string{"A", "B", "C"},
			separator: "-",
			expected:  "(A-B-C)",
		},
		{
			name:      "with-other-formatters",
			format:    "%s ${labels} %d %[1]s",
			values:    []string{"A", "B", "C"},
			separator: "-",
			expected:  "%s A-B-C %d %[1]s",
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			labels := buildLabels(parseLabelsFormat(test.format), test.values, test.separator)
			if labels.formatted != test.expected {
				t.Errorf("expected %q, but got %q", test.expected, labels.formatted)
			}
		})
	}
}

func Test_labels(t *testing.T) {
	t.Parallel()

	labels := buildLabels("%[1]s", nil, " ")
	newLabels := labels.add("a", "1")
	if labels.values != nil {
		t.Fatalf("original labels values were updated during add(), but expected to stay unmodified")
	}
	if labels.formatted != "" {
		t.Fatalf("original labels formatted was updated during add(), but expected to stay unmodified")
	}
	if newLabels.formatted != "a 1" {
		t.Fatalf("new labels formatted expected %q, but got %q", "a 1", newLabels.formatted)
	}

	formattedLabels := newLabels.setFormat("(%[1]s)")
	if newLabels.formatted != "a 1" {
		t.Fatalf("original labels formatted was updated during setFormat(), but expected to stay unmodified")
	}
	if formattedLabels.formatted != "(a 1)" {
		t.Fatalf("new labels formatted expected %q, but got %q", "(a 1)", newLabels.formatted)
	}
	clearLabels := formattedLabels.clear()
	if formattedLabels.formatted != "(a 1)" {
		t.Fatalf("original labels formatted was updated during clear(), but expected to stay unmodified")
	}
	if clearLabels.formatted != "" {
		t.Fatalf("new labels formatted expected %q, but got %q", "", newLabels.formatted)
	}
}
