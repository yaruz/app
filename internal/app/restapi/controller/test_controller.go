package controller

import (
	"fmt"
	"strconv"

	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform"

	"github.com/minipkg/log"

	routing "github.com/go-ozzo/ozzo-routing/v2"
)

type testController struct {
	Logger        log.ILogger
	yaruzPlatform yarus_platform.IPlatform
}

// RegisterHandlers sets up the routing of the HTTP handlers.
//	GET /api/models/ - список всех моделей
//	GET /api/model/{ID} - детали модели
func RegisterTestHandlers(r *routing.RouteGroup, yaruzPlatform yarus_platform.IPlatform, logger log.ILogger, authHandler routing.Handler) {
	c := testController{
		Logger:        logger,
		yaruzPlatform: yaruzPlatform,
	}

	r.Get("/test", c.Test)
	r.Get("/property-unit", c.PropertyUnit)
	r.Get("/property-group", c.PropertyGroup)
	r.Get("/property-view-type", c.PropertyViewType)
	r.Get("/property-type", c.propertyType)
	r.Get("/property", c.property)
}

func (c testController) Test(ctx *routing.Context) error {
	return c.yaruzPlatform.ReferenceSubsystem().PropertyUnit.Repository.Test()

}

func (c testController) PropertyUnit(ctx *routing.Context) error {
	res := make([]map[string]interface{}, 0, 10)
	res = append(res, map[string]interface{}{"test": "property-unit"})
	cntx := ctx.Request.Context()

	entity := c.yaruzPlatform.ReferenceSubsystem().PropertyUnit.Service.NewEntity()
	entity.Sysname = "WrongName"

	err := entity.Validate()
	if err != nil {
		res = append(res, map[string]interface{}{"1. errValidate": err.Error()})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().PropertyUnit.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"2. errCreate1": err.Error()})
	}

	entity.Sysname = "property_unit"
	err = c.yaruzPlatform.ReferenceSubsystem().PropertyUnit.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"3. errCreate": err.Error()})
	}

	e, err := c.yaruzPlatform.ReferenceSubsystem().PropertyUnit.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"4. errCreate1": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"4. entity1": fmt.Sprintf("%#v", e)})
	}

	entity.Sysname = "property_unit_" + strconv.Itoa(int(entity.ID))
	err = c.yaruzPlatform.ReferenceSubsystem().PropertyUnit.Service.Update(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"5. errCreate": err.Error()})
	}

	e, err = c.yaruzPlatform.ReferenceSubsystem().PropertyUnit.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"6. errGet": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"6. entity2": fmt.Sprintf("%#v", e)})
	}

	list, err := c.yaruzPlatform.ReferenceSubsystem().PropertyUnit.Service.Query(cntx, &selection_condition.SelectionCondition{
		Where: selection_condition.WhereCondition{
			Field:     "Sysname",
			Condition: "eq",
			Value:     entity.Sysname,
		},
	})
	if err != nil {
		res = append(res, map[string]interface{}{"6. errGet": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"6. list": fmt.Sprintf("%#v", list)})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().PropertyUnit.Service.Delete(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"7. errDelete": err.Error()})
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
