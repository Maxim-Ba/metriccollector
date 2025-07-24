package buildinfo

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrintBuildInfo(t *testing.T) {
	tests := []struct {
		name         string
		version      string
		date         string
		commit       string
		expectOutput string
	}{
		{
			name:         "all fields set",
			version:      "v1.0.0",
			date:         "2023-01-01",
			commit:       "abc123",
			expectOutput: "Build version: v1.0.0\nBuild date: 2023-01-01\nBuild commit: abc123\n",
		},
		{
			name:         "all fields empty",
			version:      "",
			date:         "",
			commit:       "",
			expectOutput: "Build version: N/A\nBuild date: N/A\nBuild commit: N/A\n",
		},
		{
			name:         "some fields empty",
			version:      "v1.0.0",
			date:         "",
			commit:       "abc123",
			expectOutput: "Build version: v1.0.0\nBuild date: N/A\nBuild commit: abc123\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildVersion := tt.version
			buildDate := tt.date
			buildCommit := tt.commit

			// Redirect stdout to capture output
			old := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				require.NoError(t, err)
			}
			os.Stdout = w

			PrintBuildInfo(buildVersion, buildDate, buildCommit)

			err = w.Close()
			if err != nil {
				require.NoError(t, err)
			}
			os.Stdout = old

			var buf bytes.Buffer
			_, err = buf.ReadFrom(r)
			if err != nil {
				require.NoError(t, err)
			}
			output := buf.String()

			if output != tt.expectOutput {
				t.Errorf("Expected output:\n%s\nGot:\n%s", tt.expectOutput, output)
			}
		})
	}
}
