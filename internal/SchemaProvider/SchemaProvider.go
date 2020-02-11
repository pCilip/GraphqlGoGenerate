package SchemaProvider

import "GraphqlGoGenerate/internal/Schema"

type SchemaProvider interface {
	MustProvide() Schema.IntrospectionData
}
