package model

import (
	"go-work-01/utils"

	"gorm.io/gorm"

	"gorm.io/datatypes"
)

type EthBlock struct {
	gorm.Model
	BlockNumber      int64          `gorm:"column:block_number;type:bigint;not null;uniqueIndex:idx_block_number"`
	BlockHash        string         `gorm:"column:block_hash;type:varchar(66);not null;uniqueIndex:idx_block_hash"`
	ParentHash       string         `gorm:"column:parent_hash;type:varchar(66);not null"`
	MinerAddress     string         `gorm:"column:miner_address;type:varchar(42);not null"`
	GasUsed          int64          `gorm:"column:gas_used;type:bigint;not null"`
	GasLimit         int64          `gorm:"column:gas_limit;type:bigint;not null"`
	TimeStamp        int64          `gorm:"column:timestamp;type:bigint;not null"`
	TransactionCount int            `gorm:"column:transaction_count;type:int;not null"`
	BlockData        datatypes.JSON `gorm:"column:block_data;type:json;not null"`
}

func GetLastBlock() (*EthBlock, error) {
	var block EthBlock
	err := utils.Db.Order("block_number desc").First(&block).Error
	if err != nil {
		return nil, err
	}
	return &block, nil

}

func GetBlockByNumber(number int64) (*EthBlock, error) {
	var block EthBlock
	err := utils.Db.Where("block_number = ?", number).First(&block).Error
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func GetBlockByHash(hash string) (*EthBlock, error) {
	var block EthBlock
	err := utils.Db.Where("block_hash = ?", hash).First(&block).Error
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func InsertBlock(e *EthBlock) error {
	err := utils.Db.Create(e).Error
	if err != nil {
		return err
	}
	return nil
}
