package model

import (
	"time"
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type Account struct {
	Base `valid:"required"`
	OwnerName string `json:"owner_name" valid:"notnull"`
	Bank *Bank `valid:"-"`
	Number string `json:"number" valid:"notnull"`
}

func (account * Account) isValid() error {
	_, err := govalidator.ValidateStruct(account)
	if err != nil {
		return err
	}
	return nil
}

func NewAccount(bank *Bank, number string, ownerName string) (*Account, error) {
	account := Account{
		Bank: bank,
		Number: number,
		OwnerName: ownerName,
	}

	account.ID = uuid.NewV4().String()
	account.CreatedAt = time.Now()

	return &account, nil
}