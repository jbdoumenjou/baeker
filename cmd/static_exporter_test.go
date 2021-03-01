package cmd

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traefik/traefik/v2/pkg/config/static"
	"github.com/traefik/traefik/v2/pkg/provider/docker"
	"github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd"
	"github.com/traefik/traefik/v2/pkg/types"
)

func TestTOMLExport(t *testing.T) {
	t.Parallel()
	configuration := &static.Configuration{
		Log: &types.TraefikLog{
			Level: "debug",
		},
	}
	exportedConf := new(bytes.Buffer)
	err := ExportToml(configuration, exportedConf)
	require.NoError(t, err)

	expectedConf, err := ioutil.ReadFile(filepath.FromSlash("./fixtures/static.toml"))

	require.NoError(t, err)

	assert.Equal(t, string(expectedConf), exportedConf.String())
}

func TestYAMLExport(t *testing.T) {
	t.Parallel()
	configuration := &static.Configuration{
		Log: &types.TraefikLog{
			Level: "debug",
		},
	}
	exportedConf := new(bytes.Buffer)
	err := ExportYaml(configuration, exportedConf)
	require.NoError(t, err)

	expectedConf, err := ioutil.ReadFile(filepath.FromSlash("./fixtures/static.yml"))

	require.NoError(t, err)

	assert.Equal(t, string(expectedConf), exportedConf.String())
}

func TestCLIExport(t *testing.T) {
	t.Parallel()
	configuration := &static.Configuration{
		Log: &types.TraefikLog{
			Level: "DEBUG",
		},
	}
	exportedConf := new(bytes.Buffer)
	err := ExportCLI(configuration, exportedConf)
	require.NoError(t, err)

	expectedConf, err := ioutil.ReadFile(filepath.FromSlash("./fixtures/static.cli"))

	require.NoError(t, err)

	assert.Equal(t, string(expectedConf), exportedConf.String())
}

func TestKubernetesExport(t *testing.T) {
	t.Parallel()
	configuration := &static.Configuration{
		EntryPoints: map[string]*static.EntryPoint{
			"web":       {Address: ":8000"},
			"websecure": {Address: ":8443"},
		},
		Providers: &static.Providers{
			KubernetesCRD: &crd.Provider{},
		},
	}
	exportedConf := new(bytes.Buffer)
	err := ExportKubernetes(configuration, "traefik-lb-svc-tpl.yml", exportedConf)
	require.NoError(t, err)

	expectedConf, err := ioutil.ReadFile(filepath.FromSlash("./fixtures/traefik-lb-svc.yml"))

	require.NoError(t, err)

	assert.Equal(t, string(expectedConf), exportedConf.String())
}

func TestDockerExport(t *testing.T) {
	t.Parallel()
	configuration := &static.Configuration{
		EntryPoints: map[string]*static.EntryPoint{
			"web":       {Address: ":8000"},
			"websecure": {Address: ":8443"},
		},
		Providers: &static.Providers{
			Docker: &docker.Provider{},
		},
	}
	exportedConf := new(bytes.Buffer)
	err := ExportKubernetes(configuration, "docker-compose-tpl.yml", exportedConf)
	require.NoError(t, err)

	expectedConf, err := ioutil.ReadFile(filepath.FromSlash("./fixtures/docker-compose.yml"))

	require.NoError(t, err)

	assert.Equal(t, string(expectedConf), exportedConf.String())
}
