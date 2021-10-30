package config

type Metadata struct {
	IsAutomigrate bool
	Languages     Languages
	PropertyUnits PropertyUnits
	EntityTypes   EntityTypes
}

type Languages map[string]string

type PropertyUnits map[string]PropertyUnit

type PropertyUnit map[string]NameAndDescriptionText

type NameAndDescriptionText struct {
	Name        string
	Description string
}

type EntityTypes map[string]EntityType

type EntityType struct {
	Texts      map[string]NameAndDescriptionText
	Properties Properties
}

type Properties map[string]Property

type Property struct {
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
