func (m *default{{.upperStartCamelObject}}Model) FindOneBy{{.upperField}}(session *gorm.DB, {{.in}}) (*Custom{{.upperStartCamelObject}}, error) {
	var resp Custom{{.upperStartCamelObject}}
	err := session.Model(&Custom{{.upperStartCamelObject}}{}).Last(&resp,  "{{.originalField}}",{{.lowerStartCamelField}} ).Error
	if err != nil {
		return &Custom{{.upperStartCamelObject}}{
			{{.upperStartCamelObject}}: &{{.upperStartCamelObject}}{},
		},err
	} else {
		return &resp,nil
	}
}

func (m *default{{.upperStartCamelObject}}Model) FindOneBy{{.upperField}}ForUpdate(session *gorm.DB, {{.in}}) (*Custom{{.upperStartCamelObject}}, error) {
	var resp Custom{{.upperStartCamelObject}}
	err := session.Model(&Custom{{.upperStartCamelObject}}{}).Clauses(clause.Locking{Strength: "UPDATE"}).Last(&resp,  "{{.originalField}}",{{.lowerStartCamelField}} ).Error
	if err != nil {
		return &Custom{{.upperStartCamelObject}}{
			{{.upperStartCamelObject}}: &{{.upperStartCamelObject}}{},
		},err
	} else {
		return &resp,nil
	}
}

func (m *default{{.upperStartCamelObject}}Model) UpdateBy{{.upperField}}(session *gorm.DB, {{.in}}, newData interface{}) error {
	return session.Model(&{{.upperStartCamelObject}}{}).Where("{{.originalField}}",{{.lowerStartCamelField}}).Updates(newData).Error
}

func (m *default{{.upperStartCamelObject}}Model) DeleteBy{{.upperField}}(session *gorm.DB, {{.in}}) error {
	return session.Where("{{.originalField}}",{{.lowerStartCamelField}}).Delete(&{{.upperStartCamelObject}}{}).Error
}
