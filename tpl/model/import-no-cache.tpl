import (
	"context"
	"database/sql"
	"github.com/SpectatorNan/gorm-zero/gormc"
	{{if .time}}"time"{{end}}

	"gorm.io/gorm"
    "github.com/SpectatorNan/gorm-zero/gormc/pagex"
	{{if .third}}{{.third}}{{end}}
)
