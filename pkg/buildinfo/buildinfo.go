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
func PrintBuildInfo(buildVersion,buildDate,  buildCommit string) {
	if buildVersion == "" {
		buildVersion = "undefined"
	}
	if buildDate == "" {
		buildDate = "undefined"
	}
	if buildCommit == "" {
		buildCommit = "undefined"
	}

	// Форматированный вывод информации
	fmt.Fprintf(os.Stdout, "Build version: %s\n", buildVersion)
	fmt.Fprintf(os.Stdout, "Build date: %s\n", buildDate)
	fmt.Fprintf(os.Stdout, "Build commit: %s\n", buildCommit)
}
