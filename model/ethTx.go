package model

import (
	"go-work-01/utils"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type EthTx struct {
	gorm.Model
	TxHash      string          `gorm:"column:tx_hash;type:varchar(66);not null;uniqueIndex:idx_tx_hash"`
	BlockNumber int64           `gorm:"column:block_number;type:bigint;not null"`
	FromAddress string          `gorm:"column:from_address;type:varchar(42);not null"`
	ToAddress   string          `gorm:"column:to_address;type:varchar(42)"`
	Value       decimal.Decimal `gorm:"column:value;type:DECIMAL(38,18);not null"`
	GasPrice    int64           `gorm:"column:gas_price;type:bigint;not null"`
	GasLimit    int64           `gorm:"column:gas_limit;type:bigint;not null"`
	Gas         int64           `gorm:"column:gas;type:bigint;not null"`
	InputData   string          `gorm:"column:input_data;type:text"`
}

func CreateTx(tx *EthTx) error {
	return utils.Db.Create(tx).Error
}

func GetTxByHash(hash string) (tx *EthTx, e error) {
	e = utils.Db.Where("tx_hash = ?", hash).First(&tx).Error
	return
}
