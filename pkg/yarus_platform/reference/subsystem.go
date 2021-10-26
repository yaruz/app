package reference

import (
	"context"
	"fmt"
	golog "log"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/yaruz/app/pkg/yarus_platform/config"

	"github.com/yaruz/app/pkg/yarus_platform/infrastructure"

	"github.com/minipkg/selection_condition"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/log"
	"github.com/pkg/errors"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type2property"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_group"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type2property_view_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_unit"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_view_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_lang"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_source"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_value"
	"github.com/yaruz/app/pkg/yarus_platform/reference/infrastructure/repository/gorm"
)

type ReferenceSubsystem struct {
	EntityType                    ReferenceDomainEntityType
	EntityType2Property           ReferenceDomainEntityType2Property
	Property                      ReferenceDomainProperty
	Relation                      ReferenceDomainRelation
	PropertyGroup                 ReferenceDomainPropertyGroup
	PropertyType                  ReferenceDomainPropertyType
	PropertyType2PropertyViewType ReferenceDomainPropertyType2PropertyViewType
	PropertyUnit                  ReferenceDomainPropertyUnit
	PropertyViewType              ReferenceDomainPropertyViewType
	TextSource                    ReferenceDomainTextSource
	TextValue                     ReferenceDomainTextValue
	TextLang                      ReferenceDomainTextLang
}

type ReferenceDomainEntityType struct {
	Service    entity_type.IService
	Repository entity_type.Repository
}

type ReferenceDomainEntityType2Property struct {
	Service    entity_type2property.IService
	Repository entity_type2property.Repository
}

type ReferenceDomainProperty struct {
	Service    property.IService
	Repository property.Repository
}

type ReferenceDomainRelation struct {
	Service    entity_type.RelationService
	Repository entity_type.RelationRepository
}

type ReferenceDomainPropertyGroup struct {
	Service    property_group.IService
	Repository property_group.Repository
}

type ReferenceDomainPropertyType struct {
	Service    property_type.IService
	Repository property_type.Repository
}

type ReferenceDomainPropertyType2PropertyViewType struct {
	Repository property_type2property_view_type.Repository
}

type ReferenceDomainPropertyUnit struct {
	Service    property_unit.IService
	Repository property_unit.Repository
}

type ReferenceDomainPropertyViewType struct {
	Service    property_view_type.IService
	Repository property_view_type.Repository
}

type ReferenceDomainTextSource struct {
	Service    text_source.IService
	Repository text_source.Repository
}

type ReferenceDomainTextValue struct {
	Service    text_value.IService
	Repository text_value.Repository
}

type ReferenceDomainTextLang struct {
	Service    text_lang.IService
	Repository text_lang.Repository
}

func NewReferenceSubsystem(infra *infrastructure.Infrastructure, referenceData config.ReferenceData) (*ReferenceSubsystem, error) {
	s := &ReferenceSubsystem{}
	if err := s.setupRepositories(infra); err != nil {
		return nil, err
	}
	s.setupServices(infra.Logger)

	if err := s.autoMigrate(infra.ReferenceDB); err != nil {
		return nil, err
	}

	if err := s.dbDataInit(referenceData); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *ReferenceSubsystem) autoMigrate(db minipkg_gorm.IDB) error {
	if db.IsAutoMigrate() {
		err := db.DB().SetupJoinTable(&entity_type.EntityType{}, "Properties", &entity_type2property.EntityType2Property{})
		if err != nil {
			return err
		}

		err = db.DB().AutoMigrate(
			&text_lang.TextLang{},
			&text_source.TextSource{},
			&text_value.TextValue{},
			&property_unit.PropertyUnit{},
			&property_group.PropertyGroup{},
			&property_type.PropertyType{},
			&property_view_type.PropertyViewType{},
			&property.Property{},
			&entity_type.EntityType{},
		)
		if err != nil {
			return err
		}

		if err = s.dbStructFix(db); err != nil {
			return err
		}
	}
	return nil
}

func (s *ReferenceSubsystem) setupRepositories(infra *infrastructure.Infrastructure) (err error) {
	var ok bool

	repo, err := gorm.GetRepository(infra.Logger, infra.ReferenceDB, text_source.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", text_source.EntityName, err)
	}
	s.TextSource.Repository, ok = repo.(text_source.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", text_source.EntityName, text_source.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, property_type.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_type.EntityName, err)
	}
	s.PropertyType.Repository, ok = repo.(property_type.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_type.EntityName, property_type.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, property_view_type.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_view_type.EntityName, err)
	}
	s.PropertyViewType.Repository, ok = repo.(property_view_type.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_view_type.EntityName, property_view_type.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, entity_type.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", entity_type.EntityName, err)
	}
	s.EntityType.Repository, ok = repo.(entity_type.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", entity_type.EntityName, entity_type.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, property_unit.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_unit.EntityName, err)
	}
	s.PropertyUnit.Repository, ok = repo.(property_unit.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_unit.EntityName, property_unit.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, property_group.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_group.EntityName, err)
	}
	s.PropertyGroup.Repository, ok = repo.(property_group.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_group.EntityName, property_group.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, property_type2property_view_type.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_type2property_view_type.EntityName, err)
	}
	s.PropertyType2PropertyViewType.Repository, ok = repo.(property_type2property_view_type.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_type2property_view_type.EntityName, property_type2property_view_type.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, property.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property.EntityName, err)
	}
	s.Property.Repository, ok = repo.(property.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property.EntityName, property.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, text_value.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", text_value.EntityName, err)
	}
	s.TextValue.Repository, ok = repo.(text_value.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", text_value.EntityName, text_value.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, entity_type2property.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", entity_type2property.EntityName, err)
	}
	s.EntityType2Property.Repository, ok = repo.(entity_type2property.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", entity_type2property.EntityName, entity_type2property.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, entity_type.RelationEntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", entity_type.RelationEntityName, err)
	}
	s.Relation.Repository, ok = repo.(entity_type.RelationRepository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", entity_type.RelationEntityName, entity_type.RelationEntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, text_lang.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", text_lang.EntityName, err)
	}
	s.TextLang.Repository, ok = repo.(text_lang.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", text_lang.EntityName, text_lang.EntityName, repo)
	}

	return nil
}

