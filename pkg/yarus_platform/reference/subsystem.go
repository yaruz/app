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

type EntityIDRepository interface {
	AutoMigrate() error
	NextVal(entityTypeSysname string) (id uint, err error)
	LastVal(entityTypeSysname string) (id uint, err error)
}

type ReferenceSubsystem struct {
	EntityType                              entity_type.IService
	entityTypeRepository                    entity_type.Repository
	EntityType2Property                     entity_type2property.IService
	entityType2PropertyRepository           entity_type2property.Repository
	Property                                property.IService
	propertyRepository                      property.Repository
	Relation                                entity_type.RelationService
	relationRepository                      entity_type.RelationRepository
	PropertyGroup                           property_group.IService
	propertyGroupRepository                 property_group.Repository
	PropertyType                            property_type.IService
	propertyTypeRepository                  property_type.Repository
	propertyType2PropertyViewTypeRepository property_type2property_view_type.Repository
	PropertyUnit                            property_unit.IService
	propertyUnitRepository                  property_unit.Repository
	PropertyViewType                        property_view_type.IService
	propertyViewTypeRepository              property_view_type.Repository
	TextSource                              text_source.IService
	textSourceRepository                    text_source.Repository
	TextValue                               text_value.IService
	textValueRepository                     text_value.Repository
	TextLang                                text_lang.IService
	textLangRepository                      text_lang.Repository
	EntityIDRepository                      EntityIDRepository
}

func NewReferenceSubsystem(infra *infrastructure.Infrastructure, cfg *config.Configuration) (*ReferenceSubsystem, error) {
	s := &ReferenceSubsystem{}
	s.EntityIDRepository = gorm.GetEntityIDRepository(infra.Logger, infra.ReferenceDB, cfg.Infrastructure.DataSharding.GetEntityTypesByClusterSysnames())

	if err := s.setupRepositories(infra); err != nil {
		return nil, err
	}
	s.setupServices(infra.Logger)

	if err := s.autoMigrate(infra.ReferenceDB); err != nil {
		return nil, err
	}

	if err := s.dbDataInit(cfg.Metadata); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *ReferenceSubsystem) autoMigrate(db minipkg_gorm.IDB) error {
	if db.IsAutoMigrate() {
		s.EntityIDRepository.AutoMigrate()

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
	s.textSourceRepository, ok = repo.(text_source.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", text_source.EntityName, text_source.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, property_type.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_type.EntityName, err)
	}
	s.propertyTypeRepository, ok = repo.(property_type.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_type.EntityName, property_type.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, property_view_type.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_view_type.EntityName, err)
	}
	s.propertyViewTypeRepository, ok = repo.(property_view_type.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_view_type.EntityName, property_view_type.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, entity_type.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", entity_type.EntityName, err)
	}
	s.entityTypeRepository, ok = repo.(entity_type.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", entity_type.EntityName, entity_type.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, property_unit.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_unit.EntityName, err)
	}
	s.propertyUnitRepository, ok = repo.(property_unit.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_unit.EntityName, property_unit.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, property_group.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_group.EntityName, err)
	}
	s.propertyGroupRepository, ok = repo.(property_group.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_group.EntityName, property_group.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, property_type2property_view_type.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_type2property_view_type.EntityName, err)
	}
	s.propertyType2PropertyViewTypeRepository, ok = repo.(property_type2property_view_type.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_type2property_view_type.EntityName, property_type2property_view_type.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, property.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property.EntityName, err)
	}
	s.propertyRepository, ok = repo.(property.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property.EntityName, property.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, text_value.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", text_value.EntityName, err)
	}
	s.textValueRepository, ok = repo.(text_value.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", text_value.EntityName, text_value.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, entity_type2property.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", entity_type2property.EntityName, err)
	}
	s.entityType2PropertyRepository, ok = repo.(entity_type2property.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", entity_type2property.EntityName, entity_type2property.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, entity_type.RelationEntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", entity_type.RelationEntityName, err)
	}
	s.relationRepository, ok = repo.(entity_type.RelationRepository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", entity_type.RelationEntityName, entity_type.RelationEntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, text_lang.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", text_lang.EntityName, err)
	}
	s.textLangRepository, ok = repo.(text_lang.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", text_lang.EntityName, text_lang.EntityName, repo)
	}

	return nil
}

