func (m *default{{.upperStartCamelObject}}Model) Updates(session *gorm.DB, newData interface{}) error {
	return session.Model(&{{.upperStartCamelObject}}{}).Updates(newData).Error
}
