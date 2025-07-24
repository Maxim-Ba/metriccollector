// Package provides custom static analysis tools for Go code.
//
// This package includes analyzers that extend the standard set of Go static analysis tools.
// It is designed to be used with multichecker to perform comprehensive code checks.
//
// Usage:
//
//	go build -o multichecker
//	multichecker [-flag] [package]
//
// The analyzers in this package focus on detecting specific patterns or potential issues
// in Go code that aren't covered by standard analyzers.
package main

import (
	"golang.org/x/tools/go/analysis"

	"golang.org/x/tools/go/analysis/multichecker"
	"honnef.co/go/tools/staticcheck"

	"github.com/Maxim-Ba/metriccollector/pkg/staticanalizer"
)

func main() {
	var analyzers []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		if len(v.Analyzer.Name) >= 2 && v.Analyzer.Name[:2] == "SA" {
			analyzers = append(analyzers, v.Analyzer)
		}
		if staticanalizer.StaticcheckOthers[v.Analyzer.Name] {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	analyzers = append(analyzers, staticanalizer.StaticchecAnalyzers...)

	multichecker.Main(
		analyzers...,
	)
}
