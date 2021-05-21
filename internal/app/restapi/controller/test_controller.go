package controller

import (
	"fmt"
	"strconv"
	"time"

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
	r.Get("/text-source", c.textSource)
	r.Get("/text-value", c.textValue)
	r.Get("/property-unit", c.propertyUnit)
	r.Get("/property-group", c.propertyGroup)
	r.Get("/property-view-type", c.propertyViewType)
	r.Get("/property-type", c.propertyType)
	r.Get("/property", c.property)
	r.Get("/entity-type", c.entityType)
}

func (c testController) Test(ctx *routing.Context) error {
	return c.yaruzPlatform.ReferenceSubsystem().PropertyUnit.Repository.Test()

}

func (c testController) textSource(ctx *routing.Context) error {
	res := make([]map[string]interface{}, 0, 10)
	res = append(res, map[string]interface{}{"test": "text-source"})
	cntx := ctx.Request.Context()

	entity := c.yaruzPlatform.ReferenceSubsystem().TextSource.Service.NewEntity()

	err := c.yaruzPlatform.ReferenceSubsystem().TextSource.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"1. errCreate": err.Error()})
	}

	e, err := c.yaruzPlatform.ReferenceSubsystem().TextSource.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"4. errCreate1": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"4. entity1": fmt.Sprintf("%#v", e)})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().TextSource.Service.Delete(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"7. errDelete": err.Error()})
	}

	return ctx.Write(res)
}

func (c testController) textValue(ctx *routing.Context) error {
	res := make([]map[string]interface{}, 0, 10)
	res = append(res, map[string]interface{}{"test": "text-value"})
	cntx := ctx.Request.Context()

	source := c.yaruzPlatform.ReferenceSubsystem().TextSource.Service.NewEntity()
	err := c.yaruzPlatform.ReferenceSubsystem().TextSource.Service.Create(cntx, source)
	if err != nil {
		res = append(res, map[string]interface{}{"1. errCreate": err.Error()})
	}

	entity := c.yaruzPlatform.ReferenceSubsystem().TextValue.Service.NewEntity()
	entity.TextSourceID = source.ID
	entity.LangID = 1
	entity.Value = "text value"
	err = c.yaruzPlatform.ReferenceSubsystem().TextValue.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"2. Create err": err.Error()})
	}

	e, err := c.yaruzPlatform.ReferenceSubsystem().TextValue.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"3. Get err": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"3. entity1": fmt.Sprintf("%#v", e)})
	}

	entity.Value = "updated text value - " + strconv.Itoa(int(entity.ID))
	err = c.yaruzPlatform.ReferenceSubsystem().TextValue.Service.Update(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"4. Update err": err.Error()})
	}

	e, err = c.yaruzPlatform.ReferenceSubsystem().TextValue.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"5. Get err": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"5. entity2": fmt.Sprintf("%#v", e)})
	}

	list, err := c.yaruzPlatform.ReferenceSubsystem().TextValue.Service.Query(cntx, &selection_condition.SelectionCondition{
		Where: selection_condition.WhereCondition{
			Field:     "Value",
			Condition: "eq",
			Value:     entity.Value,
		},
	})
	if err != nil {
		res = append(res, map[string]interface{}{"6. errGet": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"6. list": fmt.Sprintf("%#v", list)})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().TextValue.Service.Delete(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"7. errDelete": err.Error()})
	}

	return ctx.Write(res)
}

