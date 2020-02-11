package Schema

type TypeKind string

const (
	SCALAR       TypeKind = "SCALAR"
	OBJECT       TypeKind = "OBJECT"
	INTERFACE    TypeKind = "INTERFACE"
	UNION        TypeKind = "UNION"
	ENUM         TypeKind = "ENUM"
	INPUT_OBJECT TypeKind = "INPUT_OBJECT"
	LIST         TypeKind = "LIST"
	NON_NULL     TypeKind = "NON_NULL"
)

type Type struct {
	Name *string
}

type FullType struct {
	Type
	Kind          TypeKind
	Name          *string
	Description   *string
	Fields        []Field
	InputFields   []InputValue
	Interfaces    []TypeRef
	EnumValues    []EnumValue
	PossibleTypes []TypeRef
}
