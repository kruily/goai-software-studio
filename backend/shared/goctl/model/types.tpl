type (
	{{.lowerStartCamelObject}}Model interface{
		{{.method}}
	}

	default{{.upperStartCamelObject}}Model struct {
	}

	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)
