package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackagesInfo_Get(t *testing.T) {
	t.Parallel()

	pkgsInfo := PackagesInfo{
		Packages: map[string]PackageInfo{
			"github.com/smartcontractkit/package": {
				ImportPath: "github.com/smartcontractkit/package",
			},
		},
	}

	_, _, err := pkgsInfo.Get("github.com/smartcontractkit/package")
	require.NoError(t, err, "error getting existing package")

	_, _, err = pkgsInfo.Get("github.com/smartcontractkit/non_existent")
	require.Error(t, err, "error getting non-existent package")
}
