package Generator

import (
	"GraphqlGoGenerate/internal/Schema"
	"GraphqlGoGenerate/internal/Utils"
	"github.com/dave/jennifer/jen"
)

type EndType int

const (
	RETURN EndType = iota
	CONTINUE
	OK
)

func BuildParamType(generator *Generator, inputValue Schema.InputValue, generatedField *jen.Statement) EndType {
	return BuildTypeFromReference(generator, &inputValue.Type, generatedField, true)
}

func BuildFieldType(generator *Generator, field Schema.Field, generatedField *jen.Statement) EndType {
	return BuildTypeFromReference(generator, &field.Type, generatedField, true)
}

func BuildTypeFromReference(generator *Generator, typeRef Schema.TypeReference, generatedField *jen.Statement, addNull bool) EndType {

	if typeRef == nil {
		return OK
	}

	switch typeRef.GetKind() {
	// non null any typeRef...
	case Schema.LIST:
		if addNull {
			generatedField.Id("*")
		}
		generatedField.Id("[]")
		return BuildTypeFromReference(generator, typeRef.SubType(), generatedField, true)
	// nullable type
	case Schema.SCALAR, Schema.ENUM, Schema.INPUT_OBJECT, Schema.OBJECT:
		if addNull {
			generatedField.Id("*")
		}

		return GenerateName(generator, generatedField, typeRef.GetName())
	// non nullable type
	case Schema.NON_NULL:
		return BuildTypeFromReference(generator, typeRef.SubType(), generatedField, false)

	default:
		return CONTINUE
	}
}

func GenerateName(generator *Generator, generatedField *jen.Statement, name *string) EndType {

	if name == nil {
		return RETURN
	}

	typeVal, ok := Utils.ToGoType(*name)
	if ok {
		importPath, found := generator.findImportPath(typeVal)
		if found {
			generatedField.Qual(importPath, typeVal)
		} else {
			return RETURN
		}
	} else {
		generatedField.Id(typeVal)
	}
	return OK
}
