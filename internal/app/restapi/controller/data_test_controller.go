package controller

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/minipkg/log"
	"github.com/yaruz/app/pkg/yarus_platform"
)

type dataTestController struct {
	Logger        log.ILogger
	yaruzPlatform yarus_platform.IPlatform
}

// RegisterHandlers sets up the routing of the HTTP handlers.
//	GET /api/models/ - список всех моделей
//	GET /api/model/{ID} - детали модели
func RegisterDataTestHandlers(r *routing.RouteGroup, yaruzPlatform yarus_platform.IPlatform, logger log.ILogger, authHandler routing.Handler) {
	c := dataTestController{
		Logger:        logger,
		yaruzPlatform: yaruzPlatform,
	}

	r.Get("/entity", c.entity)
}

func (c dataTestController) entity(ctx *routing.Context) error {
	res := make([]map[string]interface{}, 0, 10)
	res = append(res, map[string]interface{}{"test": "text-source"})
	cntx := ctx.Request.Context()

	entity := c.yaruzPlatform.ReferenceSubsystem().TextSource.Service.NewEntity()

	err := c.yaruzPlatform.ReferenceSubsystem().TextSource.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"1. errCreate": err.Error()})
	}

	val := c.yaruzPlatform.ReferenceSubsystem().TextValue.Service.NewEntity()
	val.TextSourceID = entity.ID
	val.TextLangID = 1
	val.Value = "Тестовое значение текстового поля"

	err = c.yaruzPlatform.ReferenceSubsystem().TextValue.Service.Create(cntx, val)
	if err != nil {
		res = append(res, map[string]interface{}{"1. errCreate": err.Error()})
	}

	e, err := c.yaruzPlatform.ReferenceSubsystem().TextSource.Service.TGet(cntx, entity.ID, val.TextLangID)
	if err != nil {
		res = append(res, map[string]interface{}{"4. errCreate1": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"4. entity1": e})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().TextSource.Service.Delete(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"7. errDelete": err.Error()})
	}

	return ctx.Write(res)
}
