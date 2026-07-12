func (m *default{{.upperStartCamelObject}}Model) FindOneById(session *gorm.DB, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*Custom{{.upperStartCamelObject}}, error) {
	var resp Custom{{.upperStartCamelObject}}
	err := session.Model(&Custom{{.upperStartCamelObject}}{}).Last(&resp, "{{.lowerStartCamelPrimaryKey}} = ?",{{.lowerStartCamelPrimaryKey}}).Error
	if err != nil {
		return &Custom{{.upperStartCamelObject}}{
			{{.upperStartCamelObject}}: &{{.upperStartCamelObject}}{},
		},err
	} else {
		return &resp,nil
	}
}

func (m *default{{.upperStartCamelObject}}Model) FindOneByIdForUpdate(session *gorm.DB, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*Custom{{.upperStartCamelObject}}, error) {
	var resp Custom{{.upperStartCamelObject}}
	err := session.Model(&Custom{{.upperStartCamelObject}}{}).Clauses(clause.Locking{Strength: "UPDATE"}).Last(&resp, {{.lowerStartCamelPrimaryKey}}).Error
	if err != nil {
		return &Custom{{.upperStartCamelObject}}{
			{{.upperStartCamelObject}}: &{{.upperStartCamelObject}}{},
		},err
	} else {
		return &resp,nil
	}
}

func (m *default{{.upperStartCamelObject}}Model) FindOne(session *gorm.DB) (*Custom{{.upperStartCamelObject}}, error) {
	var resp Custom{{.upperStartCamelObject}}
	err := session.Model(&Custom{{.upperStartCamelObject}}{}).Last(&resp).Error
	if err != nil {
		return &Custom{{.upperStartCamelObject}}{
			{{.upperStartCamelObject}}: &{{.upperStartCamelObject}}{},
		},err
	} else {
		return &resp,nil
	}
}

func (m *default{{.upperStartCamelObject}}Model) FindList(session *gorm.DB, limit ...int) (models []*Custom{{.upperStartCamelObject}}, err error) {

	if len(limit) > 0 && limit[0] > 0 {
		models = make([]*Custom{{.upperStartCamelObject}}, 0, limit[0])
	} else {
		limit = []int{-1}
		models = make([]*Custom{{.upperStartCamelObject}}, 0)
	}
	err = session.Model(&Custom{{.upperStartCamelObject}}{}).Limit(limit[0]).Find(&models).Error
	return
}

func (m *default{{.upperStartCamelObject}}Model) FindListForPage(session *gorm.DB, page int, pageSize int) (models []*Custom{{.upperStartCamelObject}}, err error) {
	models = make([]*Custom{{.upperStartCamelObject}}, 0, pageSize)
	err = session.Model(&Custom{{.upperStartCamelObject}}{}).Limit(pageSize).Offset((page - 1) * pageSize).Find(&models).Error
	return
}