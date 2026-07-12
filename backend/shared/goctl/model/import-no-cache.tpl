import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	{{if .time}}"time"{{end}}

    {{if .containsPQ}}"github.com/lib/pq"{{end}}
)
