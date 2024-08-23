package database

import (
	"gorm.io/gorm"
)

// UsersEntity represents users of the marina.
type UsersEntity struct {
	gorm.Model
	Email       string         `json:"email" gorm:"unique"`
	Password    string         `json:"-"` // Don't expose the password in JSON responses
	LastName    string         `json:"last_name"`
	FirstName   string         `json:"first_name"`
	Phone       string         `json:"phone"`
	Address     string         `json:"address"`
	NIP         string         `json:"nip"`
	CompanyName string         `json:"company_name"`
	Notes       string         `json:"notes"`
	Boats       []*BoatsEntity `json:"boats" gorm:"many2many:user_boats;"`
}

func (UsersEntity) TableName() string {
	return "users"
}

// SlipsEntity represents parking spots for boats in the marina.
type SlipsEntity struct {
	gorm.Model
	Number     int            `json:"number" gorm:"unique"`
	IsOccupied bool           `json:"is_occupied"`
	Notes      string         `json:"notes"`
	Boats      []*BoatsEntity `json:"boats" gorm:"many2many:slip_boats;"`
}

func (SlipsEntity) TableName() string {
	return "slips"
}

// BoatsEntity represents boats docked at the marina.
type BoatsEntity struct {
	gorm.Model
	Name   string         `json:"name"`
	Type   string         `json:"type"`
	Length string         `json:"length"`
	Width  string         `json:"width"`
	Weight string         `json:"weight"`
	Draft  string         `json:"draft"`
	Owners []*UsersEntity `json:"owners" gorm:"many2many:user_boats;"`
	Slips  []*SlipsEntity `json:"slips" gorm:"many2many:slip_boats;"`
	Notes  string         `json:"notes"`
}

func (BoatsEntity) TableName() string {
	return "boats"
}