func (s *ReferenceSubsystem) setupServices(logger log.ILogger) {
	s.TextLang = text_lang.NewService(logger, s.textLangRepository)
	s.TextSource = text_source.NewService(logger, s.textSourceRepository)
	s.TextValue = text_value.NewService(logger, s.textValueRepository)
	s.EntityType2Property = entity_type2property.NewService(logger, s.entityType2PropertyRepository)
	s.Relation = entity_type.NewRelationService(logger, s.relationRepository)
	s.PropertyGroup = property_group.NewService(logger, s.propertyGroupRepository)
	s.PropertyType = property_type.NewService(logger, s.propertyTypeRepository, s.propertyType2PropertyViewTypeRepository, s.TextLang)
	s.PropertyUnit = property_unit.NewService(logger, s.propertyUnitRepository, s.TextLang)
	s.PropertyViewType = property_view_type.NewService(logger, s.propertyViewTypeRepository)
	s.Property = property.NewService(logger, s.propertyRepository, s.PropertyType, s.PropertyUnit, s.PropertyViewType, s.PropertyGroup, s.TextLang)
	s.EntityType = entity_type.NewService(logger, s.entityTypeRepository, s.Relation, s.Property, s.TextLang)
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

func (s *ReferenceSubsystem) dbDataInit(metadata *config.Metadata) error {
	ctx := context.Background()

	err := s.TextLang.DataInit(ctx, metadata.Languages)
	if err != nil {
		return err
	}

	err = s.PropertyType.DataInit(ctx)
	if err != nil {
		return err
	}

	err = s.PropertyUnit.DataInit(ctx, metadata.PropertyUnits)
	if err != nil {
		return err
	}

	err = s.EntityType.DataInit(ctx, metadata.EntityTypes)
	if err != nil {
		return err
	}

	return nil
}

func (s *ReferenceSubsystem) LangDataInit(ctx context.Context, langsConfig config.Languages) error {
	count, err := s.TextLang.Count(ctx, &selection_condition.SelectionCondition{})
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
		item.Cfgname = langConfig.Cfgname
		s.TextLang.Create(ctx, item)
	}

	return nil
}

