package Generator

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	_ "github.com/graph-gophers/graphql-go"
	"github.com/iancoleman/strcase"
	Schema2 "github.com/pCilip/GraphqlGoGenerate/internal/Schema"
	"github.com/pCilip/GraphqlGoGenerate/internal/Utils"
	"os"
	path2 "path"
	"strings"
)

type Generator struct {
	ProjectName      string
	Schema           Schema2.IntrospectionSchema
	RootDirectory    string
	EnumDirectory    string
	ScalarsDirectory string
	InputsDirectory  string
	ObjectsDirectory string
	Imports          map[string]string
}

func NewGenerator(Schema Schema2.IntrospectionSchema, projectName string) *Generator {
	return &Generator{
		ProjectName:      projectName,
		Schema:           Schema,
		RootDirectory:    "./generated",
		ScalarsDirectory: "./generated/scalars",
		EnumDirectory:    "./generated/enums",
		InputsDirectory:  "./generated/inputs",
		ObjectsDirectory: "./generated/objects",
		Imports:          map[string]string{},
	}
}

func (generator *Generator) prepare() {
	err := os.RemoveAll(generator.RootDirectory)

	err = os.Mkdir(generator.RootDirectory, 0777)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(generator.EnumDirectory, 0777)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(generator.InputsDirectory, 0777)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(generator.ScalarsDirectory, 0777)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(generator.ObjectsDirectory, 0777)
	if err != nil {
		panic(err)
	}

	generator.addImportPath("graphql.ID", "graphql.Id")
}

func (generator *Generator) addImportPath(key, path string) {
	generator.Imports[key] = fmt.Sprintf("%s/%s", generator.ProjectName, path2.Dir(path))
}

func (generator *Generator) findImportPath(key string) (string, bool) {
	value, ok := generator.Imports[key]
	return value, ok
}

func (generator *Generator) RenderEnum(enum Schema2.FullType) {
	if enum.Name == nil {
		return
	}

	if len(enum.EnumValues) <= 0 {
		return
	}

	file := jen.NewFile(fmt.Sprintf("%s", *enum.Name))
	file.Comment("Code generated by GraphqlGenerator. DO NOT EDIT.")

	file.Type().Id(*enum.Name).String().Line()

	for _, value := range enum.EnumValues {
		// Value name must start with capital letter
		file.Const().Id(strcase.ToCamel(value.Name)).Id(*enum.Name).Op("=").Lit(value.Name)
	}

	filePath := fmt.Sprintf("%s/%s/%s.go", generator.EnumDirectory, *enum.Name, *enum.Name)
	generator.addImportPath(*enum.Name, filePath)

	Utils.EnsureDir(filePath)
	systemFile, err := os.Create(filePath)

	if err != nil {
		panic(err)
	}

	err = file.Render(systemFile)

	if err != nil {
		panic(err)
	}
}

func (generator *Generator) RenderScalar(scalar Schema2.FullType) {
	if scalar.Name == nil {
		return
	}
	file := jen.NewFile(fmt.Sprintf("%s", *scalar.Name))
	file.Comment("Code generated by GraphqlGenerator. DO NOT EDIT.")

	file.Type().Id(*scalar.Name).Struct(
		jen.Id("Data").String(),
		jen.Id("MarshallFunction").Func().Parens(nil).Params(jen.Id("[]byte"), jen.Qual("", "error")),
	).Line()

	// Implements graphql Type
	file.Func().
		Parens(jen.Id("_").Id(*scalar.Name)).
		Id("ImplementsGraphQLType").
		Params(jen.Id("name").Id("string")).
		Id("bool").
		Block(jen.Id("return").Id("name").Op("==").Lit(*scalar.Name)).
		Line()

	// UnmarshalGraphQL
	file.Func().
		Parens(jen.Id("id").Id("*").Id(*scalar.Name)).
		Id("UnmarshalGraphQL").                    // function name
		Params(jen.Id("input").Id("interface{}")). // parameters
		Qual("", "error").                         // return type
		Block(
			jen.Switch(jen.Id("input.(type)")).
				Block(
					jen.Case(jen.Id("string")),
					jen.Id("*id").Op("=").Id(*scalar.Name).
						Block(jen.Id("Data:"), jen.Id("input.(string),")).
						Line().
						Return(jen.Nil()),
					jen.Case(jen.Id("int")),
					jen.Id("*id").Op("=").Id(*scalar.Name).
						Block(jen.Id("Data:"), jen.Qual("strconv", "FormatInt").Parens(jen.Id("int64(input.(int)), 10")).Id(",")).
						Line().
						Return(jen.Nil()),
					jen.Case(jen.Id("int64")),
					jen.Id("*id").Op("=").Id(*scalar.Name).
						Block(jen.Id("Data:"), jen.Qual("strconv", "FormatInt").Parens(jen.Id("input.(int64), 10")).Id(",")).
						Line().
						Return(jen.Nil()),
					jen.Case(jen.Id("float64")),
					jen.Id("*id").Op("=").Id(*scalar.Name).
						Block(jen.Id("Data:"), jen.Qual("strconv", "FormatFloat").Parens(jen.Id("float64(input.(float64)), 'g', -1, 64")).Id(",")).
						Line().
						Return(jen.Nil()),
				),
			jen.Return(jen.Qual("errors", "").Id("New").Parens(jen.Lit(fmt.Sprintf("cannot unmarshal: %s", *scalar.Name))))).
		Line()

	file.Func().
		Parens(jen.Id("id").Id(*scalar.Name)).
		Id("MarshalJSON"). // function name
		Params().
		Params(jen.Id("[]byte"), jen.Qual("", "error")).
		Block(
			jen.If(jen.Id("id.MarshallFunction").Op("!=").Nil()).
				Block(
					jen.Return(jen.Id("id.MarshallFunction()")),
				).
				Else().
				Block(
					jen.Return(jen.Qual("encoding/json", "Marshal").Parens(jen.Id("id.Data"))),
				),
		)

	filePath := fmt.Sprintf("%s/%s/%s.go", generator.ScalarsDirectory, *scalar.Name, *scalar.Name)
	generator.addImportPath(*scalar.Name, filePath)

	Utils.EnsureDir(filePath)
	systemFile, err := os.Create(filePath)

	if err != nil {
		panic(err)
	}

	err = file.Render(systemFile)

	if err != nil {
		panic(err)
	}
}

