package staticanalizer

import (
	"golang.org/x/tools/go/analysis"

	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"

	"github.com/Maxim-Ba/metriccollector/cmd/staticlint"
	"github.com/kisielk/errcheck/errcheck"
	"github.com/timakin/bodyclose/passes/bodyclose"
)

var StaticcheckOthers = map[string]bool{
	"ST1000": true,
	"S1000":  true,
	"QF1001": true,
}

var StaticchecAnalyzers = []*analysis.Analyzer{
	shadow.Analyzer,
	printf.Analyzer,
	structtag.Analyzer,
	copylock.Analyzer,
	loopclosure.Analyzer,
	lostcancel.Analyzer,
	unreachable.Analyzer,
	unmarshal.Analyzer,
	// ----- others -----
	errcheck.Analyzer,
	bodyclose.Analyzer,
	staticlint.OsExitAnalyzer,
}
