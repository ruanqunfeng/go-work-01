package model

import (
	"go-work-01/utils"
	"strconv"

	"github.com/ethereum/go-ethereum/rpc"
	"gorm.io/gorm"
)

type EthBlockTag struct {
	gorm.Model
	TagName     string `gorm:"column:tag_name;type:varchar(20);not null"`
	BlockNumber int64  `gorm:"column:block_number;type:bigint;not null"`
}

func GetBlock(s string) (*EthBlock, error) {
	number, err := ParseBLockNumber(s)
	if err != nil {
		return nil, err
	}
	return GetBlockByNumber(number)
}

func InsertBlockTag(tagName string, blockNumber int64) error {
	block := EthBlockTag{TagName: tagName, BlockNumber: blockNumber}
	tx := utils.Db.Create(block)
	return tx.Error
}

func ParseBLockNumber(s string) (int64, error) {
	var err error
	var number int64
	switch s {
	case "head":
		number = int64(rpc.LatestBlockNumber)
	case "finalized":
		number = int64(rpc.FinalizedBlockNumber)
	case "safe":
		number = int64(rpc.SafeBlockNumber)
	default:
		number, err = strconv.ParseInt(s, 10, 64)
	}

	return number, err
}
