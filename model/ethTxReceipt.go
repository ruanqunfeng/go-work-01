package model

import (
	"go-work-01/utils"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type EthTxReceipt struct {
	gorm.Model
	TxHash          string         `gorm:"column:tx_hash;type:varchar(66);not null;uniqueIndex:idx_tx_hash"`
	Status          int            `gorm:"column:status;type:int;not null"`
	GasUsed         int64          `gorm:"column:gas_used;type:bigint;not null"`
	ContractAddress string         `gorm:"column:contract_address;type:varchar(42)"`
	Logs            datatypes.JSON `gorm:"column:logs;type:json"`
}

func CreateTxReceipt(tx *EthTxReceipt) error {
	return utils.Db.Create(tx).Error
}

func GetTxReceipt(txHash string) (*EthTxReceipt, error) {
	var txReceipt EthTxReceipt
	err := utils.Db.Where("tx_hash = ?", txHash).First(&txReceipt).Error
	if err != nil {
		return nil, err
	}
	return &txReceipt, nil
}
