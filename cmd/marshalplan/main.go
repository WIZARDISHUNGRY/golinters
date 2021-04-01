package main

import (
	"github.com/WIZARDISHUNGRY/golinters/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.MarshalPlan)
}
