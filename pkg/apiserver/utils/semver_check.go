// Copyright 2021 PingCAP, Inc. Licensed under Apache-2.0.

package utils

import (
	"strings"

	"github.com/Masterminds/semver"

	"github.com/pingcap/tidb-dashboard/pkg/utils/version"
)

// IsVersionSupport checks if a semantic version fits within a set of constraints
// pdVersion, standaloneVersion examples: "v5.2.2", "v5.3.0", "v5.4.0-alpha-xxx", "5.3.0" (semver can handle `v` prefix by itself)
// constraints examples: "~5.2.2", ">= 5.3.0", see semver docs to get more information.
func IsVersionSupport(standaloneVersion string, constraints []string) bool {
	curVersion := standaloneVersion
	if version.Standalone == "No" {
		curVersion = version.PDVersion
	}
	// drop "-alpha-xxx" suffix
	versionWithoutSuffix := strings.Split(curVersion, "-")[0]
	v, err := semver.NewVersion(versionWithoutSuffix)
	if err != nil {
		return false
	}
	for _, ver := range constraints {
		c, err := semver.NewConstraint(ver)
		if err != nil {
			continue
		}
		if c.Check(v) {
			return true
		}
	}
	return false
}
