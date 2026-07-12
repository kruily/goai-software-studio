func (m *default{{.upperStartCamelObject}}Model) Create(session *gorm.DB, data *{{.upperStartCamelObject}}) error {
	return session.Create(data).Error
}

func (m *default{{.upperStartCamelObject}}Model) CreateBatch(session *gorm.DB, data []*{{.upperStartCamelObject}}) error {
	return session.CreateInBatches(data,100).Error
}