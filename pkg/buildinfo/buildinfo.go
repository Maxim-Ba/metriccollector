package buildinfo

import (
	"fmt"
	"os"
)

// printBuildInfo prints the build information to stdout.
// It displays the build version, date, and commit hash.
// If any of these values are not set (empty string),
// it will display "N/A" instead.
// The output format is:
// Build version: <value>
// Build date: <value>
// Build commit: <value>
func PrintBuildInfo(buildVersion, buildDate, buildCommit string) {
	getValue := func(s string) string {
		if s == "" {
			return "N/A"
		}
		return s
	}
	// Форматированный вывод информации
	fmt.Fprintf(os.Stdout, "Build version: %s\n", getValue(buildVersion))
	fmt.Fprintf(os.Stdout, "Build date: %s\n", getValue(buildDate))
	fmt.Fprintf(os.Stdout, "Build commit: %s\n", getValue(buildCommit))
}
