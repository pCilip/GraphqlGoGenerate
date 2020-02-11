package HttpProvider

import (
	"bytes"
	"encoding/json"
	"github.com/pCilip/GraphqlGoGenerate/internal/Schema"
	"io/ioutil"
	"net/http"
)

type Provider struct {
	HttpEndpoint string
}

type graphQlData struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func (p *Provider) MustProvide() Schema.IntrospectionData {
	data := graphQlData{
		Query: IntrospectionQuery,
	}

	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(data); err != nil {
		panic("Error encoding introspection query to json")
	}

	resp, err := http.Post(p.HttpEndpoint, "application/json", body)

	if err != nil {
		panic(err)
	}

	if resp == nil {
		panic("response not set")
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	out := Schema.IntrospectionData{}

	err = json.Unmarshal(responseBody, &out)

	if err != nil {
		panic(err)
	}

	return out
}
