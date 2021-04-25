package controller

import (
	"fmt"

	"github.com/yaruz/app/pkg/yarus_platform"

	"github.com/minipkg/log"

	routing "github.com/go-ozzo/ozzo-routing/v2"
)

type testController struct {
	Logger          log.ILogger
	yaruzRepository yarus_platform.IPlatform
}

// RegisterHandlers sets up the routing of the HTTP handlers.
//	GET /api/models/ - список всех моделей
//	GET /api/model/{ID} - детали модели
func RegisterTestHandlers(r *routing.RouteGroup, yaruzRepository yarus_platform.IPlatform, logger log.ILogger, authHandler routing.Handler) {
	c := testController{
		Logger:          logger,
		yaruzRepository: yaruzRepository,
	}

	r.Get("/property-unit", c.PropertyUnit)
	r.Get("/property-group", c.PropertyGroup)
	r.Get("/property-view-type", c.PropertyViewType)
	r.Get("/property-type", c.propertyType)
	r.Get("/property", c.property)
}

func (c testController) PropertyUnit(ctx *routing.Context) error {
	res := make(map[string]interface{}, 10)
	res["test"] = "property-unit"
	cntx := ctx.Request.Context()

	entity := c.yaruzRepository.Reference().PropertyUnit.Service.NewEntity()
	entity.Sysname = "WrongName"

	err := entity.Validate()
	if err != nil {
		res["1. errValidate"] = err
	}

	err = c.yaruzRepository.Reference().PropertyUnit.Service.Create(cntx, entity)
	if err != nil {
		res["2. errCreate1"] = err.Error()
	}

	entity.Sysname = "property_unit_1"
	err = c.yaruzRepository.Reference().PropertyUnit.Service.Create(cntx, entity)
	if err != nil {
		res["3. errCreate"] = err.Error()
	}

	e, err := c.yaruzRepository.Reference().PropertyUnit.Service.Get(cntx, entity.ID)
	if err != nil {
		res["4. errCreate1"] = err.Error()
	} else {
		res["4. entity1"], _ = fmt.Printf("%#v", e)
	}

	entity.Sysname = "property_unit_1"
	err = c.yaruzRepository.Reference().PropertyUnit.Service.Update(cntx, entity)
	if err != nil {
		res["5. errCreate"] = err.Error()
	}

	e, err = c.yaruzRepository.Reference().PropertyUnit.Service.Get(cntx, entity.ID)
	if err != nil {
		res["6. errGet"] = err.Error()
	} else {
		res["6. entity2"], _ = fmt.Printf("%#v", e)
	}

	err = c.yaruzRepository.Reference().PropertyUnit.Service.Delete(cntx, entity.ID)
	if err != nil {
		res["7. errDelete"] = err.Error()
	}

	return ctx.Write(res)
}

func (c testController) PropertyGroup(ctx *routing.Context) error {

	return ctx.Write("test: property-group")
}

func (c testController) PropertyViewType(ctx *routing.Context) error {

	return ctx.Write("test: property-view-type")
}

func (c testController) propertyType(ctx *routing.Context) error {

	return ctx.Write("test: property-type")
}

func (c testController) property(ctx *routing.Context) error {

	return ctx.Write("test: property")
}
