package property

import (
	"context"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_group"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_unit"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_view_type"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/yaruz/app/pkg/yarus_platform/config"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_lang"

	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"

	"github.com/minipkg/selection_condition"

	"github.com/pkg/errors"

	"github.com/minipkg/log"
)

// IService encapsulates usecase logic.
type IService interface {
	NewEntity() *Property
	PropertyInit(ctx context.Context, PropertiesConfig config.Properties, entityTypeSysname string) (propertyIDs []uint, err error)
	Get(ctx context.Context, id uint) (*Property, error)
	First(ctx context.Context, entity *Property) (*Property, error)
	Query(ctx context.Context, query *selection_condition.SelectionCondition) ([]Property, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *Property) error
	Update(ctx context.Context, entity *Property) error
	Delete(ctx context.Context, entity *Property) error
	TGet(ctx context.Context, id uint, langID uint) (*Property, error)
	TFirst(ctx context.Context, entity *Property, langID uint) (*Property, error)
	TQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]Property, error)
	TCreate(ctx context.Context, entity *Property, langID uint) (err error)
	TUpdate(ctx context.Context, entity *Property, langID uint) (err error)
	GetSysnames(ctx context.Context) ([]string, error)
	GetSysnamesEmptyInterfaceSlice(ctx context.Context) ([]interface{}, error)
	GetMapSysnameID(ctx context.Context) (map[string]uint, error)
	GetMapIDSysname(ctx context.Context) (map[uint]string, error)
	GetMapSysnameTypeID(ctx context.Context) (map[string]uint, error)
	GetMapIDTypeID(ctx context.Context) (map[uint]uint, error)
	GetIDBySysname(ctx context.Context, sysname string) (uint, error)
	GetSysnameByID(ctx context.Context, id uint) (string, error)
	GetTypeIDBySysname(ctx context.Context, sysname string) (uint, error)
	GetTypeIDByID(ctx context.Context, id uint) (uint, error)
	GetBySysname(ctx context.Context, sysname string, langID uint) (*Property, error)
}

type service struct {
	logger                  log.ILogger
	repository              Repository
	propertyTypeService     property_type.IService
	propertyUnitService     property_unit.IService
	propertyViewTypeService property_view_type.IService
	propertyGroupService    property_group.IService
	langFinder              text_lang.LangFinder
}

var _ IService = (*service)(nil)

// NewService creates a new service.
func NewService(logger log.ILogger, repo Repository, propertyTypeService property_type.IService, propertyUnitService property_unit.IService, propertyViewTypeService property_view_type.IService, propertyGroupService property_group.IService, langFinder text_lang.LangFinder) IService {
	s := &service{
		logger:                  logger,
		repository:              repo,
		propertyTypeService:     propertyTypeService,
		propertyUnitService:     propertyUnitService,
		propertyViewTypeService: propertyViewTypeService,
		propertyGroupService:    propertyGroupService,
		langFinder:              langFinder,
	}
	repo.SetDefaultConditions(s.defaultConditions())
	return s
}

// Defaults returns defaults params
func (s *service) defaultConditions() *selection_condition.SelectionCondition {
	return &selection_condition.SelectionCondition{}
}

func (s *service) NewEntity() *Property {
	return New()
}

