package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	yaml3 "gopkg.in/yaml.v2"
	yaml2 "gopkg.in/yaml.v3"
	json2 "k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/yaml"
)

type A struct {
	ValueA int `json:"a"`
}

type Ab struct {
	ValueAb int `json:"ab"`
}

type B struct {
	A  `json:",inline" yaml:",inline"`
	Ab `json:""`
	M  map[string]string `json:",inline" yaml:",inline"`
}

func Test_Gopkg(t *testing.T) {
	b := &B{A: A{4}, Ab: Ab{42}, M: map[string]string{"one": "un"}}

	bytes, err := yaml2.Marshal(b)
	require.NoError(t, err)
	fmt.Println("With gopkg.in/yaml.v3")
	fmt.Printf("%s\n", string(bytes))

	bytes, err = yaml3.Marshal(b)
	require.NoError(t, err)
	fmt.Println("With gopkg.in/yaml.v2")
	fmt.Printf("%s\n", string(bytes))

	fmt.Println("With sigs.k8s.io/yaml")
	bytes, err = yaml.Marshal(b)
	require.NoError(t, err)
	fmt.Printf("%s\n", string(bytes))

	fmt.Println("With apimachinery json")
	bytes, err = json2.Marshal(b)
	require.NoError(t, err)
	fmt.Printf("%s\n", string(bytes))
}
