package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traefik/traefik/v2/pkg/config/static"
)

func TestAddKubernetesProvider(t *testing.T) {
	t.Parallel()
	builder, err := NewStaticConfBuilder().AddKubernetesProvider()
	require.NoError(t, err)

	configuration := builder.GetConfiguration()
	assert.NotNil(t, configuration.Providers)
	assert.NotNil(t, configuration.Providers.KubernetesCRD)
}

func TestAddDockerProvider(t *testing.T) {
	t.Parallel()
	builder, err := NewStaticConfBuilder().AddDockerProvider()
	require.NoError(t, err)

	configuration := builder.GetConfiguration()
	assert.NotNil(t, configuration.Providers)
	assert.NotNil(t, configuration.Providers.Docker)
}

func TestAddFileProvider(t *testing.T) {
	t.Parallel()
	directory := "/foo/bar"
	builder, err := NewStaticConfBuilder().AddFileProvider(directory)
	require.NoError(t, err)

	configuration := builder.GetConfiguration()
	assert.NotNil(t, configuration.Providers)
	assert.NotNil(t, configuration.Providers.File)
	assert.Equal(t, directory, configuration.Providers.File.Directory)
}

func TestAddEntryPoint(t *testing.T) {
	t.Parallel()
	builder, err := NewStaticConfBuilder().AddEntryPoint("web", ":8000")
	require.NoError(t, err)
	_, err = builder.AddEntryPoint("websecure", ":8443")
	require.NoError(t, err)

	configuration := builder.GetConfiguration()
	require.Equal(t, 2, len(configuration.EntryPoints))

	ep, ok := configuration.EntryPoints["web"]
	require.True(t, ok)
	assert.Equal(t, &static.EntryPoint{Address: ":8000"}, ep)

	ep, ok = configuration.EntryPoints["websecure"]
	require.True(t, ok)
	assert.Equal(t, &static.EntryPoint{Address: ":8443"}, ep)
}
