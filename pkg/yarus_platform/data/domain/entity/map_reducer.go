package entity

import (
	minipkg_gorm "github.com/minipkg/db/gorm"
)

type MapReducer interface {
	GetDBByID(ID uint, typeID uint) minipkg_gorm.IDB
	GetDBForNew(typeID uint) minipkg_gorm.IDB
}
