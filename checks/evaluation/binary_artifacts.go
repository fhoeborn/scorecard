// Copyright 2021 OpenSSF Scorecard Authors
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
	"github.com/ossf/scorecard/v4/probes/freeOfUnverifiedBinaryArtifacts"
)

// BinaryArtifacts applies the score policy for the Binary-Artifacts check.
func BinaryArtifacts(name string,
	findings []finding.Finding,
	dl checker.DetailLogger,
) checker.CheckResult {
	expectedProbes := []string{
		freeOfUnverifiedBinaryArtifacts.Probe,
	}

	if !finding.UniqueProbesEqual(findings, expectedProbes) {
		e := sce.WithMessage(sce.ErrScorecardInternal, "invalid probe results")
		return checker.CreateRuntimeErrorResult(name, e)
	}

	if findings[0].Outcome == finding.OutcomeTrue {
		return checker.CreateMaxScoreResult(name, "no binaries found in the repo")
	}

	for i := range findings {
		f := &findings[i]
		if f.Outcome != finding.OutcomeFalse {
			continue
		}
		dl.Warn(&checker.LogMessage{
			Path:   f.Location.Path,
			Type:   f.Location.Type,
			Offset: *f.Location.LineStart,
			Text:   "binary detected",
		})
	}

	// There are only false findings.
	// Deduct the number of findings from max score
	numberOfBinaryFilesFound := len(findings)

	score := checker.MaxResultScore - numberOfBinaryFilesFound

	if score < checker.MinResultScore {
		score = checker.MinResultScore
	}

	return checker.CreateResultWithScore(name, "binaries present in source code", score)
}
