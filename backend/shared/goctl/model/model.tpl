package {{.pkg}}
import (
	"gorm.io/gorm"
)
var _ {{.upperStartCamelObject}}Model = (*custom{{.upperStartCamelObject}}Model)(nil)

type (
	// {{.upperStartCamelObject}}Model is an interface to be customized, add more methods here,
	// and implement the added methods in custom{{.upperStartCamelObject}}Model.
	{{.upperStartCamelObject}}Model interface {
		{{.lowerStartCamelObject}}Model
	}

	custom{{.upperStartCamelObject}}Model struct {
		*default{{.upperStartCamelObject}}Model
	}

	Custom{{.upperStartCamelObject}} struct {
		*{{.upperStartCamelObject}}
	}
)

// New{{.upperStartCamelObject}}Model returns a model for the database table.
func New{{.upperStartCamelObject}}Model() {{.upperStartCamelObject}}Model {
	return &custom{{.upperStartCamelObject}}Model{
		default{{.upperStartCamelObject}}Model: new{{.upperStartCamelObject}}Model(),
	}
}

func NewCustom{{.upperStartCamelObject}}() *Custom{{.upperStartCamelObject}} {
	return &Custom{{.upperStartCamelObject}}{
		{{.upperStartCamelObject}}: new({{.upperStartCamelObject}}),
	}
}

func (m *Custom{{.upperStartCamelObject}}) AfterFind(*gorm.DB) (err error) {
	return nil
}