package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traefik/traefik/v2/pkg/config/static"
	"github.com/traefik/traefik/v2/pkg/types"
)

func TestImportToml(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		description string
		filePath    string
		expected    *static.Configuration
		err         bool
	}{
		{
			description: "empty",
			filePath:    "./fixtures/empty.toml",
			expected:    &static.Configuration{},
		},
		{
			description: "Use valid file",
			filePath:    "./fixtures/static.toml",
			expected: &static.Configuration{
				Log: &types.TraefikLog{
					Level: "debug",
				},
			},
		},
		{
			description: "Use unsopported format (yaml instead of toml)",
			filePath:    "./fixtures/static.yml",
			err:         true,
		},
	}

	for _, test := range testcases {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			confReader, err := os.Open(filepath.FromSlash(test.filePath))
			require.NoError(t, err)

			conf, err := ImportToml(confReader)
			if test.err {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expected, conf)
		})
	}
}

func TestImportYaml(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		description string
		filePath    string
		expected    *static.Configuration
		err         bool
	}{
		{
			description: "empty",
			filePath:    "./fixtures/empty.yml",
			expected:    &static.Configuration{},
		},
		{
			description: "Use valid file",
			filePath:    "./fixtures/static.yml",
			expected: &static.Configuration{
				Log: &types.TraefikLog{
					Level: "debug",
				},
			},
		},
		{
			description: "Use unsopported format (toml instead of yaml)",
			filePath:    "./fixtures/static.toml",
			err:         true,
		},
	}

	for _, test := range testcases {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			confReader, err := os.Open(filepath.FromSlash(test.filePath))
			require.NoError(t, err)

			conf, err := ImportYaml(confReader)
			if test.err {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expected, conf)
		})
	}
}
