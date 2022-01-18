package config

type Metadata struct {
	IsAutomigrate bool
	Languages     Languages
	PropertyUnits PropertyUnits
	EntityTypes   EntityTypes
	Relations     Relations
}

type Languages []Language

type Language struct {
	Code    string
	Name    string
	Cfgname string
}

type PropertyUnits []PropertyUnit

type PropertyUnit struct {
	Sysname string
	Texts   map[string]NameAndDescriptionText
}

type NameAndDescriptionText struct {
	Name        string
	Description string
}

type EntityTypes []EntityType

type EntityType struct {
	Sysname    string
	Texts      map[string]NameAndDescriptionText
	Properties Properties
}

type Properties []Property

type Property struct {
	Sysname          string
	PropertyType     string
	PropertyUnit     string
	PropertyViewType string
	PropertyGroup    string
	IsSpecific       bool
	IsRange          bool
	IsMultiple       bool
	SortOrder        uint
	Options          []map[string]interface{}
	Texts            map[string]NameAndDescriptionText
}

type Relations []Relation

type Relation struct {
	Sysname              string
	PropertyViewType     string
	PropertyGroup        string
	IsSpecific           bool
	IsMultiple           bool
	SortOrder            uint
	Options              []map[string]interface{}
	Texts                map[string]NameAndDescriptionText
	UndependedEntityType string
	DependedEntityType   string
}
