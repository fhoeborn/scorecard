// Copyright 2024 OpenSSF Scorecard Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package evaluation

import (
	"github.com/ossf/scorecard/v4/checker"
	sce "github.com/ossf/scorecard/v4/errors"
	"github.com/ossf/scorecard/v4/finding"
	"github.com/ossf/scorecard/v4/probes/hasPermissiveLicense"
)

// License applies the score policy for the License check.
func PermissiveLicense(name string,
	findings []finding.Finding,
	dl checker.DetailLogger,
) checker.CheckResult {
	// We have 3 unique probes, each should have a finding.
	expectedProbes := []string{
		hasPermissiveLicense.Probe,
	}

	if !finding.UniqueProbesEqual(findings, expectedProbes) {
		e := sce.WithMessage(sce.ErrScorecardInternal, "invalid probe results")
		return checker.CreateRuntimeErrorResult(name, e)
	}

	// Compute the score.
	score := 0
	m := make(map[string]bool)
	for i := range findings {
		f := &findings[i]
		switch f.Outcome {
		case finding.OutcomeNotApplicable:
			dl.Info(&checker.LogMessage{
				Type:   finding.FileTypeSource,
				Offset: 1,
				Text:   f.Message,
			})
		case finding.OutcomePositive:
			switch f.Probe {
			case hasPermissiveLicense.Probe:
				score += scoreProbeOnce(f.Probe, m, 10)
			default:
				e := sce.WithMessage(sce.ErrScorecardInternal, "unknown probe results ")
				return checker.CreateRuntimeErrorResult(name, e)
			}
		case finding.OutcomeNegative:
			switch f.Probe {
			case hasPermissiveLicense.Probe:
				dl.Info(&checker.LogMessage{
					Text: "Non-permissive license found " + f.Message,
				})
			}
		default:
			continue // for linting
		}
	}
	_, defined := m[hasPermissiveLicense.Probe]
	if !defined {
		if score > 0 {
			e := sce.WithMessage(sce.ErrScorecardInternal, "score calculation problem")
			return checker.CreateRuntimeErrorResult(name, e)
		}
		return checker.CreateMinScoreResult(name, "Non-Permissive license detected")
	}
	return checker.CreateMaxScoreResult(name, "Permissive license detected")
}
