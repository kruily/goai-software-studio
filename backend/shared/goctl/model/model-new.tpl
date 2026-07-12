func new{{.upperStartCamelObject}}Model() *default{{.upperStartCamelObject}}Model {
	return &default{{.upperStartCamelObject}}Model{
	}
}

func (m *Custom{{.upperStartCamelObject}}) TableName() string {
	return {{.table}}
}
