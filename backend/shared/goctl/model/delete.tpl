func (m *default{{.upperStartCamelObject}}Model) Update(session *gorm.DB, {{.lowerStartCamelPrimaryKey}} {{.dataType}}, newData interface{}) error {
	return session.Model(&{{.upperStartCamelObject}}{}).Where("{{.originalPrimaryKey}} = ?",{{.lowerStartCamelPrimaryKey}}).Updates(newData).Error
}

func (m *default{{.upperStartCamelObject}}Model) Deletes(session *gorm.DB) error {
	return session.Delete(&{{.upperStartCamelObject}}{}).Error
}

func (m *default{{.upperStartCamelObject}}Model) Delete(session *gorm.DB, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error {
	return session.Delete(&{{.upperStartCamelObject}}{},{{.lowerStartCamelPrimaryKey}}).Error
}

func (m *default{{.upperStartCamelObject}}Model) Count(session *gorm.DB) (count int64) {
	session.Model(&{{.upperStartCamelObject}}{}).Count(&count)
	return
}
