package main

import (
	"golang.org/x/tools/go/analysis"

	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"honnef.co/go/tools/staticcheck"

	"github.com/Maxim-Ba/metriccollector/cmd/staticlint"
	"github.com/kisielk/errcheck/errcheck"
	"github.com/timakin/bodyclose/passes/bodyclose"
)

func main() {
	staticcheckOthers := map[string]bool{
		"ST1000": true,
		"S1000":  true,
		"QF1001": true,
	}
	var analyzers []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		if len(v.Analyzer.Name) >= 2 && v.Analyzer.Name[:2] == "SA" {
			analyzers = append(analyzers, v.Analyzer)
		}
		if staticcheckOthers[v.Analyzer.Name] {
			analyzers = append(analyzers, v.Analyzer)
		}
	}
	analyzers = append(analyzers,
		shadow.Analyzer,
		printf.Analyzer,
		structtag.Analyzer,
		copylock.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		unreachable.Analyzer,
		unmarshal.Analyzer,
	)
	analyzers = append(analyzers,
		errcheck.Analyzer, 
		bodyclose.Analyzer, 
		staticlint.OsExitAnalyzer,
	)

	multichecker.Main(
		analyzers...,
	)
}