func (c testController) propertyUnit(ctx *routing.Context) error {
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

func (c testController) propertyGroup(ctx *routing.Context) error {
	res := make([]map[string]interface{}, 0, 10)
	res = append(res, map[string]interface{}{"test": "property-group"})
	cntx := ctx.Request.Context()

	entity := c.yaruzPlatform.ReferenceSubsystem().PropertyGroup.Service.NewEntity()
	entity.Sysname = "WrongName"

	err := entity.Validate()
	if err != nil {
		res = append(res, map[string]interface{}{"1. errValidate": err.Error()})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().PropertyGroup.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"2. errCreate1": err.Error()})
	}

	entity.Sysname = "property_group"
	err = c.yaruzPlatform.ReferenceSubsystem().PropertyGroup.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"3. errCreate": err.Error()})
	}

	e, err := c.yaruzPlatform.ReferenceSubsystem().PropertyGroup.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"4. errCreate1": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"4. entity1": fmt.Sprintf("%#v", e)})
	}

	entity.Sysname = "property_group_" + strconv.Itoa(int(entity.ID))
	err = c.yaruzPlatform.ReferenceSubsystem().PropertyGroup.Service.Update(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"5. errCreate": err.Error()})
	}

	e, err = c.yaruzPlatform.ReferenceSubsystem().PropertyGroup.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"6. errGet": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"6. entity2": fmt.Sprintf("%#v", e)})
	}

	list, err := c.yaruzPlatform.ReferenceSubsystem().PropertyGroup.Service.Query(cntx, &selection_condition.SelectionCondition{
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

	err = c.yaruzPlatform.ReferenceSubsystem().PropertyGroup.Service.Delete(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"7. errDelete": err.Error()})
	}

	return ctx.Write(res)
}

func (c testController) propertyViewType(ctx *routing.Context) error {
	res := make([]map[string]interface{}, 0, 10)
	res = append(res, map[string]interface{}{"test": "property-view-type"})
	cntx := ctx.Request.Context()

	entity := c.yaruzPlatform.ReferenceSubsystem().PropertyViewType.Service.NewEntity()
	entity.Sysname = "WrongName"

	err := entity.Validate()
	if err != nil {
		res = append(res, map[string]interface{}{"1. errValidate": err.Error()})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().PropertyViewType.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"2. errCreate1": err.Error()})
	}

	entity.Sysname = "property_view_type"
	err = c.yaruzPlatform.ReferenceSubsystem().PropertyViewType.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"3. errCreate": err.Error()})
	}

	e, err := c.yaruzPlatform.ReferenceSubsystem().PropertyViewType.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"4. errCreate1": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"4. entity1": fmt.Sprintf("%#v", e)})
	}

	entity.Sysname = "property_view_type_" + strconv.Itoa(int(entity.ID))
	err = c.yaruzPlatform.ReferenceSubsystem().PropertyViewType.Service.Update(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"5. errCreate": err.Error()})
	}

	e, err = c.yaruzPlatform.ReferenceSubsystem().PropertyViewType.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"6. errGet": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"6. entity2": fmt.Sprintf("%#v", e)})
	}

	list, err := c.yaruzPlatform.ReferenceSubsystem().PropertyViewType.Service.Query(cntx, &selection_condition.SelectionCondition{
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

	err = c.yaruzPlatform.ReferenceSubsystem().PropertyViewType.Service.Delete(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"7. errDelete": err.Error()})
	}

	return ctx.Write(res)
}

