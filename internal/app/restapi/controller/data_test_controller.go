package controller

import (
	"context"
	"time"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_unit"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/minipkg/log"
	"github.com/yaruz/app/pkg/yarus_platform"
)

type dataTestController struct {
	Logger        log.ILogger
	yaruzPlatform yarus_platform.IPlatform
}

var langEngID = uint(1)
var langRusID = uint(2)
var propertyUnitMMSysname = "mm"
var entityTypeObjSysname = "obj"
var propertyNumberSysname = "num"
var propertyLenSysname = "len"
var propertyBoolSysname = "bool"
var propertyDateSysname = "date"
var propertyTSSysname = "time"
var propertyNameSysname = "name"
var propertyOptSysname = "opt"
var propertyOpt0ID = 0
var propertyOpt0Key = "one"
var propertyOpt0Val = 1
var propertyOpt1ID = 1
var propertyOpt1Key = "two"
var propertyOpt1Val = 2
var propertyOpt2ID = 2
var propertyOpt2Key = "three"
var propertyOpt2Val = 3

// RegisterHandlers sets up the routing of the HTTP handlers.
//	GET /api/models/ - список всех моделей
//	GET /api/model/{ID} - детали модели
func RegisterDataTestHandlers(r *routing.RouteGroup, yaruzPlatform yarus_platform.IPlatform, logger log.ILogger, authHandler routing.Handler) {
	c := dataTestController{
		Logger:        logger,
		yaruzPlatform: yaruzPlatform,
	}

	r.Get("/entity", c.entity)
	r.Get("/entity-text", c.entityText)
}

func (c dataTestController) entity(ctx *routing.Context) error {
	res := make([]map[string]interface{}, 0, 10)
	res = append(res, map[string]interface{}{"test": "entity"})
	cntx := ctx.Request.Context()

	if err := c.propertyUnitsInit(cntx); err != nil {
		res = append(res, map[string]interface{}{"propertyUnitsInit": err.Error()})
	}

	if err := c.propertiesInit(cntx); err != nil {
		res = append(res, map[string]interface{}{"propertiesInit": err.Error()})
	}

	if err := c.entityTypesInit(cntx); err != nil {
		res = append(res, map[string]interface{}{"entityTypesInit": err.Error()})
	}

	propertyNumber, err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.TFirst(cntx, &property.Property{Sysname: propertyNumberSysname}, langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"finding for propertyNumber": err.Error()})
	}

	propertyLen, err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.TFirst(cntx, &property.Property{Sysname: propertyLenSysname}, langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"finding for propertyLen": err.Error()})
	}

	propertyOpt, err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.TFirst(cntx, &property.Property{Sysname: propertyOptSysname}, langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"finding for propertyOpt": err.Error()})
	}

	propertyBool, err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.TFirst(cntx, &property.Property{Sysname: propertyBoolSysname}, langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"finding for propertyBool": err.Error()})
	}

	propertyDate, err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.TFirst(cntx, &property.Property{Sysname: propertyDateSysname}, langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"finding for propertyDate": err.Error()})
	}

	propertyTime, err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.TFirst(cntx, &property.Property{Sysname: propertyTSSysname}, langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"finding for propertyTS": err.Error()})
	}

	entity := c.yaruzPlatform.DataSubsystem().Entity.Service.NewEntity()
	//entity.PropertiesValuesMap = map[uint]interface{}{
	//	propertyNumber.ID: int(158),
	//	propertyLen.ID:    float64(32.543),
	//	propertyOpt.ID:    propertyOpt1Val,
	//	propertyBool.ID:   true,
	//	propertyDate.ID:   "2021-09-22",
	//	propertyTime.ID:     time.Now().Unix(),
	//}

	err = c.yaruzPlatform.DataSubsystem().Entity.Service.EntitySetPropertyValue(cntx, entity, propertyNumber, int(158), langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"Set propertyNumber": err.Error()})
	}

	err = c.yaruzPlatform.DataSubsystem().Entity.Service.EntitySetPropertyValue(cntx, entity, propertyLen, float64(32.543), langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"Set propertyLen": err.Error()})
	}

	err = c.yaruzPlatform.DataSubsystem().Entity.Service.EntitySetPropertyValue(cntx, entity, propertyOpt, propertyOpt1Val, langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"Set propertyOpt": err.Error()})
	}

	err = c.yaruzPlatform.DataSubsystem().Entity.Service.EntitySetPropertyValue(cntx, entity, propertyBool, true, langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"Set propertyBool": err.Error()})
	}

	err = c.yaruzPlatform.DataSubsystem().Entity.Service.EntitySetPropertyValue(cntx, entity, propertyDate, time.Now(), langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"Set propertyDate": err.Error()})
	}

	err = c.yaruzPlatform.DataSubsystem().Entity.Service.EntitySetPropertyValue(cntx, entity, propertyTime, time.Now(), langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"Set propertyTime": err.Error()})
	}

	err = c.yaruzPlatform.DataSubsystem().Entity.Service.Create(cntx, entity, langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"entity.Create()": err.Error()})
	}

	entityG, err := c.yaruzPlatform.DataSubsystem().Entity.Service.Get(cntx, entity.ID, langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"entity.Get()": err.Error()})
	}
	res = append(res, map[string]interface{}{"entity": entityG})

	return ctx.Write(res)
}

