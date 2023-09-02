package templator

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTemplateEnv(t *testing.T) {
	sample := path.Join(basepath, "testdata/sample.template")
	f, err := os.Open(sample)
	require.NoError(t, err)
	defer f.Close()

	Debug = true
	err = TemplateIOEnv(f, os.Stdout, "", "")
	require.NoError(t, err)
}
