package entities

import (
	"sailormoon/backend/database"
	"time"

	"gorm.io/gorm"
)

type UsersEntity struct {
	gorm.Model
	Name                  string     `json:"name"`
	Email                 string     `json:"email" gorm:"unique"`
	Password              string     `json:"-"` // Don't expose the password in JSON responses
	LastName              string     `json:"last_name"`
	FirstName             string     `json:"first_name"`
	Phone                 string     `json:"phone"`
	Address               string     `json:"address"`
	NIP                   string     `json:"nip"`
	CompanyName           string     `json:"company_name"`
	ContractStart         time.Time  `json:"contract_start"`
	ContractEnd           time.Time  `json:"contract_end"`
	Notes                 string     `json:"notes"`
	Checked               bool       `json:"checked"`
	CraneOperationDate    *time.Time `json:"crane_operation_date"`
	SecondInstallmentPaid bool       `json:"second_installment_paid"`
	ContractDuration      int        `json:"contract_duration"`
	Remarks               string     `json:"remarks"`
}

func (UsersEntity) TableName() string {
	return "users"
}

func AutoMigrate() error {
	return database.DB.AutoMigrate(&UsersEntity{})
}
