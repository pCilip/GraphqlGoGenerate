package Schema

type TypeRef struct {
	Kind   TypeKind
	Name   *string
	OfType struct {
		Kind   TypeKind
		Name   *string
		OfType struct {
			Kind   TypeKind
			Name   *string
			OfType struct {
				Kind TypeKind
				Name *string
			}
		}
	}
}

func (t *TypeRef) GetKind() TypeKind {
	return t.Kind
}

func (t *TypeRef) GetName() *string {
	return t.Name
}

func (t *TypeRef) SubType() TypeReference {
	if t.Kind == "" {
		return nil
	}

	typeRef := &TypeRef{
		Kind: t.OfType.Kind,
		Name: t.OfType.Name,
		OfType: struct {
			Kind   TypeKind
			Name   *string
			OfType struct {
				Kind   TypeKind
				Name   *string
				OfType struct {
					Kind TypeKind
					Name *string
				}
			}
		}{
			Kind: t.OfType.OfType.Kind,
			Name: t.OfType.OfType.Name,
			OfType: struct {
				Kind   TypeKind
				Name   *string
				OfType struct {
					Kind TypeKind
					Name *string
				}
			}{
				Kind: t.OfType.OfType.OfType.Kind,
				Name: t.OfType.OfType.OfType.Name,
				OfType: struct {
					Kind TypeKind
					Name *string
				}{
					Kind: "",
					Name: nil,
				},
			},
		},
	}

	var out TypeReference
	out = typeRef

	return out
}