func (generator *Generator) Render(object Schema2.FullType, dir string) bool {
	filePath := fmt.Sprintf("%s/%s/%s.go", dir, *object.Name, *object.Name)
	generator.addImportPath(*object.Name, filePath)

	file := jen.NewFile(fmt.Sprintf("%s", *object.Name))
	file.PackageComment("Code generated by GraphqlGenerator. DO NOT EDIT.")

	var fields []jen.Code
	for _, value := range object.InputFields {
		field := jen.Id(strcase.ToCamel(value.Name))

		endType := BuildParamType(generator, value, field)

		switch endType {
		case RETURN:
			return false
		case CONTINUE:
			continue
		default:
		}

		fields = append(fields, field)
	}
	file.Type().Id(*object.Name).Struct(fields...)

	Utils.EnsureDir(filePath)
	systemFile, err := os.Create(filePath)

	if err != nil {
		panic(err)
	}

	err = file.Render(systemFile)

	if err != nil {
		panic(err)
	}

	return true
}

func (generator *Generator) RenderObject(object Schema2.FullType, dir string) bool {
	filePath := fmt.Sprintf("%s/%s/%s.go", dir, *object.Name, *object.Name)
	generator.addImportPath(*object.Name, filePath)

	file := jen.NewFile(fmt.Sprintf("%s", *object.Name))
	file.PackageComment("Code generated by GraphqlGenerator. DO NOT EDIT.")

	var fields []jen.Code
	for _, value := range object.Fields {
		field := jen.Id(strcase.ToCamel(value.Name))

		var params []jen.Code

		for _, param := range value.Args {
			genParam := jen.Id(strcase.ToCamel(param.Name))

			endType := BuildParamType(generator, param, genParam)

			switch endType {
			case RETURN:
				return false
			case CONTINUE:
				continue
			default:
			}

			params = append(params, genParam)
		}
		// first param ctx and + wrapper around params
		ctxParam := jen.Id("ctx").Qual("context", "Context")

		var fieldParams []jen.Code

		fieldParams = append(fieldParams, ctxParam)

		if len(params) > 0 {
			args := jen.Id("args").Struct(params...)
			fieldParams = append(fieldParams, args)
		}

		field.Params(fieldParams...)

		returnTypes := jen.Id("")

		endType := BuildFieldType(generator, value, returnTypes)

		switch endType {
		case RETURN:
			return false
		case CONTINUE:
			continue
		default:
		}

		field.Parens(jen.List(returnTypes, jen.Qual("", "error")))
		fields = append(fields, field)
	}
	file.Type().Id(fmt.Sprintf("%s", *object.Name)).Interface(fields...)

	Utils.EnsureDir(filePath)
	systemFile, err := os.Create(filePath)

	if err != nil {
		panic(err)
	}

	err = file.Render(systemFile)

	if err != nil {
		panic(err)
	}

	return true
}

func (generator *Generator) RenderInputType(inputType Schema2.FullType) bool {
	if inputType.Name == nil {
		return true
	}

	if len(inputType.InputFields) <= 0 {
		return true
	}

	return generator.Render(inputType, generator.InputsDirectory)
}

func (generator *Generator) Generate() {
	generator.prepare()

	var enums []Schema2.FullType
	var inputs []Schema2.FullType
	var scalars []Schema2.FullType
	var objects []Schema2.FullType

	for _, itemType := range generator.Schema.Schema.Types {

		if itemType.Kind == Schema2.ENUM {
			enums = append(enums, itemType)
		}

		if itemType.Kind == Schema2.INPUT_OBJECT {
			inputs = append(inputs, itemType)
		}

		if itemType.Kind == Schema2.OBJECT {
			if itemType.Name != nil && strings.HasPrefix(*itemType.Name, "__") {
				continue
			}

			objects = append(objects, itemType)
		}

		if itemType.Kind == Schema2.SCALAR {
			if itemType.Name != nil {
				_, ok := Utils.ToGoType(*itemType.Name)
				if ok {
					scalars = append(scalars, itemType)
				}
			}
		}
	}

	for _, scalar := range scalars {
		generator.RenderScalar(scalar)
	}

	for _, enum := range enums {
		generator.RenderEnum(enum)
	}

	for len(inputs) > 0 {
		for i, inputType := range inputs {
			ok := generator.RenderInputType(inputType)

			if ok {
				inputs[i] = inputs[len(inputs)-1]
				inputs = inputs[:len(inputs)-1]
				break
			}
		}
	}

	for len(objects) > 0 {
		for i, inputType := range objects {
			ok := generator.RenderObject(inputType, generator.ObjectsDirectory)

			if ok {
				objects[i] = objects[len(objects)-1]
				objects = objects[:len(objects)-1]
				break
			}
		}
	}
}
