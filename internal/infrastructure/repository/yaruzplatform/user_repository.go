package yaruzplatform

import (
	"context"

	"github.com/yaruz/app/internal/pkg/config"
	yaruz_config "github.com/yaruz/app/pkg/yarus_platform/config"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"

	"github.com/yaruz/app/internal/domain/user"

	"github.com/yaruz/app/internal/domain/task"

	"github.com/minipkg/selection_condition"
)

// TaskRepository is a repository for the model entity
type UserRepository struct {
	repository
}

var _ user.Repository = (*UserRepository)(nil)

// New creates a new UserRepository
func NewUserRepository(repository *repository) (*UserRepository, error) {
	return &UserRepository{repository: *repository}, nil
}

func (r *UserRepository) PropertiesConfig(ctx context.Context) property.Configs {
	return property.Configs{
		"Name": property.Config{
			PropertyTypeID: property_type.IDText,
			Texts: map[string]yaruz_config.NameAndDescriptionText{
				config.LangEng: {
					Name:        "Name",
					Description: "Name",
				},
				config.LangRus: {
					Name:        "Имя",
					Description: "Имя",
				},
			},
		},
	}
}

// Get reads the album with the specified ID from the database.
func (r *UserRepository) Get(ctx context.Context, id uint) (*task.Task, error) {
	entity := &task.Task{}

	//err := r.DB().First(entity, id).Error
	//if err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		return entity, yaruserror.ErrNotFound
	//	}
	//}

	return entity, nil
}

func (r *UserRepository) First(ctx context.Context, entity *task.Task) (*task.Task, error) {
	//err := r.DB().Where(entity).First(entity).Error
	//if err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		return entity, yaruserror.ErrNotFound
	//	}
	//}

	return entity, nil
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *UserRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]task.Task, error) {
	items := []task.Task{}
	//db := minipkg_gorm.Conditions(r.DB(), cond)
	//if db.Error != nil {
	//	return nil, db.Error
	//}
	//
	//err := db.Find(&items).Error
	//if err != nil {
	//	if err == gorm.ErrRecordNotFound {
	//		return items, yaruserror.ErrNotFound
	//	}
	//}

	return items, nil
}

func (r *UserRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
	var count int64
	//c := cond
	//c.Limit = 0
	//c.Offset = 0
	//db := minipkg_gorm.Conditions(r.DB(), cond)
	//if db.Error != nil {
	//	return 0, db.Error
	//}
	//
	//err := db.Count(&count).Error
	return count, nil
}

// Create saves a new record in the database.
func (r *UserRepository) Create(ctx context.Context, entity *task.Task) error {

	//if entity.ID > 0 {
	//	return errors.New("entity is not new")
	//}
	//return r.db.DB().Create(entity).Error
	return nil
}

// Update saves a changed Maintenance record in the database.
func (r *UserRepository) Update(ctx context.Context, entity *task.Task) error {

	//if entity.ID == 0 {
	//	return errors.New("entity is new")
	//}
	//
	//return r.Save(ctx, entity)
	return nil
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *UserRepository) Save(ctx context.Context, entity *task.Task) error {
	//return r.db.DB().Save(entity).Error
	return nil
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *UserRepository) Delete(ctx context.Context, id uint) error {

	//err := r.db.DB().Delete(&task.Task{}, id).Error
	//if err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		return apperror.ErrNotFound
	//	}
	//}
	//return err
	return nil
}
