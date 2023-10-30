package entity

type Product struct {
	ProductId   int    `gorm:"primaryKey,not null"`
	ProductName string `gorm:"size:50,not null"`
	Total       int    `gorm:"not null,default:0"`
	Status      int    `gorm:"not null,default:0"`
}