func (c testController) propertyType(ctx *routing.Context) error {
	res := make([]map[string]interface{}, 0, 10)
	res = append(res, map[string]interface{}{"test": "property-type"})
	cntx := ctx.Request.Context()

	entity := c.yaruzPlatform.ReferenceSubsystem().PropertyType.Service.NewEntity()
	entity.Sysname = "WrongName"

	err := entity.Validate()
	if err != nil {
		res = append(res, map[string]interface{}{"1. errValidate": err.Error()})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().PropertyType.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"2. errCreate1": err.Error()})
	}

	entity.Sysname = "property_type"
	err = c.yaruzPlatform.ReferenceSubsystem().PropertyType.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"3. errCreate": err.Error()})
	}

	e, err := c.yaruzPlatform.ReferenceSubsystem().PropertyType.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"4. errCreate1": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"4. entity1": fmt.Sprintf("%#v", e)})
	}

	entity.Sysname = "property_type_" + strconv.Itoa(int(entity.ID))
	err = c.yaruzPlatform.ReferenceSubsystem().PropertyType.Service.Update(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"5. errCreate": err.Error()})
	}

	e, err = c.yaruzPlatform.ReferenceSubsystem().PropertyType.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"6. errGet": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"6. entity2": fmt.Sprintf("%#v", e)})
	}

	list, err := c.yaruzPlatform.ReferenceSubsystem().PropertyType.Service.Query(cntx, &selection_condition.SelectionCondition{
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

	// References
	viewType := c.yaruzPlatform.ReferenceSubsystem().PropertyViewType.Service.NewEntity()
	viewType.Sysname = "view_type_" + strconv.Itoa(int(time.Now().Unix()))

	err = c.yaruzPlatform.ReferenceSubsystem().PropertyViewType.Service.Create(cntx, viewType)
	if err != nil {
		res = append(res, map[string]interface{}{"7. errCreate ViewType": err.Error()})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().PropertyType.Service.BindPropertyViewType(cntx, entity, viewType.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"8. BindPropertyViewType err": err.Error()})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().PropertyType.Service.InitPropertyViewTypes(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"9. GetPropertyViewTypes err": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"9. GetPropertyViewTypes list": fmt.Sprintf("%#v", entity.PropertyViewTypes)})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().PropertyType.Service.UnbindPropertyViewType(cntx, entity, viewType.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"10. UnbindPropertyViewType err": err.Error()})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().PropertyType.Service.Delete(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"11. errDelete": err.Error()})
	}

	return ctx.Write(res)
}

func (c testController) property(ctx *routing.Context) error {
	res := make([]map[string]interface{}, 0, 10)
	res = append(res, map[string]interface{}{"test": "property"})
	cntx := ctx.Request.Context()

	entity := c.yaruzPlatform.ReferenceSubsystem().Property.Service.NewEntity()
	entity.Sysname = "WrongName"
	entity.PropertyTypeID = 3
	entity.PropertyUnitID = 1

	err := entity.Validate()
	if err != nil {
		res = append(res, map[string]interface{}{"1. errValidate": err.Error()})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().Property.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"2. errCreate1": err.Error()})
	}

	entity.Sysname = "property"
	err = c.yaruzPlatform.ReferenceSubsystem().Property.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"3. errCreate": err.Error()})
	}

	e, err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"4. errCreate1": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"4. entity1": fmt.Sprintf("%#v", e)})
	}

	entity.Sysname = "property_" + strconv.Itoa(int(entity.ID))
	err = c.yaruzPlatform.ReferenceSubsystem().Property.Service.Update(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"5. errCreate": err.Error()})
	}

	e, err = c.yaruzPlatform.ReferenceSubsystem().Property.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"6. errGet": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"6. entity2": fmt.Sprintf("%#v", e)})
	}

	list, err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.Query(cntx, &selection_condition.SelectionCondition{
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

	err = c.yaruzPlatform.ReferenceSubsystem().Property.Service.Delete(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"7. errDelete": err.Error()})
	}

	return ctx.Write(res)
}

func (c testController) entityType(ctx *routing.Context) error {
	res := make([]map[string]interface{}, 0, 10)
	res = append(res, map[string]interface{}{"test": "entity-type"})
	cntx := ctx.Request.Context()

	entity := c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.NewEntity()
	entity.Sysname = "WrongName"

	err := entity.Validate()
	if err != nil {
		res = append(res, map[string]interface{}{"1. errValidate": err.Error()})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"2. errCreate1": err.Error()})
	}

	entity.Sysname = "entity_type"
	err = c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"3. errCreate": err.Error()})
	}

	e, err := c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"4. errCreate1": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"4. entity1": fmt.Sprintf("%#v", e)})
	}

	entity.Sysname = "entity_type_" + strconv.Itoa(int(entity.ID))
	err = c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.Update(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"5. errCreate": err.Error()})
	}

	e, err = c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"6. errGet": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"6. entity2": fmt.Sprintf("%#v", e)})
	}

	list, err := c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.Query(cntx, &selection_condition.SelectionCondition{
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

	err = c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.Delete(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"7. errDelete": err.Error()})
	}

	return ctx.Write(res)
}