func (s *ReferenceSubsystem) PropertyTypeDataInit(ctx context.Context) error {
	count, err := s.PropertyType.Count(ctx, &selection_condition.SelectionCondition{})
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
		err = s.PropertyType.Create(ctx, &i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ReferenceSubsystem) PropertyUnitInit(ctx context.Context, unitsConfig config.PropertyUnits) error {
	count, err := s.PropertyUnit.Count(ctx, &selection_condition.SelectionCondition{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	langsSl, err := s.TextLang.GetCodesEmptyInterfaceSlice(ctx)
	if err != nil {
		return err
	}

	langsIDsMap, err := s.TextLang.GetMapCodeID(ctx)
	if err != nil {
		return err
	}

	for sysname, unitConfig := range unitsConfig {
		unit := property_unit.New()
		unit.Sysname = sysname
		if err := s.PropertyUnit.TCreate(ctx, unit, 1); err != nil {
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

			name := texts.Name
			description := texts.Description
			unit.Name = &name
			unit.Description = &description
			if err := s.PropertyUnit.TUpdate(ctx, unit, langID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *ReferenceSubsystem) EntityTypeInit(ctx context.Context, EntityTypesConfig config.EntityTypes) error {
	count, err := s.EntityType.Count(ctx, &selection_condition.SelectionCondition{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	propertiesCount, err := s.Property.Count(ctx, &selection_condition.SelectionCondition{})
	if err != nil {
		return err
	}

	langsSl, err := s.TextLang.GetCodesEmptyInterfaceSlice(ctx)
	if err != nil {
		return err
	}

	langsIDsMap, err := s.TextLang.GetMapCodeID(ctx)
	if err != nil {
		return err
	}

	for sysname, entityTypeConfig := range EntityTypesConfig {
		entityType := entity_type.New()
		entityType.Sysname = sysname
		if err := s.EntityType.TCreate(ctx, entityType, 1); err != nil {
			return err
		}

		for lang, texts := range entityTypeConfig.Texts {
			if err := validation.Validate(lang, validation.In(langsSl...)); err != nil {
				return errors.Wrapf(err, "EntityTypeInit error: invalid lang = %q", lang)
			}
			langID, ok := langsIDsMap[lang]
			if !ok {
				return errors.Errorf("EntityTypeInit error: not found lang = %q", lang)
			}

			name := texts.Name
			description := texts.Description
			entityType.Name = &name
			entityType.Description = &description
			if err := s.EntityType.TUpdate(ctx, entityType, langID); err != nil {
				return err
			}
		}
		if propertiesCount == 0 {
			err := s.PropertyInit(ctx, entityType, entityTypeConfig.Properties)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *ReferenceSubsystem) PropertyInit(ctx context.Context, entityType *entity_type.EntityType, PropertiesConfig config.Properties) error {
	count, err := s.Property.Count(ctx, &selection_condition.SelectionCondition{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	langsSl, err := s.TextLang.GetCodesEmptyInterfaceSlice(ctx)
	if err != nil {
		return err
	}

	langsIDsMap, err := s.TextLang.GetMapCodeID(ctx)
	if err != nil {
		return err
	}

	for sysname, propertyConfig := range PropertiesConfig {
		prop := property.New()
		prop.Sysname = entityType.Sysname + "." + sysname

		if propertyConfig.PropertyType != "" {
			if prop.PropertyTypeID, err = s.PropertyType.GetIDBySysname(ctx, propertyConfig.PropertyType); err != nil {
				return err
			}
		}

		if propertyConfig.PropertyUnit != "" {
			propertyUnitID, err := s.PropertyUnit.GetIDBySysname(ctx, propertyConfig.PropertyUnit)
			if err != nil {
				return err
			}
			prop.PropertyUnitID = &propertyUnitID
		}

		if propertyConfig.PropertyViewType != "" {
			propertyViewTypeID, err := s.PropertyViewType.GetIDBySysname(ctx, propertyConfig.PropertyViewType)
			if err != nil {
				return err
			}
			prop.PropertyViewTypeID = &propertyViewTypeID
		}

		if propertyConfig.PropertyGroup != "" {
			propertyGroupID, err := s.PropertyGroup.GetIDBySysname(ctx, propertyConfig.PropertyGroup)
			if err != nil {
				return err
			}
			prop.PropertyGroupID = &propertyGroupID
		}

		prop.IsSpecific = propertyConfig.IsSpecific
		prop.IsRange = propertyConfig.IsRange
		prop.IsMultiple = propertyConfig.IsMultiple
		prop.SortOrder = propertyConfig.SortOrder
		prop.Options = propertyConfig.Options
		if err := s.Property.TCreate(ctx, prop, 1); err != nil {
			return err
		}

		if err := s.EntityType.BindProperty(ctx, entityType.ID, prop.ID); err != nil {
			return err
		}

		for lang, texts := range propertyConfig.Texts {
			if err := validation.Validate(lang, validation.In(langsSl...)); err != nil {
				return errors.Wrapf(err, "PropertyInit error: invalid lang = %q", lang)
			}
			langID, ok := langsIDsMap[lang]
			if !ok {
				return errors.Errorf("PropertyInit error: not found lang = %q", lang)
			}

			name := texts.Name
			description := texts.Description
			prop.Name = &name
			prop.Description = &description
			if err := s.Property.TUpdate(ctx, prop, langID); err != nil {
				return err
			}
		}
	}

	return nil
}