func (c dataTestController) entityText(ctx *routing.Context) error {
	res := make([]map[string]interface{}, 0, 10)
	res = append(res, map[string]interface{}{"test": "entity"})
	cntx := ctx.Request.Context()

	if err := c.propertyUnitsInit(cntx); err != nil {
		res = append(res, map[string]interface{}{"propertyUnitsInit": err.Error()})
	}

	if err := c.propertiesInit(cntx); err != nil {
		res = append(res, map[string]interface{}{"propertiesInit": err.Error()})
	}

	if err := c.entityTypesInit(cntx); err != nil {
		res = append(res, map[string]interface{}{"entityTypesInit": err.Error()})
	}

	propertyName, err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.TFirst(cntx, &property.Property{Sysname: propertyNameSysname}, langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"finding for propertyName": err.Error()})
	}

	entity := c.yaruzPlatform.DataSubsystem().Entity.Service.NewEntity()

	err = c.yaruzPlatform.DataSubsystem().Entity.Service.EntitySetPropertyValue(cntx, entity, propertyName, "Андрей", langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"Set propertyBool": err.Error()})
	}

	err = c.yaruzPlatform.DataSubsystem().Entity.Service.Create(cntx, entity, langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"entity.Create()": err.Error()})
	}

	err = c.yaruzPlatform.DataSubsystem().Entity.Service.EntitySetPropertyValue(cntx, entity, propertyName, "Андрей Калинин", langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"Set propertyBool": err.Error()})
	}

	err = c.yaruzPlatform.DataSubsystem().Entity.Service.Update(cntx, entity, langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"entity.Update()": err.Error()})
	}

	entity1, err := c.yaruzPlatform.DataSubsystem().Entity.Service.Get(cntx, entity.ID, langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"entity1.Get()": err.Error()})
	}
	res = append(res, map[string]interface{}{"entity1": entity1})

	err = c.yaruzPlatform.DataSubsystem().Entity.Service.EntitySetPropertyValue(cntx, entity1, propertyName, "Андрей Калинин 1", langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"Set propertyBool": err.Error()})
	}

	err = c.yaruzPlatform.DataSubsystem().Entity.Service.Update(cntx, entity, langRusID)
	if err != nil {
		res = append(res, map[string]interface{}{"entity.Update()": err.Error()})
	}

	entity2, err := c.yaruzPlatform.DataSubsystem().Entity.Service.Get(cntx, entity.ID, langEngID)
	if err != nil {
		res = append(res, map[string]interface{}{"entity2.Get()": err.Error()})
	}
	res = append(res, map[string]interface{}{"entity2": entity2})

	err = c.yaruzPlatform.DataSubsystem().Entity.Service.EntitySetPropertyValue(cntx, entity2, propertyName, "Andrey", langEngID)
	if err != nil {
		res = append(res, map[string]interface{}{"Set propertyBool": err.Error()})
	}

	err = c.yaruzPlatform.DataSubsystem().Entity.Service.EntitySetPropertyValue(cntx, entity2, propertyName, "Andrey Kalinin", langEngID)
	if err != nil {
		res = append(res, map[string]interface{}{"Set propertyBool": err.Error()})
	}

	err = c.yaruzPlatform.DataSubsystem().Entity.Service.Update(cntx, entity2, langEngID)
	if err != nil {
		res = append(res, map[string]interface{}{"entity2.Update()": err.Error()})
	}

	entity3, err := c.yaruzPlatform.DataSubsystem().Entity.Service.Get(cntx, entity.ID, langEngID)
	if err != nil {
		res = append(res, map[string]interface{}{"entity3.Get()": err.Error()})
	}
	res = append(res, map[string]interface{}{"entity3": entity3})

	return ctx.Write(res)
}

