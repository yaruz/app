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

var langID = uint(1)
var propertyUnitMMSysname = "mm"
var entityTypeObjSysname = "obj"
var propertyNumberSysname = "num"
var propertyLenSysname = "len"
var propertyBoolSysname = "bool"
var propertyDateSysname = "date"
var propertyTSSysname = "ts"
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

	propertyNumber, err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.TFirst(cntx, &property.Property{Sysname: propertyNumberSysname}, langID)
	if err != nil {
		res = append(res, map[string]interface{}{"finding for propertyNumber": err.Error()})
	}

	propertyLen, err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.TFirst(cntx, &property.Property{Sysname: propertyLenSysname}, langID)
	if err != nil {
		res = append(res, map[string]interface{}{"finding for propertyLen": err.Error()})
	}

	propertyOpt, err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.TFirst(cntx, &property.Property{Sysname: propertyOptSysname}, langID)
	if err != nil {
		res = append(res, map[string]interface{}{"finding for propertyOpt": err.Error()})
	}

	propertyBool, err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.TFirst(cntx, &property.Property{Sysname: propertyBoolSysname}, langID)
	if err != nil {
		res = append(res, map[string]interface{}{"finding for propertyBool": err.Error()})
	}

	propertyDate, err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.TFirst(cntx, &property.Property{Sysname: propertyDateSysname}, langID)
	if err != nil {
		res = append(res, map[string]interface{}{"finding for propertyDate": err.Error()})
	}

	propertyTS, err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.TFirst(cntx, &property.Property{Sysname: propertyTSSysname}, langID)
	if err != nil {
		res = append(res, map[string]interface{}{"finding for propertyTS": err.Error()})
	}

	//propertyName, err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.TFirst(cntx, &property.Property{Sysname: propertyNameSysname}, langID)
	//if err != nil {
	//	res = append(res, map[string]interface{}{"finding for propertyName": err.Error()})
	//}

	entity := c.yaruzPlatform.DataSubsystem().Entity.Service.NewEntity()
	entity.PropertiesValuesMap = map[uint]interface{}{
		propertyNumber.ID: int(158),
		propertyLen.ID:    float64(32.543),
		propertyOpt.ID:    propertyOpt1Val,
		propertyBool.ID:   true,
		propertyDate.ID:   "2021-09-22",
		propertyTS.ID:     time.Now().Unix(),
	}

	err = c.yaruzPlatform.DataSubsystem().Entity.Service.EntityInit(cntx, entity, langID)
	if err != nil {
		res = append(res, map[string]interface{}{"entity.AfterFind()": err.Error()})
	}

	err = c.yaruzPlatform.DataSubsystem().Entity.Service.Create(cntx, entity)
	if err != nil {
		res = append(res, map[string]interface{}{"entity.Create()": err.Error()})
	}

	entityG, err := c.yaruzPlatform.DataSubsystem().Entity.Service.Get(cntx, entity.ID, langID)
	if err != nil {
		res = append(res, map[string]interface{}{"entity.Get()": err.Error()})
	}
	res = append(res, map[string]interface{}{"entity": entityG})

	return ctx.Write(res)
}

func (c dataTestController) propertyUnitsInit(ctx context.Context) error {

	propertyUnitMM := c.yaruzPlatform.ReferenceSubsystem().PropertyUnit.Service.NewEntity()
	propertyUnitMM.Sysname = propertyUnitMMSysname
	propertyUnitLenName := "мм"
	propertyUnitLenDesc := "миллиметры"
	propertyUnitMM.Name = &propertyUnitLenName
	propertyUnitMM.Description = &propertyUnitLenDesc

	err := c.yaruzPlatform.ReferenceSubsystem().PropertyUnit.Service.TCreate(ctx, propertyUnitMM, langID)
	if err != nil {
		return err
	}
	return nil
}

func (c dataTestController) propertiesInit(ctx context.Context) error {

	propertyUnitItem, err := c.yaruzPlatform.ReferenceSubsystem().PropertyUnit.Service.TFirst(ctx, &property_unit.PropertyUnit{Sysname: property_unit.SysnameItem}, langID)
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

	if err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.TCreate(ctx, propertyNumber, langID); err != nil {
		return err
	}

	propertyUnitMM, err := c.yaruzPlatform.ReferenceSubsystem().PropertyUnit.Service.TFirst(ctx, &property_unit.PropertyUnit{Sysname: "mm"}, langID)
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

	if err = c.yaruzPlatform.ReferenceSubsystem().Property.Service.TCreate(ctx, propertyLen, langID); err != nil {
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

	if err = c.yaruzPlatform.ReferenceSubsystem().Property.Service.TCreate(ctx, propertyBool, langID); err != nil {
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

	if err = c.yaruzPlatform.ReferenceSubsystem().Property.Service.TCreate(ctx, propertyDate, langID); err != nil {
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

	if err = c.yaruzPlatform.ReferenceSubsystem().Property.Service.TCreate(ctx, propertyTS, langID); err != nil {
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

	if err = c.yaruzPlatform.ReferenceSubsystem().Property.Service.TCreate(ctx, propertyName, langID); err != nil {
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

	if err := c.yaruzPlatform.ReferenceSubsystem().Property.Service.TCreate(ctx, propertyOpts, langID); err != nil {
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

	if err := c.yaruzPlatform.ReferenceSubsystem().EntityType.Service.TCreate(ctx, entityTypeObj, langID); err != nil {
		return err
	}

	return nil
}
