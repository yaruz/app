package text_lang

import (
	"context"

	"github.com/yaruz/app/pkg/yarus_platform/config"

	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"

	"github.com/minipkg/selection_condition"

	"github.com/pkg/errors"

	"github.com/minipkg/log"
)

type LangFinder interface {
	GetCodes(ctx context.Context) ([]string, error)
	GetCodesEmptyInterfaceSlice(ctx context.Context) ([]interface{}, error)
	GetMapCodeID(ctx context.Context) (map[string]uint, error)
	GetMapIDCode(ctx context.Context) (map[uint]string, error)
	GetMapIDCfgname(ctx context.Context) (map[uint]string, error)
	GetIDByCode(ctx context.Context, code string) (uint, error)
	GetCfgnameByID(ctx context.Context, id uint) (string, error)
}

// IService encapsulates usecase logic.
type IService interface {
	NewEntity() *TextLang
	DataInit(ctx context.Context, langsConfig config.Languages) error
	Get(ctx context.Context, id uint) (*TextLang, error)
	Query(ctx context.Context, query *selection_condition.SelectionCondition) ([]TextLang, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *TextLang) error
	Upsert(ctx context.Context, entity *TextLang) error
	Update(ctx context.Context, entity *TextLang) error
	Save(ctx context.Context, entity *TextLang) error
	Delete(ctx context.Context, id uint) error
	GetCodes(ctx context.Context) ([]string, error)
	GetCodesEmptyInterfaceSlice(ctx context.Context) ([]interface{}, error)
	GetMapCodeID(ctx context.Context) (map[string]uint, error)
	GetMapIDCode(ctx context.Context) (map[uint]string, error)
	GetMapIDCfgname(ctx context.Context) (map[uint]string, error)
	GetIDByCode(ctx context.Context, code string) (uint, error)
	GetCfgnameByID(ctx context.Context, id uint) (string, error)
}

const (
	CodeEng = "eng"
)

type service struct {
	logger     log.Logger
	repository Repository
}

var _ IService = (*service)(nil)

// NewService creates a new service.
func NewService(logger log.Logger, repo Repository) IService {
	s := &service{
		logger:     logger,
		repository: repo,
	}
	repo.SetDefaultConditions(s.defaultConditions())
	return s
}

// Defaults returns defaults params
func (s *service) defaultConditions() *selection_condition.SelectionCondition {
	return &selection_condition.SelectionCondition{}
}

func (s *service) NewEntity() *TextLang {
	return New()
}

func (s *service) DataInit(ctx context.Context, langsConfig config.Languages) error {

	for _, langConfig := range langsConfig {
		item := New()
		item.Code = langConfig.Code
		item.Name = langConfig.Name
		item.Cfgname = langConfig.Cfgname
		s.Upsert(ctx, item)
	}

	return nil
}

// Get returns the entity with the specified ID.
func (s *service) Get(ctx context.Context, id uint) (*TextLang, error) {
	entity, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// Query returns the items with the specified selection condition.
func (s *service) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]TextLang, error) {
	items, err := s.repository.Query(ctx, cond)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of items by query: %v", cond)
	}
	return items, nil
}

func (s *service) GetCodes(ctx context.Context) ([]string, error) {
	items, err := s.Query(ctx, &selection_condition.SelectionCondition{})
	if err != nil {
		return nil, err
	}

	res := make([]string, len(items))
	for i, item := range items {
		res[i] = item.Code
	}
	return res, nil
}

func (s *service) GetCodesEmptyInterfaceSlice(ctx context.Context) ([]interface{}, error) {
	items, err := s.Query(ctx, &selection_condition.SelectionCondition{})
	if err != nil {
		return nil, err
	}

	res := make([]interface{}, len(items))
	for i, item := range items {
		res[i] = item.Code
	}
	return res, nil
}

func (s *service) GetMapCodeID(ctx context.Context) (map[string]uint, error) {
	items, err := s.Query(ctx, &selection_condition.SelectionCondition{})
	if err != nil {
		return nil, err
	}

	res := make(map[string]uint, len(items))
	for _, item := range items {
		res[item.Code] = item.ID
	}
	return res, nil
}

func (s *service) GetMapIDCode(ctx context.Context) (map[uint]string, error) {
	items, err := s.Query(ctx, &selection_condition.SelectionCondition{})
	if err != nil {
		return nil, err
	}

	res := make(map[uint]string, len(items))
	for _, item := range items {
		res[item.ID] = item.Code
	}
	return res, nil
}

func (s *service) GetMapIDCfgname(ctx context.Context) (map[uint]string, error) {
	items, err := s.Query(ctx, &selection_condition.SelectionCondition{})
	if err != nil {
		return nil, err
	}

	res := make(map[uint]string, len(items))
	for _, item := range items {
		res[item.ID] = item.Cfgname
	}
	return res, nil
}

func (s *service) GetIDByCode(ctx context.Context, code string) (uint, error) {
	mapSysnameID, err := s.GetMapCodeID(ctx)
	if err != nil {
		return 0, err
	}

	id, ok := mapSysnameID[code]
	if !ok {
		return 0, yaruserror.ErrNotFound
	}

	return id, nil
}

func (s *service) GetCfgnameByID(ctx context.Context, id uint) (string, error) {
	mapIDCfgname, err := s.GetMapIDCfgname(ctx)
	if err != nil {
		return "", err
	}

	cfgname, ok := mapIDCfgname[id]
	if !ok {
		return "", yaruserror.ErrNotFound
	}

	return cfgname, nil
}

// GetMapIDCfgname(ctx context.Context) (map[uint]string, error)

func (s *service) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
	count, err := s.repository.Count(ctx, cond)
	if err != nil {
		return 0, errors.Wrapf(err, "Can not count a list of items by query: %v", cond)
	}
	return count, nil
}

func (s *service) Create(ctx context.Context, entity *TextLang) error {
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

func (s *service) Upsert(ctx context.Context, entity *TextLang) error {
	err := entity.Validate()
	if err != nil {
		return errors.Wrapf(err, "Validation error: %v", err)
	}

	err = s.repository.Upsert(ctx, entity)
	if err != nil {
		return errors.Wrapf(err, "Can not create an entity: %v", entity)
	}
	return nil
}

func (s *service) Update(ctx context.Context, entity *TextLang) error {
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

func (s *service) Save(ctx context.Context, entity *TextLang) error {
	err := entity.Validate()
	if err != nil {
		return errors.Wrapf(err, "Validation error: %v", err)
	}

	err = s.repository.Save(ctx, entity)
	if err != nil {
		return errors.Wrapf(err, "Can not save an entity: %v", entity)
	}
	return nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "Can not delete an entity by ID: %v", id)
	}
	return nil
}
