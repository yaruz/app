package controller

import (
	"fmt"
	"strconv"
	"time"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/minipkg/log"
	"github.com/minipkg/selection_condition"
	"github.com/yaruz/app/pkg/yarus_platform"
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

	r.Get("/text-source", c.textSource)
	r.Get("/text-value", c.textValue)
	r.Get("/property-unit", c.propertyUnit)
	r.Get("/property-group", c.propertyGroup)
	r.Get("/property-view-type", c.propertyViewType)
	r.Get("/property-type", c.propertyType)
	r.Get("/property-type_with-text", c.propertyTypeWithText)
	r.Get("/property", c.property)
	r.Get("/relation", c.relation)
	r.Get("/property-options-validation", c.propertyOptionsValidation)
	r.Get("/entity-type", c.entityType)
	r.Get("/entity-type-binding", c.entityTypeBinding)
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
	entity.TextLangID = 1
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

	err = c.yaruzPlatform.ReferenceSubsystem().PropertyType.Service.Delete(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"11. errDelete": err.Error()})
	}

	return ctx.Write(res)
}

func (c testController) propertyTypeWithText(ctx *routing.Context) error {
	res := make([]map[string]interface{}, 0, 10)
	res = append(res, map[string]interface{}{"test": "property-type-with-text"})
	cntx := ctx.Request.Context()
	langID := uint(1)

	entity := c.yaruzPlatform.ReferenceSubsystem().PropertyType.Service.NewEntity()
	entity.Sysname = "property_type_" + strconv.Itoa(int(entity.ID))
	name := "name1_" + strconv.Itoa(int(entity.ID))
	entity.Name = &name

	err := c.yaruzPlatform.ReferenceSubsystem().PropertyType.Service.TCreate(cntx, entity, langID)
	if err != nil {
		res = append(res, map[string]interface{}{"2. errCreate1": err.Error()})
	}

	e, err := c.yaruzPlatform.ReferenceSubsystem().PropertyType.Service.TGet(cntx, entity.ID, langID)
	if err != nil {
		res = append(res, map[string]interface{}{"4. errCreate1": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"4. entity1": fmt.Sprintf("%#v", e)})
	}

	desc := "desc_" + strconv.Itoa(int(entity.ID))
	entity.Description = &desc
	*entity.Name = "name_" + strconv.Itoa(int(entity.ID))
	err = c.yaruzPlatform.ReferenceSubsystem().PropertyType.Service.TUpdate(cntx, entity, langID)
	if err != nil {
		res = append(res, map[string]interface{}{"5. errCreate": err.Error()})
	}

	e, err = c.yaruzPlatform.ReferenceSubsystem().PropertyType.Service.TGet(cntx, entity.ID, langID)
	if err != nil {
		res = append(res, map[string]interface{}{"6. errGet": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"6. entity2": fmt.Sprintf("%#v", e)})
	}

	list, err := c.yaruzPlatform.ReferenceSubsystem().PropertyType.Service.TQuery(cntx, &selection_condition.SelectionCondition{
		Where: selection_condition.WhereCondition{
			Field:     "Sysname",
			Condition: "eq",
			Value:     entity.Sysname,
		},
	}, langID)
	if err != nil {
		res = append(res, map[string]interface{}{"6. errGet": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"6. list": fmt.Sprintf("%#v", list)})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().PropertyType.Service.Delete(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"11. errDelete": err.Error()})
	}

	return ctx.Write(res)
}

func (c testController) property(ctx *routing.Context) error {
	res := make([]map[string]interface{}, 0, 10)
	res = append(res, map[string]interface{}{"test": "property"})
	cntx := ctx.Request.Context()
	// ViewType
	viewType := c.yaruzPlatform.ReferenceSubsystem().PropertyViewType.Service.NewEntity()
	viewType.Sysname = "view_type_" + strconv.Itoa(int(time.Now().Unix()))

	err := c.yaruzPlatform.ReferenceSubsystem().PropertyViewType.Service.Create(cntx, viewType)
	if err != nil {
		res = append(res, map[string]interface{}{"0. Create ViewType err": err.Error()})
	}
	// Group
	group := c.yaruzPlatform.ReferenceSubsystem().PropertyGroup.Service.NewEntity()
	group.Sysname = "group_" + strconv.Itoa(int(time.Now().Unix()))

	err = c.yaruzPlatform.ReferenceSubsystem().PropertyGroup.Service.Create(cntx, group)
	if err != nil {
		res = append(res, map[string]interface{}{"0. Create Group err": err.Error()})
	}
	var one uint = 1
	entity := c.yaruzPlatform.ReferenceSubsystem().Property.Service.NewEntity()
	entity.Sysname = "WrongName"
	entity.PropertyTypeID = property_type.IDInt
	entity.PropertyUnitID = &one
	entity.PropertyViewTypeID = &one
	entity.PropertyGroupID = &one

	err = entity.Validate()
	if err != nil {
		res = append(res, map[string]interface{}{"1. errValidate": err.Error()})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().Property.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"2. errCreate1": err.Error()})
	}

	entity.Sysname = "property" + strconv.Itoa(int(time.Now().Unix()))
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

