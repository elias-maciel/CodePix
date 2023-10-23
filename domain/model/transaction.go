package model

import (
	"time"
	"errors"
	"github.com/asaskevich/govalidator"
)

const (
	TransactionPending = "pending"
	TransactionCompleted = "completed"
	TransactionConfirmed = "confirmed"
	TransactionError = "error"
)

type TransactionsRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction, error)
}

type Transactions struct {

}

type Transaction struct {
	Base `valid:"required"`
	AccountFrom *Account `valid:"-"`
	AccountFromID string `gorm:"column:account_from_id;type:uuid;not null" vaid:"notnull"`
	Amount float64 `json:"amount" gorm:"type:float" valid:"notnull"`
	PixKeyTo *PixKey `valid:"-"`
	PixKeyIDTo string `gorm:"column:pix_key_id_to;type:uuid;not null" vaid:"notnull"`
	Status string `json:"status" gorm:"type:varchar(20)" valid:"notnull"`
	Description string `json:"description" gorm:"type:varchar(255)" valid:"notnull"`
	CancelDescription string `json:"cancel_description" gorm:"type:varchar(255)" valid:"-"`
}

func (transaction *Transaction) isValid() error {
	_, err := govalidator.ValidateStruct(transaction)

	if transaction.Amount <= 0 {
		return errors.New("the amount must be greater than 0")
	} 

	if transaction.Status != TransactionPending && transaction.Status != TransactionCompleted && transaction.Status != TransactionConfirmed && transaction.Status != TransactionError {
		return errors.New("invalid status for the transaction")
	}

	if transaction.PixKeyTo.AccountID == transaction.AccountFrom.ID {
		return errors.New("the source and the designation account cannot be the same")
	}

	if err != nil {
		return err
	}
	return nil
}

func NewTransaction(accountFrom *Account, amount float64, pixKeyTo *PixKey, description string) (*Transaction, error) {
	transaction := Transaction{
		AccountFrom: accountFrom,
		Amount: amount,
		PixKeyTo: pixKeyTo,
		Status: TransactionPending,
		Description: description,
	}

	if err := transaction.isValid(); err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (transaction *Transaction) Complete() error {
	transaction.Status = TransactionCompleted
	transaction.UpdatedAt = time.Now()
	err := transaction.isValid()
	return err
}

func (transaction *Transaction) Confirm() error {
	transaction.Status = TransactionConfirmed
	transaction.UpdatedAt = time.Now()
	err := transaction.isValid()
	return err
}

func (transation *Transaction) Cancel(description string) error {
	transation.Status = TransactionError
	transation.UpdatedAt = time.Now()
	transation.CancelDescription = description
	err := transation.isValid()
	return err
}