func (s *service) PropertyInit(ctx context.Context, PropertiesConfig config.Properties, entityTypeSysname string) (propertyIDs []uint, err error) {
	var isNeedUpdate bool

	langIDEng, err := s.langFinder.GetIDByCode(ctx, text_lang.CodeEng)

	langsSl, err := s.langFinder.GetCodesEmptyInterfaceSlice(ctx)
	if err != nil {
		return nil, err
	}

	langsIDsMap, err := s.langFinder.GetMapCodeID(ctx)
	if err != nil {
		return nil, err
	}

	for _, propertyConfig := range PropertiesConfig {
		prop := New()
		prop.Sysname = propertyConfig.Sysname

		if propertyConfig.PropertyType != "" {
			if prop.PropertyTypeID, err = s.propertyTypeService.GetIDBySysname(ctx, propertyConfig.PropertyType); err != nil {
				return nil, err
			}
		}

		if propertyConfig.PropertyUnit != "" {
			propertyUnitID, err := s.propertyUnitService.GetIDBySysname(ctx, propertyConfig.PropertyUnit)
			if err != nil {
				return nil, err
			}
			prop.PropertyUnitID = &propertyUnitID
		}

		if propertyConfig.PropertyViewType != "" {
			propertyViewTypeID, err := s.propertyViewTypeService.GetIDBySysname(ctx, propertyConfig.PropertyViewType)
			if err != nil {
				return nil, err
			}
			prop.PropertyViewTypeID = &propertyViewTypeID
		}

		if propertyConfig.PropertyGroup != "" {
			propertyGroupID, err := s.propertyGroupService.GetIDBySysname(ctx, propertyConfig.PropertyGroup)
			if err != nil {
				return nil, err
			}
			prop.PropertyGroupID = &propertyGroupID
		}

		prop.IsSpecific = propertyConfig.IsSpecific
		prop.IsRange = propertyConfig.IsRange
		prop.IsMultiple = propertyConfig.IsMultiple
		prop.SortOrder = propertyConfig.SortOrder
		prop.Options = propertyConfig.Options
		if err := s.UpsertBySysname(ctx, prop, langIDEng); err != nil {
			return nil, err
		}

		for lang, texts := range propertyConfig.Texts {
			if err := validation.Validate(lang, validation.In(langsSl...)); err != nil {
				return nil, errors.Wrapf(err, "PropertyInit error: invalid lang = %q", lang)
			}
			langID, ok := langsIDsMap[lang]
			if !ok {
				return nil, errors.Errorf("PropertyInit error: not found lang = %q", lang)
			}

			if prop, err = s.TGet(ctx, prop.ID, langID); err != nil {
				return nil, err
			}

			if texts.Name != "" {
				name := texts.Name
				prop.Name = &name
				isNeedUpdate = true
			}

			if texts.Description != "" {
				description := texts.Description
				prop.Description = &description
				isNeedUpdate = true
			}

			if isNeedUpdate {
				if err := s.TUpdate(ctx, prop, langID); err != nil {
					return nil, err
				}
			}
		}
		propertyIDs = append(propertyIDs, prop.ID)
	}

	return propertyIDs, nil
}

// Get returns the entity with the specified ID.
func (s *service) Get(ctx context.Context, id uint) (*Property, error) {
	entity, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *service) GetBySysname(ctx context.Context, sysname string, langID uint) (*Property, error) {
	return s.repository.GetBySysname(ctx, sysname, langID)
}

func (s *service) First(ctx context.Context, entity *Property) (*Property, error) {
	entity, err := s.repository.First(ctx, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *service) TGet(ctx context.Context, id uint, langID uint) (*Property, error) {
	entity, err := s.repository.TGet(ctx, id, langID)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *service) TFirst(ctx context.Context, entity *Property, langID uint) (*Property, error) {
	entity, err := s.repository.TFirst(ctx, entity, langID)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// Query returns the items with the specified selection condition.
func (s *service) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]Property, error) {
	items, err := s.repository.Query(ctx, cond)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of items by query: %v", cond)
	}
	return items, nil
}

func (s *service) TQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]Property, error) {
	items, err := s.repository.TQuery(ctx, cond, langID)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of items by query: %v", cond)
	}
	return items, nil
}

func (s *service) GetSysnames(ctx context.Context) ([]string, error) {
	items, err := s.Query(ctx, &selection_condition.SelectionCondition{})
	if err != nil {
		return nil, err
	}

	res := make([]string, len(items))
	for i, item := range items {
		res[i] = item.Sysname
	}
	return res, nil
}

func (s *service) GetSysnamesEmptyInterfaceSlice(ctx context.Context) ([]interface{}, error) {
	items, err := s.Query(ctx, &selection_condition.SelectionCondition{})
	if err != nil {
		return nil, err
	}

	res := make([]interface{}, len(items))
	for i, item := range items {
		res[i] = item.Sysname
	}
	return res, nil
}

func (s *service) GetMapSysnameID(ctx context.Context) (map[string]uint, error) {
	items, err := s.Query(ctx, &selection_condition.SelectionCondition{})
	if err != nil {
		return nil, err
	}

	res := make(map[string]uint, len(items))
	for _, item := range items {
		res[item.Sysname] = item.ID
	}
	return res, nil
}

func (s *service) GetMapSysnameTypeID(ctx context.Context) (map[string]uint, error) {
	items, err := s.Query(ctx, &selection_condition.SelectionCondition{})
	if err != nil {
		return nil, err
	}

	res := make(map[string]uint, len(items))
	for _, item := range items {
		res[item.Sysname] = item.PropertyTypeID
	}
	return res, nil
}

func (s *service) GetMapIDTypeID(ctx context.Context) (map[uint]uint, error) {
	items, err := s.Query(ctx, &selection_condition.SelectionCondition{})
	if err != nil {
		return nil, err
	}

	res := make(map[uint]uint, len(items))
	for _, item := range items {
		res[item.ID] = item.PropertyTypeID
	}
	return res, nil
}

