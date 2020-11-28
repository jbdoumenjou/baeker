package cmd

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traefik/traefik/v2/pkg/config/static"
)

func TestExportConf(t *testing.T) {
	testCases := []struct {
		desc     string
		conf     static.Configuration
		expected string
		err      bool
	}{
		{
			"Default docker configuration",
			GetDefaultConf("docker"),
			"docker-compose.yml",
			false,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			exportedConf := new(bytes.Buffer)
			err := ExportConf(test.conf, "./docker-compose-tpl.yaml", exportedConf)
			if test.err {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			expectedConf, err := ioutil.ReadFile(filepath.FromSlash("./fixtures/" + test.expected))
			require.NoError(t, err)

			assert.Equal(t, string(expectedConf), exportedConf.String())
		})
	}

}
