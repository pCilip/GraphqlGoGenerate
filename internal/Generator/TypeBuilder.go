package Generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/pCilip/GraphqlGoGenerate/internal/Schema"
	"github.com/pCilip/GraphqlGoGenerate/internal/Utils"
)

type EndType int

const (
	RETURN EndType = iota
	CONTINUE
	OK
)

func BuildParamType(generator *Generator, parentPackage string, inputValue Schema.InputValue, generatedField *jen.Statement) EndType {
	return BuildTypeFromReference(generator, parentPackage, &inputValue.Type, generatedField, true)
}

func BuildFieldType(generator *Generator, parentPackage string, field Schema.Field, generatedField *jen.Statement) EndType {
	return BuildTypeFromReference(generator, parentPackage, &field.Type, generatedField, true)
}

func BuildTypeFromReference(generator *Generator, parentPackage string, typeRef Schema.TypeReference, generatedField *jen.Statement, addNull bool) EndType {

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
		return BuildTypeFromReference(generator, parentPackage, typeRef.SubType(), generatedField, true)
	// nullable type
	case Schema.OBJECT:
		return GenerateName(generator, parentPackage, generatedField, typeRef.GetName())

	case Schema.SCALAR, Schema.ENUM, Schema.INPUT_OBJECT:
		if addNull {
			generatedField.Id("*")
		}

		return GenerateName(generator, parentPackage, generatedField, typeRef.GetName())
	// non nullable type
	case Schema.NON_NULL:
		return BuildTypeFromReference(generator, parentPackage, typeRef.SubType(), generatedField, false)

	default:
		return CONTINUE
	}
}

func GenerateName(generator *Generator, parentPackage string, generatedField *jen.Statement, name *string) EndType {

	if name == nil {
		return RETURN
	}

	typeVal, ok := Utils.ToGoType(*name)
	if ok {
		importPath, found := generator.findImportPath(typeVal)
		if found {
			// importing model to same model (parent-child relations)
			// import without global path
			if importPath == parentPackage {
				generatedField.Qual("", typeVal)
			} else {
				generatedField.Qual(importPath, typeVal)
			}

		} else {
			return RETURN
		}
	} else {
		// TODO #1 fix
		if typeVal == "graphql.ID" {
			generatedField.Qual("github.com/graph-gophers/graphql-go", "ID")
		} else {
			generatedField.Id(typeVal)
		}

	}
	return OK
}