type PropertyOptionsTestCases [][][]map[string]interface{}

func (c testController) propertyOptionsValidation(ctx *routing.Context) error {
	res := make([]map[string]interface{}, 0, 10)
	res = append(res, map[string]interface{}{"test": "property option validation"})

	entity := c.yaruzPlatform.ReferenceSubsystem().Property.Service.NewEntity()
	entity.Sysname = "options_validation_test_" + strconv.Itoa(int(time.Now().Unix()))
	var one uint = 1
	entity.PropertyUnitID = &one

	var typeTestCases = [][][]map[string]interface{}{
		{ // boolean
			{
				{"1": 1},
				{"2": 2},
			},
			{
				{"1": "one"},
				{"2": "two"},
			},
			{
				{"1": true},
				{"2": false},
			},
		},
		{ // int
			{
				{"1": "one"},
				{"2": "two"},
			},
			{
				{"1": true},
				{"2": false},
			},
			{
				{"1": 1},
				{"2": 2},
			},
		},
		{ // bigint
			{
				{"1": "one"},
				{"2": "two"},
			},
			{
				{"1": true},
				{"2": false},
			},
			{
				{"1": 1},
				{"2": 2},
			},
		},
		{ // float
			{
				{"1": "one"},
				{"2": "two"},
			},
			{
				{"1": true},
				{"2": false},
			},
			{
				{"1": 1},
				{"2": 2},
			},
			{
				{"1": 1.1},
				{"2": 2.1},
			},
		},
		{ // date
			{
				{"1": 1},
				{"2": 2},
			},
			{
				{"1": "one"},
				{"2": "two"},
			},
			{
				{"1": "2021-05-27"},
				{"2": "2021-05-30"},
			},
		},
		{ // timestamp
			{
				{"1": "one"},
				{"2": "two"},
			},
			{
				{"1": true},
				{"2": false},
			},
			{
				{"1": 1405544146},
				{"2": 1405545146},
			},
		},
		{ // text
			{
				{"1": 1},
				{"2": 2},
			},
			{
				{"1": true},
				{"2": false},
			},
			{
				{"1": "one"},
				{"2": "two"},
			},
		},
	}

	for propertyType, testCases := range typeTestCases {
		propertyType++
		entity.PropertyTypeID = uint(propertyType)

		for caseNum, options := range testCases {
			entity.Options = options

			err := entity.Validate()
			if err != nil {
				res = append(res, map[string]interface{}{"PropertyType:" + strconv.Itoa(propertyType) + " test case #" + strconv.Itoa(caseNum) + ". error": err.Error()})
			} else {
				res = append(res, map[string]interface{}{"PropertyType:" + strconv.Itoa(propertyType) + " test case #" + strconv.Itoa(caseNum) + ".": "success"})
			}
		}
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
		res = append(res, map[string]interface{}{"1. EntityType Validation error: ": err.Error()})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"2. EntityType Creation error: ": err.Error()})
	}

	entity.Sysname = "entity_type"
	err = c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"3. EntityType Creation error: ": err.Error()})
	}

	e, err := c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"3. EntityType Getting error: ": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"4. EntityType: ": fmt.Sprintf("%#v", e)})
	}

	entity.Sysname = "entity_type_" + strconv.Itoa(int(entity.ID))
	err = c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.Update(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"5. EntityType Updating error: ": err.Error()})
	}

	e, err = c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"6. EntityType Getting error: ": err.Error()})
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
		res = append(res, map[string]interface{}{"7. EntityType Getting a list error: ": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"7. list": fmt.Sprintf("%#v", list)})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.Delete(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"8. EntityType Deleting error: ": err.Error()})
	}

	return ctx.Write(res)
}

