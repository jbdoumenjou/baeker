package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traefik/traefik/v2/pkg/config/static"
)

func TestExportConf(t *testing.T) {
	t.Parallel()

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

func TestTag(t *testing.T) {
	//conf := static.EntryPoints{
	//	"web":       &static.EntryPoint{Address: "8000"},
	//	"websecure": &static.EntryPoint{Address: "8443"},
	//}

	value := reflect.ValueOf(static.EntryPoint{Address: ":8000", HTTP: static.HTTPConfig{
		Middlewares: []string{"test"},
	}})

	ty := value.Type()
	sf := make([]reflect.StructField, 0)
	for i := 0; i < ty.NumField(); i++ {
		fmt.Println(ty.Field(i).Tag)
		sf = append(sf, ty.Field(i))
		if ty.Field(i).Name == "TOTO" {
			sf[i].Tag = `toml:"-"`
		}
	}

	newType := reflect.StructOf(sf)
	newValue := value.Convert(newType)
	fmt.Printf("New object:\n%v\n", newValue.Interface())
	//fmt.Printf("New object:\n%v\n", newValue.Interface().(static.EntryPoint))
	//nv := newValue.Interface().(static.EntryPoint)
	//conf["web"] = &nv
	//st := reflect.TypeOf(ep)
	//field := st.Field(4)
	//fmt.Println(field.Tag.Get("toml"))

	exportedConf := new(bytes.Buffer)
	err := toml.NewEncoder(exportedConf).Encode(newValue.Interface())
	require.NoError(t, err)
	fmt.Println(exportedConf.String())

}

type LocalEntryPoint struct {
	Address          string                       `description:"Entry point address." json:"address,omitempty" toml:"address,omitempty" yaml:"address,omitempty"`
	Transport        *static.EntryPointsTransport `description:"Configures communication between clients and Traefik." json:"transport,omitempty" toml:"transport,omitempty" yaml:"transport,omitempty" export:"true"`
	ProxyProtocol    *static.ProxyProtocol        `description:"Proxy-Protocol configuration." json:"proxyProtocol,omitempty" toml:"proxyProtocol,omitempty" yaml:"proxyProtocol,omitempty" label:"allowEmpty" file:"allowEmpty" export:"true"`
	ForwardedHeaders *static.ForwardedHeaders     `description:"Trust client forwarding headers." json:"forwardedHeaders,omitempty" toml:"forwardedHeaders,omitempty" yaml:"forwardedHeaders,omitempty" export:"true"`
	HTTP             static.HTTPConfig            `description:"HTTP configuration." json:"-" toml:"-" yaml:"http,omitempty" export:"true"`
}
type LocalEntryPoints map[string]*LocalEntryPoint

func TestPouet(t *testing.T) {
	exportedConf := new(bytes.Buffer)

	defaultConf := GetDefaultConf("toml")
	//defaultConf := dynamic.Router{
	//	TLS: &dynamic.RouterTLSConfig{},
	//}
	myConf := static.Configuration{
		EntryPoints: static.EntryPoints{
			"ep1": &static.EntryPoint{
				Address: "8443",
				HTTP:    static.HTTPConfig{},
			},
		},
	}
	err := toml.NewEncoder(exportedConf).Encode(defaultConf)
	require.NoError(t, err)
	fmt.Println(exportedConf.String())
	fmt.Println("========================")
	err = toml.NewEncoder(exportedConf).Encode(myConf)
	require.NoError(t, err)
	fmt.Println(exportedConf.String())
}

func TestName(t *testing.T) {
	exportedConf := new(bytes.Buffer)
	//conf := GetDefaultConf("toml")
	//conf := dynamic.Router{
	//	TLS: &dynamic.RouterTLSConfig{},
	//}
	conf := static.EntryPoints{
		"web":       &static.EntryPoint{Address: "8000"},
		"websecure": &static.EntryPoint{Address: "8443"},
	}
	err := toml.NewEncoder(exportedConf).Encode(conf)
	require.NoError(t, err)
	fmt.Println(exportedConf.String())
	exportedConf.Reset()
	fmt.Println("======")

	localConf := LocalEntryPoints{
		"web":       &LocalEntryPoint{Address: "8000"},
		"websecure": &LocalEntryPoint{Address: "8443"},
	}
	err = toml.NewEncoder(exportedConf).Encode(localConf)
	require.NoError(t, err)
	fmt.Println(exportedConf.String())

	//data, err := yaml.Marshal(&conf)
	//require.NoError(t, err)
	//
	//fmt.Println(string(data))
	//fmt.Println("=======")
	//fmt.Println(exportedConf.String())
	//
	fmt.Println("====================")
	exportedConf2 := new(bytes.Buffer)
	conf2 := getEmptyConf(GetDefaultConf("toml"))
	toml.NewEncoder(exportedConf2).Encode(&conf2)
	fmt.Println(exportedConf2.String())

	type poub struct {
		Next *poub
		Name string
	}
	//poub := struct {
	//	Name string
	//}{}

	toml.NewEncoder(exportedConf2).Encode(&poub{Next: &poub{Name: "under"}})
	fmt.Println(exportedConf2.String())

}