func (s *service) GetMapIDSysname(ctx context.Context) (map[uint]string, error) {
	items, err := s.Query(ctx, &selection_condition.SelectionCondition{})
	if err != nil {
		return nil, err
	}

	res := make(map[uint]string, len(items))
	for _, item := range items {
		res[item.ID] = item.Sysname
	}
	return res, nil
}

func (s *service) GetIDBySysname(ctx context.Context, sysname string) (uint, error) {
	mapSysnameID, err := s.GetMapSysnameID(ctx)
	if err != nil {
		return 0, err
	}

	id, ok := mapSysnameID[sysname]
	if !ok {
		return 0, yaruserror.ErrNotFound
	}

	return id, nil
}

func (s *service) GetSysnameByID(ctx context.Context, id uint) (string, error) {
	mapIDSysname, err := s.GetMapIDSysname(ctx)
	if err != nil {
		return "", err
	}

	sysname, ok := mapIDSysname[id]
	if !ok {
		return "", yaruserror.ErrNotFound
	}

	return sysname, nil
}

func (s *service) GetTypeIDBySysname(ctx context.Context, sysname string) (uint, error) {
	mapSysnameTypeID, err := s.GetMapSysnameTypeID(ctx)
	if err != nil {
		return 0, err
	}

	typeID, ok := mapSysnameTypeID[sysname]
	if !ok {
		return 0, yaruserror.ErrNotFound
	}

	return typeID, nil
}

func (s *service) GetTypeIDByID(ctx context.Context, id uint) (uint, error) {
	mapIDTypeID, err := s.GetMapIDTypeID(ctx)
	if err != nil {
		return 0, err
	}

	typeID, ok := mapIDTypeID[id]
	if !ok {
		return 0, yaruserror.ErrNotFound
	}

	return typeID, nil
}

func (s *service) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
	count, err := s.repository.Count(ctx, cond)
	if err != nil {
		return 0, errors.Wrapf(err, "Can not count a list of items by query: %v", cond)
	}
	return count, nil
}

func (s *service) UpsertBySysname(ctx context.Context, entity *Property, langID uint) (err error) {
	found, err := s.repository.First(ctx, &Property{
		Sysname: entity.Sysname,
	})

	if err != nil {
		if err != yaruserror.ErrNotFound {
			return err
		}
		err = s.TCreate(ctx, entity, langID)
	} else {
		entity.ID = found.ID
		entity.NameSourceID = found.NameSourceID
		entity.DescriptionSourceID = found.DescriptionSourceID
		entity.CreatedAt = found.CreatedAt
		if found.SortOrder != SortOrderDefault {
			entity.SortOrder = found.SortOrder
		}
		err = s.TUpdate(ctx, entity, langID)
	}

	return err
}

func (s *service) Create(ctx context.Context, entity *Property) error {
	err := entity.Validate()
	if err != nil {
		return errors.Wrapf(err, "Validation error: %v", err)
	}

	err = s.repository.Create(ctx, entity)
	if err != nil {
		return errors.Wrapf(err, "Can not create an entity: %v", entity)
	}
	return nil
}

func (s *service) TCreate(ctx context.Context, entity *Property, langID uint) error {
	err := entity.Validate()
	if err != nil {
		return errors.Wrapf(err, "Validation error: %v", err)
	}

	err = s.repository.TCreate(ctx, entity, langID)
	if err != nil {
		return errors.Wrapf(err, "Can not create an entity: %v", entity)
	}
	return nil
}

func (s *service) Update(ctx context.Context, entity *Property) error {
	err := entity.Validate()
	if err != nil {
		return errors.Wrapf(err, "Validation error: %v", err)
	}

	err = s.repository.Update(ctx, entity)
	if err != nil {
		return errors.Wrapf(err, "Can not update an entity: %v", entity)
	}
	return nil
}

func (s *service) TUpdate(ctx context.Context, entity *Property, langID uint) error {
	err := entity.Validate()
	if err != nil {
		return errors.Wrapf(err, "Validation error: %v", err)
	}

	err = s.repository.TUpdate(ctx, entity, langID)
	if err != nil {
		return errors.Wrapf(err, "Can not update an entity: %v", entity)
	}
	return nil
}

func (s *service) Delete(ctx context.Context, entity *Property) error {
	err := s.repository.Delete(ctx, entity)
	if err != nil {
		return errors.Wrapf(err, "Can not delete an entity: %v", entity)
	}
	return nil
}
