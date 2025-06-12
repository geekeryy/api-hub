package svc

import (
	{{.configImport}}
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
