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
		desc          string
		conf          static.Configuration
		templatePath  string
		expected      string
		expectedError bool
	}{
		{
			desc:          "unknown template",
			conf:          GetDefaultConf("docker"),
			templatePath:  "unknown",
			expectedError: true,
		},
		{
			desc: "bad entrypoint syntax",
			conf: static.Configuration{EntryPoints: map[string]*static.EntryPoint{
				"test": {Address: "bad syntax"},
			}},
			templatePath:  "docker-compose-tpl.yml",
			expectedError: true,
		},
		{
			desc:         "Default docker configuration",
			conf:         GetDefaultConf("docker"),
			templatePath: "docker-compose-tpl.yml",
			expected:     "docker-compose.yml",
		},
		{
			desc:         "Default kubernetes configuration",
			conf:         GetDefaultConf("kubernetes"),
			templatePath: "traefik-lb-svc-tpl.yml",
			expected:     "traefik-lb-svc.yml",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			exportedConf := new(bytes.Buffer)
			err := ExportConf(test.conf, test.templatePath, exportedConf)
			if test.expectedError {
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