func (s *ReferenceSubsystem) setupServices(logger log.ILogger) {
	s.EntityType2Property.Service = entity_type2property.NewService(logger, s.EntityType2Property.Repository)
	s.Relation.Service = entity_type.NewRelationService(logger, s.Relation.Repository)
	s.EntityType.Service = entity_type.NewService(logger, s.EntityType.Repository, s.Relation.Service)
	s.PropertyGroup.Service = property_group.NewService(logger, s.PropertyGroup.Repository)
	s.Property.Service = property.NewService(logger, s.Property.Repository)
	s.PropertyType.Service = property_type.NewService(logger, s.PropertyType.Repository, s.PropertyType2PropertyViewType.Repository)
	s.PropertyUnit.Service = property_unit.NewService(logger, s.PropertyUnit.Repository)
	s.PropertyViewType.Service = property_view_type.NewService(logger, s.PropertyViewType.Repository)
	s.TextSource.Service = text_source.NewService(logger, s.TextSource.Repository)
	s.TextValue.Service = text_value.NewService(logger, s.TextValue.Repository)
	s.TextLang.Service = text_lang.NewService(logger, s.TextLang.Repository)
}

func (s *ReferenceSubsystem) dbStructFix(db minipkg_gorm.IDB) error {
	typeSpec := "DROP NOT NULL"
	fields := []string{
		"property_unit_id",
		"property_view_type_id",
		"property_group_id",
	}
	for _, f := range fields {
		sql := fmt.Sprintf("ALTER TABLE %q ALTER COLUMN %q %s",
			property.TableName,
			f,
			typeSpec,
		)
		err := db.DB().Exec(sql).Error

		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ReferenceSubsystem) dbDataInit(referenceData config.ReferenceData) error {
	ctx := context.Background()

	err := s.LangDataInit(ctx, referenceData.Languages)
	if err != nil {
		return err
	}

	err = s.PropertyTypeDataInit(ctx)
	if err != nil {
		return err
	}

	err = s.PropertyUnitDataInit(ctx, referenceData.PropertyUnits)
	if err != nil {
		return err
	}

	return nil
}

func (s *ReferenceSubsystem) LangDataInit(ctx context.Context, langsConfig config.Languages) error {
	count, err := s.TextLang.Service.Count(ctx, &selection_condition.SelectionCondition{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	for _, langConfig := range langsConfig {
		item := text_lang.New()
		item.Code = langConfig.Code
		item.Name = langConfig.Name
		s.TextLang.Service.Create(ctx, item)
	}

	return nil
}

func (s *ReferenceSubsystem) PropertyTypeDataInit(ctx context.Context) error {
	count, err := s.PropertyType.Service.Count(ctx, &selection_condition.SelectionCondition{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	items := []property_type.PropertyType{
		{
			Sysname: property_type.SysnameRelation,
		},
		{
			Sysname: property_type.SysnameBoolean,
		},
		{
			Sysname: property_type.SysnameInt,
		},
		{
			Sysname: property_type.SysnameFloat,
		},
		{
			Sysname: property_type.SysnameDate,
		},
		{
			Sysname: property_type.SysnameTime,
		},
		{
			Sysname: property_type.SysnameText,
		},
	}
	for _, i := range items {
		err = s.PropertyType.Service.Create(ctx, &i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ReferenceSubsystem) PropertyUnitDataInit(ctx context.Context, unitsConfig config.PropertyUnits) error {
	count, err := s.PropertyUnit.Service.Count(ctx, &selection_condition.SelectionCondition{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	langsSl, err := s.TextLang.Service.GetCodesEmptyInterfaceSlice(ctx)
	if err != nil {
		return err
	}

	langsIDsMap, err := s.TextLang.Service.GetMapCodeID(ctx)
	if err != nil {
		return err
	}

	for _, unitConfig := range unitsConfig {
		unit := property_unit.New()
		unit.Sysname = unitConfig.Sysname
		if err := s.PropertyUnit.Service.TCreate(ctx, unit, 1); err != nil {
			return err
		}

		for lang, texts := range unitConfig.Texts {
			if err := validation.Validate(lang, validation.In(langsSl...)); err != nil {
				return errors.Wrapf(err, "PropertyUnitDataInit error: invalid lang = %q", lang)
			}
			langID, ok := langsIDsMap[lang]
			if !ok {
				return errors.Errorf("PropertyUnitDataInit error: not found lang = %q", lang)
			}

			unitName := texts.Name
			unitDescription := texts.Description
			unit.Name = &unitName
			unit.Description = &unitDescription
			if err := s.PropertyUnit.Service.TUpdate(ctx, unit, langID); err != nil {
				return err
			}
		}
	}

	return nil
}
