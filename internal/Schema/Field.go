package Schema

type Field struct {
	Name              string
	Description       *string
	Args              []InputValue
	Type              TypeRef
	IsDeprecated      bool
	DeprecationReason *string
}