func (c dataTestController) propertyUnitsInit(ctx context.Context) error {

	propertyUnitMM := c.yaruzPlatform.ReferenceSubsystem().PropertyUnit.Service.NewEntity()
	propertyUnitMM.Sysname = propertyUnitMMSysname
	propertyUnitLenName := "мм"
	propertyUnitLenDesc := "миллиметры"
	propertyUnitMM.Name = &propertyUnitLenName
	propertyUnitMM.Description = &propertyUnitLenDesc

	err := c.yaruzPlatform.ReferenceSubsystem().PropertyUnit.Service.TCreate(ctx, propertyUnitMM, langRusID)
	if err != nil {
		return err
	}
	return nil
}

func (c dataTestController) propertiesInit(ctx context.Context) error {

	propertyUnitItem, err := c.yaruzPlatform.ReferenceSubsystem().PropertyUnit.Service.TFirst(ctx, &property_unit.PropertyUnit{Sysname: property_unit.SysnameItem}, langRusID)
	if err != nil {
		return err
	}

	propertyNumber := c.yaruzPlatform.ReferenceSubsystem().Property.Service.NewEntity()
	propertyNumber.Sysname = propertyNumberSysname
	propertyNumber.PropertyTypeID = property_type.IDInt
	propertyNumber.PropertyUnitID = &propertyUnitItem.ID
	propertyNumberName := "количество"
	propertyNumberDesc := "количество (шт)"
	propertyNumber.Name = &propertyNumberName
	propertyNumber.Description = &propertyNumberDesc

	if err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.TCreate(ctx, propertyNumber, langRusID); err != nil {
		return err
	}

	propertyUnitMM, err := c.yaruzPlatform.ReferenceSubsystem().PropertyUnit.Service.TFirst(ctx, &property_unit.PropertyUnit{Sysname: "mm"}, langRusID)
	if err != nil {
		return err
	}

	propertyLen := c.yaruzPlatform.ReferenceSubsystem().Property.Service.NewEntity()
	propertyLen.Sysname = propertyLenSysname
	propertyLen.PropertyTypeID = property_type.IDFloat
	propertyLen.PropertyUnitID = &propertyUnitMM.ID
	propertyLenName := "длина"
	propertyLenDesc := "длина (мм)"
	propertyLen.Name = &propertyLenName
	propertyLen.Description = &propertyLenDesc

	if err = c.yaruzPlatform.ReferenceSubsystem().Property.Service.TCreate(ctx, propertyLen, langRusID); err != nil {
		return err
	}

	propertyBool := c.yaruzPlatform.ReferenceSubsystem().Property.Service.NewEntity()
	propertyBool.Sysname = propertyBoolSysname
	propertyBool.PropertyTypeID = property_type.IDBoolean
	propertyBool.PropertyUnitID = &propertyUnitMM.ID
	propertyBoolName := "логическое"
	propertyBoolDesc := "логическое свойство (true/false)"
	propertyBool.Name = &propertyBoolName
	propertyBool.Description = &propertyBoolDesc

	if err = c.yaruzPlatform.ReferenceSubsystem().Property.Service.TCreate(ctx, propertyBool, langRusID); err != nil {
		return err
	}

	propertyDate := c.yaruzPlatform.ReferenceSubsystem().Property.Service.NewEntity()
	propertyDate.Sysname = propertyDateSysname
	propertyDate.PropertyTypeID = property_type.IDDate
	propertyDate.PropertyUnitID = &propertyUnitMM.ID
	propertyDateName := "дата"
	propertyDateDesc := "дата события"
	propertyDate.Name = &propertyDateName
	propertyDate.Description = &propertyDateDesc

	if err = c.yaruzPlatform.ReferenceSubsystem().Property.Service.TCreate(ctx, propertyDate, langRusID); err != nil {
		return err
	}

	propertyTS := c.yaruzPlatform.ReferenceSubsystem().Property.Service.NewEntity()
	propertyTS.Sysname = propertyTSSysname
	propertyTS.PropertyTypeID = property_type.IDTime
	propertyTS.PropertyUnitID = &propertyUnitMM.ID
	propertyTSName := "дата"
	propertyTSDesc := "дата события"
	propertyTS.Name = &propertyTSName
	propertyTS.Description = &propertyTSDesc

	if err = c.yaruzPlatform.ReferenceSubsystem().Property.Service.TCreate(ctx, propertyTS, langRusID); err != nil {
		return err
	}

	propertyName := c.yaruzPlatform.ReferenceSubsystem().Property.Service.NewEntity()
	propertyName.Sysname = propertyNameSysname
	propertyName.PropertyTypeID = property_type.IDText
	propertyName.PropertyUnitID = &propertyUnitMM.ID
	propertyNameName := "имя"
	propertyNameDesc := "ну, чьё-то имя, например"
	propertyName.Name = &propertyNameName
	propertyName.Description = &propertyNameDesc

	if err = c.yaruzPlatform.ReferenceSubsystem().Property.Service.TCreate(ctx, propertyName, langRusID); err != nil {
		return err
	}

	propertyOpts := c.yaruzPlatform.ReferenceSubsystem().Property.Service.NewEntity()
	propertyOpts.Sysname = propertyOptSysname
	propertyOpts.PropertyTypeID = property_type.IDInt
	propertyOpts.PropertyUnitID = &propertyUnitItem.ID
	propertyOpts.Options = []map[string]interface{}{
		{propertyOpt0Key: propertyOpt0Val},
		{propertyOpt1Key: propertyOpt1Val},
		{propertyOpt2Key: propertyOpt2Val},
	}
	propertyOptsName := "опции"
	propertyOptsDesc := "опции (шт)"
	propertyOpts.Name = &propertyOptsName
	propertyOpts.Description = &propertyOptsDesc

	if err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.TCreate(ctx, propertyOpts, langRusID); err != nil {
		return err
	}

	return nil
}

func (c dataTestController) entityTypesInit(ctx context.Context) error {
	entityTypeObj := c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.NewEntity()
	entityTypeObj.Sysname = entityTypeObjSysname
	entityTypeObjName := "объект"
	entityTypeObjDesc := "тип объект"
	entityTypeObj.Name = &entityTypeObjName
	entityTypeObj.Description = &entityTypeObjDesc

	if err := c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.TCreate(ctx, entityTypeObj, langRusID); err != nil {
		return err
	}

	return nil
}
