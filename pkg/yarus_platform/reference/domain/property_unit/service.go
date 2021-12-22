package property_unit

import (
	"context"

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
	NewEntity() *PropertyUnit
	DataInit(ctx context.Context, unitsConfig config.PropertyUnits) error
	Get(ctx context.Context, id uint) (*PropertyUnit, error)
	First(ctx context.Context, entity *PropertyUnit) (*PropertyUnit, error)
	Query(ctx context.Context, query *selection_condition.SelectionCondition) ([]PropertyUnit, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *PropertyUnit) error
	Update(ctx context.Context, entity *PropertyUnit) error
	Delete(ctx context.Context, entity *PropertyUnit) error
	TGet(ctx context.Context, id uint, langID uint) (*PropertyUnit, error)
	TFirst(ctx context.Context, entity *PropertyUnit, langID uint) (*PropertyUnit, error)
	TQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]PropertyUnit, error)
	TCreate(ctx context.Context, entity *PropertyUnit, langID uint) (err error)
	TUpdate(ctx context.Context, entity *PropertyUnit, langID uint) (err error)
	GetSysnames(ctx context.Context) ([]string, error)
	GetSysnamesEmptyInterfaceSlice(ctx context.Context) ([]interface{}, error)
	GetMapSysnameID(ctx context.Context) (map[string]uint, error)
	GetMapIDSysname(ctx context.Context) (map[uint]string, error)
	GetIDBySysname(ctx context.Context, sysname string) (uint, error)
	GetSysnameByID(ctx context.Context, id uint) (string, error)
}

type service struct {
	logger     log.ILogger
	repository Repository
	langFinder text_lang.LangFinder
}

var _ IService = (*service)(nil)

// NewService creates a new service.
func NewService(logger log.ILogger, repo Repository, langFinder text_lang.LangFinder) IService {
	s := &service{
		logger:     logger,
		repository: repo,
		langFinder: langFinder,
	}
	repo.SetDefaultConditions(s.defaultConditions())
	return s
}

// Defaults returns defaults params
func (s *service) defaultConditions() *selection_condition.SelectionCondition {
	return &selection_condition.SelectionCondition{}
}

func (s *service) NewEntity() *PropertyUnit {
	return New()
}

func (s *service) DataInit(ctx context.Context, unitsConfig config.PropertyUnits) error {
	langsSl, err := s.langFinder.GetCodesEmptyInterfaceSlice(ctx)
	if err != nil {
		return err
	}

	langsIDsMap, err := s.langFinder.GetMapCodeID(ctx)
	if err != nil {
		return err
	}

	for sysname, unitConfig := range unitsConfig {
		unit := New()
		unit.Sysname = sysname
		if err := s.UpsertBySysname(ctx, unit, 1); err != nil {
			return err
		}

		for lang, texts := range unitConfig {
			if err := validation.Validate(lang, validation.In(langsSl...)); err != nil {
				return errors.Wrapf(err, "PropertyUnitInit error: invalid lang = %q", lang)
			}
			langID, ok := langsIDsMap[lang]
			if !ok {
				return errors.Errorf("PropertyUnitInit error: not found lang = %q", lang)
			}

			if unit, err = s.TGet(ctx, unit.ID, langID); err != nil {
				return err
			}

			name := texts.Name
			description := texts.Description
			unit.Name = &name
			unit.Description = &description
			if err := s.TUpdate(ctx, unit, langID); err != nil {
				return err
			}
		}
	}

	return nil
}

// Get returns the entity with the specified ID.
func (s *service) Get(ctx context.Context, id uint) (*PropertyUnit, error) {
	entity, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *service) First(ctx context.Context, entity *PropertyUnit) (*PropertyUnit, error) {
	entity, err := s.repository.First(ctx, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *service) TGet(ctx context.Context, id uint, langID uint) (*PropertyUnit, error) {
	entity, err := s.repository.TGet(ctx, id, langID)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *service) TFirst(ctx context.Context, entity *PropertyUnit, langID uint) (*PropertyUnit, error) {
	entity, err := s.repository.TFirst(ctx, entity, langID)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// Query returns the items with the specified selection condition.
func (s *service) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]PropertyUnit, error) {
	items, err := s.repository.Query(ctx, cond)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of items by query: %v", cond)
	}
	return items, nil
}

func (s *service) TQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]PropertyUnit, error) {
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

func (s *service) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
	count, err := s.repository.Count(ctx, cond)
	if err != nil {
		return 0, errors.Wrapf(err, "Can not count a list of items by query: %v", cond)
	}
	return count, nil
}

func (s *service) UpsertBySysname(ctx context.Context, entity *PropertyUnit, langID uint) (err error) {
	found, err := s.repository.First(ctx, &PropertyUnit{
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
		err = s.TUpdate(ctx, entity, langID)
	}

	return err
}

func (s *service) Create(ctx context.Context, entity *PropertyUnit) error {
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

func (s *service) TCreate(ctx context.Context, entity *PropertyUnit, langID uint) error {
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

func (s *service) Update(ctx context.Context, entity *PropertyUnit) error {
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

func (s *service) TUpdate(ctx context.Context, entity *PropertyUnit, langID uint) error {
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

func (s *service) Delete(ctx context.Context, entity *PropertyUnit) error {
	err := s.repository.Delete(ctx, entity)
	if err != nil {
		return errors.Wrapf(err, "Can not delete an entity: %v", entity)
	}
	return nil
}
