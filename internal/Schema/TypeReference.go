package Schema

type TypeReference interface {
	GetKind() TypeKind
	GetName() *string
	SubType() TypeReference
}
