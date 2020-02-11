package SchemaProvider

import "github.com/pCilip/GraphqlGoGenerate/internal/Schema"

type SchemaProvider interface {
	MustProvide() Schema.IntrospectionData
}
