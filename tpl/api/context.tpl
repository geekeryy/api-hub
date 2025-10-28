// Code scaffolded by goctl. Safe to edit.
// goctl {{.version}}

package svc

import (
	{{.configImport}}
	"github.com/zeromicro/go-zero/core/proc"

	"{{.projectPkg}}/core/validate"
)

type ServiceContext struct {
	Config {{.config}}
	{{.middleware}}
	Validator *validate.Validate
}

func NewServiceContext(c {{.config}}) *ServiceContext {
	return &ServiceContext{
		Config: c,
		{{.middlewareAssignment}}
		Validator: validate.New(nil, []string{"zh", "en"}),
	}
}

func (s *ServiceContext) Close() {
	// TODO graceful stop
}