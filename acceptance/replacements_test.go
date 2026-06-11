package acceptance_test

import (
	"path/filepath"
	"regexp"
	"testing"

	"github.com/databricks/cli/internal/testutil"
	"github.com/databricks/cli/libs/testdiff"
	"github.com/stretchr/testify/assert"
)

func TestLoadUserReplacementsOverridesGenericNumberReplacement(t *testing.T) {
	tmpDir := t.TempDir()
	testutil.WriteFile(t, filepath.Join(tmpDir, userReplacementsFilename), "8781181929457644375:FOO_ID\n")

	repls := testdiff.ReplacementsContext{}
	repls.Repls = append(repls.Repls, testdiff.Replacement{
		Old:   regexp.MustCompile(`\d{17,}`),
		New:   "[NUMID]",
		Order: 10,
	})

	loadUserReplacements(t, &repls, tmpDir)

	got := repls.Replace(`{"job_id": 8781181929457644375}`)
	assert.Equal(t, `{"job_id": [FOO_ID]}`, got)
}
