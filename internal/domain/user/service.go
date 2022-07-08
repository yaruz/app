package user

import (
	"context"

	"github.com/pkg/errors"

	"github.com/minipkg/log"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/internal/domain/account"
	"github.com/yaruz/app/internal/domain/tg_account"
)

// IService encapsulates usecase logic for user.
type IService interface {
	New(ctx context.Context) (*User, error)
	Get(ctx context.Context, id uint, langID uint) (*User, error)
	GetByAccountID(ctx context.Context, accountId string, langID uint) (*User, error)
	Query(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) ([]User, error)
	First(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (*User, error)
	Count(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (uint, error)
	Create(ctx context.Context, obj *User, langID uint) error
	Update(ctx context.Context, obj *User, langID uint) error
	Delete(ctx context.Context, id uint) error
	AccountSettingsValidate(ctx context.Context, accountSettings *account.AccountSettings) error
	GetTgAccount(ctx context.Context, obj *User, langID uint) (*tg_account.TgAccount, error)
}

type service struct {
	//Domain     Domain
	logger     log.ILogger
	repository Repository
}

var _ IService = (*service)(nil)

// NewService creates a new service.
func NewService(logger log.ILogger, repo Repository) IService {
	s := &service{
		logger:     logger,
		repository: repo,
	}
	repo.SetDefaultConditions(s.defaultConditions())
	return s
}

// Defaults returns defaults params
func (s service) defaultConditions() *selection_condition.SelectionCondition {
	return &selection_condition.SelectionCondition{}
}

func (s *service) New(ctx context.Context) (*User, error) {
	return s.repository.New(ctx)
}

// Get returns the entity with the specified ID.
func (s *service) Get(ctx context.Context, id uint, langID uint) (*User, error) {
	entity, err := s.repository.Get(ctx, id, langID)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not get a %v by id: %v", EntityType, id)
	}
	return entity, nil
}

func (s *service) GetByAccountID(ctx context.Context, accountId string, langID uint) (*User, error) {
	return s.First(ctx, &selection_condition.SelectionCondition{
		Where: selection_condition.WhereCondition{
			Field:     PropertySysnameAccountID,
			Condition: selection_condition.ConditionEq,
			Value:     accountId,
		},
	}, langID)
}

// Query returns the items with the specified offset and limit.
func (s *service) Query(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) ([]User, error) {
	return s.repository.Query(ctx, condition, langID)
}

func (s *service) First(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (*User, error) {
	return s.repository.First(ctx, condition, langID)
}

// Count returns the number of items.
func (s *service) Count(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (uint, error) {
	return s.repository.Count(ctx, condition, langID)
}

func (s *service) Create(ctx context.Context, obj *User, langID uint) error {
	return s.repository.Create(ctx, obj, langID)
}

func (s *service) Update(ctx context.Context, obj *User, langID uint) error {
	return s.repository.Update(ctx, obj, langID)
}

func (s *service) Delete(ctx context.Context, id uint) error {
	return s.repository.Delete(ctx, id)
}

func (s *service) AccountSettingsValidate(ctx context.Context, accountSettings *account.AccountSettings) error {
	err := s.repository.LangIDValidate(ctx, accountSettings.LangID)
	if err != nil {
		return errors.Errorf("LangID: %v", err.Error())
	}
	return nil
}

func (s *service) GetTgAccount(ctx context.Context, obj *User, langID uint) (*tg_account.TgAccount, error) {
	return s.repository.GetTgAccount(ctx, obj, langID)
}
