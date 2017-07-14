package model

type Confirmations struct {
	BlockID int64 `gorm:"primary_key;not null"`
	Good    int32 `gorm:"not null"`
	Bad     int32 `gorm:"not null"`
	Time    int32 `gorm:"not null"`
}

func (c *Confirmations) GetMaxGoodBlock(goodCount int) error {
	return DBConn.Where("good >= ?", goodCount).Last(&c).Error
}

func (c *Confirmations) GetConfirmation(blockID int64) error {
	return DBConn.Where("blockID = ?", blockID).First(&c).Error
}

func (c *Confirmations) Save() error {
	return DBConn.Update(c).Error
}

func (c *Confirmations) Create() error {
	return DBConn.Create(c).Error
}
