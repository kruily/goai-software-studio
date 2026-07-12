Update(session *gorm.DB,{{.lowerStartCamelPrimaryKey}} {{.dataType}}, newData interface{}) error
Deletes(session *gorm.DB) error
Delete(session *gorm.DB, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error
Count(session *gorm.DB) int64
