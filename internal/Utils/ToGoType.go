package Utils

func ToGoType(val string) (string, bool) {
	switch val {
	case "String":
		return "string", false
	case "Int":
		return "int32", false
	case "Boolean":
		return "bool", false
	case "Float":
		return "float64", false
	case "ID":
		// TODO #1 hotfix
		return "graphql.ID", false
	default:
		return val, true
	}
}