func (c testController) entityTypeBinding(ctx *routing.Context) error {
	res := make([]map[string]interface{}, 0, 10)
	res = append(res, map[string]interface{}{"test": "entity-type"})
	cntx := ctx.Request.Context()

	var one uint = 1
	property1 := c.yaruzPlatform.ReferenceSubsystem().Property.Service.NewEntity()
	property1.Sysname = "property_" + strconv.Itoa(int(time.Now().Unix()))
	property1.PropertyTypeID = property_type.IDInt
	property1.PropertyUnitID = &one
	property2 := c.yaruzPlatform.ReferenceSubsystem().Property.Service.NewEntity()
	property2.Sysname = "property_" + strconv.Itoa(int(time.Now().Unix()))
	property2.PropertyTypeID = property_type.IDFloat
	property2.PropertyUnitID = &one
	err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.Create(cntx, property1)
	if err != nil {
		res = append(res, map[string]interface{}{"0. Property1 Create error: ": err.Error()})
	}
	err = c.yaruzPlatform.ReferenceSubsystem().Property.Service.Create(cntx, property2)
	if err != nil {
		res = append(res, map[string]interface{}{"0. Property2 Create error: ": err.Error()})
	}

	entity := c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.NewEntity()
	entity.Sysname = "entity_type_" + strconv.Itoa(int(time.Now().Unix()))
	err = c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"1. EntityType Creation error: ": err.Error()})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.BindProperty(cntx, entity.ID, property1.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"2. EntityType Updating error: ": err.Error()})
	}
	err = c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.BindProperty(cntx, entity.ID, property2.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"2. EntityType Updating error: ": err.Error()})
	}

	e, err := c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"3. EntityType Getting error: ": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"3. EntityType: ": fmt.Sprintf("%#v", e)})
	}

	props, rels, err := c.yaruzPlatform.ReferenceSubsystem().Relation.Service.GetPropertiesAndRelationsByEntityTypeID(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"4. EntityType Getting error: ": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"4. EntityType props: ": props})
		res = append(res, map[string]interface{}{"4. EntityType rels: ": rels})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.UnbindProperty(cntx, entity.ID, property1.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"5. EntityType Unbind error: ": err.Error()})
	}

	err = c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.UnbindAllProperty(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"6. EntityType UnbindAll error: ": err.Error()})
	}

	e, err = c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"7. EntityType Getting error: ": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"7. entity2": fmt.Sprintf("%#v", e)})
	}

	props, rels, err = c.yaruzPlatform.ReferenceSubsystem().Relation.Service.GetPropertiesAndRelationsByEntityTypeID(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"4. EntityType Getting error: ": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"4. EntityType props: ": props})
		res = append(res, map[string]interface{}{"4. EntityType rels: ": rels})
	}

	return ctx.Write(res)
}

func (c testController) relation(ctx *routing.Context) error {
	res := make([]map[string]interface{}, 0, 10)
	res = append(res, map[string]interface{}{"test": "entity-type"})
	cntx := ctx.Request.Context()

	entityType1 := c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.NewEntity()
	entityType1.Sysname = "entity_type1_" + strconv.Itoa(int(time.Now().Unix()))

	err := c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.Create(cntx, entityType1)
	if err != nil {
		res = append(res, map[string]interface{}{"1. Create EntityType1 error": err.Error()})
	}

	entityType2 := c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.NewEntity()
	entityType2.Sysname = "entity_type2_" + strconv.Itoa(int(time.Now().Unix()))

	err = c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.Create(cntx, entityType2)
	if err != nil {
		res = append(res, map[string]interface{}{"2. Create EntityType2 error": err.Error()})
	}

	entity := c.yaruzPlatform.ReferenceSubsystem().Relation.Service.NewEntity()
	entity.Sysname = "relation_" + strconv.Itoa(int(time.Now().Unix()))
	entity.PropertyTypeID = property_type.IDRelation
	entity.UndependedEntityType = entityType2
	entity.DependedEntityType = entityType1

	err = c.yaruzPlatform.ReferenceSubsystem().Relation.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"3. Create Relation error": err.Error()})
	}

	e, err := c.yaruzPlatform.ReferenceSubsystem().Relation.Service.Get(cntx, entity.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"4. Get new Relation error": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"4. New Relation": fmt.Sprintf("%#v", e)})
	}

	props, rels, err := c.yaruzPlatform.ReferenceSubsystem().Relation.Service.GetPropertiesAndRelationsByEntityTypeID(cntx, entityType1.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"5. EntityType Getting error: ": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"5. EntityType props: ": props})
		res = append(res, map[string]interface{}{"5. EntityType rels: ": rels})
	}

	entityType, err := c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.Get(cntx, entityType1.ID)
	if err != nil {
		res = append(res, map[string]interface{}{"6. EntityType Getting error: ": err.Error()})
	} else {
		res = append(res, map[string]interface{}{"6. entityType": entityType})
	}

	return ctx.Write(res)
}
