package hush

import (
	"bytes"
	"testing"

	"github.com/fatih/color"
)

func TestRun(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		input      string
		expectCode int
		expectOut  string
	}{
		{
			input:     "ls",
			expectOut: color.GreenString("➜") + " hush $ ls",
		},
	} {
		tc := tc // Enable parallel sub-tests
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			var in, out bytes.Buffer
			exitCode := run(&in, &out, &out, []string{"hush", "-c", tc.input})
			output := out.String()
			if tc.expectCode != exitCode {
				t.Errorf("Unexpected exit code.\nExpected: %d\nActual:  %d", tc.expectCode, exitCode)
			}
			if tc.expectOut != output {
				t.Errorf("Unexpected output.\nExpected: %s\nActual:   %s", tc.expectOut, output)
			}
		})
	}
}
