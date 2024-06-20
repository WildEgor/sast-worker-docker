package sast

import (
	"github.com/WildEgor/sast-worker-docker/internal/adapters/linter"
	"github.com/WildEgor/sast-worker-docker/internal/adapters/vul_checker"
)

func MapCheckResultToAnalysisResult(result []linter.CheckResult) []AnalysisItem {
	ai := make([]AnalysisItem, 0)

	for _, r := range result {
		ai = append(ai, AnalysisItem{
			Line: r.Pos.Line,
			Coll: r.Pos.Coll,
			Code: r.Err.Code,
			Msg:  r.Err.Msg,
		})
	}

	return ai
}

func MapVulItemsToAnalysisResult(result []vul_checker.VulListItem) []AnalysisItem {
	ai := make([]AnalysisItem, 0)

	for _, r := range result {
		ai = append(ai, AnalysisItem{
			Code: r.Code,
			Msg:  r.Msg,
		})
	}

	return ai
}
