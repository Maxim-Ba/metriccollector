package staticlint

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestMyAnalyzer(t *testing.T) {
	// функция analysistest.Run применяет тестируемый анализатор OsExitAnalyzer
	// к пакетам os-exit-analizer из папки testdata и проверяет ожидания
	analysistest.Run(t, analysistest.TestData(), OsExitAnalyzer, "./os-exit-analizer/os-exit-analizer.go")
}
