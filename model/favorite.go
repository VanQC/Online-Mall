package model

import "gorm.io/gorm"

// Favorite 收藏夹
type Favorite struct {
	gorm.Model
	Member    Member  `gorm:"ForeignKey:MemberID"`
	MemberID  uint    `gorm:"not null"`
	Product   Product `gorm:"ForeignKey:ProductID"`
	ProductID uint    `gorm:"not null"`
	Boss      Member  `gorm:"ForeignKey:BossID"`
	BossID    uint    `gorm:"not null"`
}
