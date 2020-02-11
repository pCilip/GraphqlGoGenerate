package JsonProvider

import (
	"GraphqlGoGenerate/internal/Schema"
	"encoding/json"
	"io/ioutil"
)

type Provider struct {
	FilePath string
}

func (p *Provider) MustProvide() Schema.IntrospectionData {
	bytes, err := ioutil.ReadFile(p.FilePath)

	if err != nil {
		panic(err)
	}

	out := Schema.IntrospectionData{}

	err = json.Unmarshal(bytes, &out)

	if err != nil {
		panic(err)
	}

	return out
}
