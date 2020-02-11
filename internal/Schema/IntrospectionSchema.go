package Schema

type IntrospectionSchema struct {
	Schema Schema `json:"__schema"`
}

type Schema struct {
	QueryType        Type
	MutationType     Type
	SubscriptionType Type
	Types            []FullType
	// TODO directives??? for generation not necessary
}
