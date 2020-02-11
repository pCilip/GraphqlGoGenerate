package Schema

type InputValue struct {
	Name         string
	Description  *string
	Type         TypeRef
	DefaultValue *string
}
