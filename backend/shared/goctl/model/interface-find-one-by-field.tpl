FindOneBy{{.upperField}}(session *gorm.DB, {{.in}}) (*Custom{{.upperStartCamelObject}}, error)
FindOneBy{{.upperField}}ForUpdate(session *gorm.DB, {{.in}}) (*Custom{{.upperStartCamelObject}}, error)
UpdateBy{{.upperField}}(session *gorm.DB, {{.in}}, newData interface{}) error
DeleteBy{{.upperField}}(session *gorm.DB, {{.in}}) error