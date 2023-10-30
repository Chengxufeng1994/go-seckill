package entity

type Activity struct {
	ActivityId   int     `gorm:"primaryKey,not null"`
	ActivityName string  `gorm:"size:50,not null"`
	ProductId    int     `gorm:"not null"`
	StartTime    int64   `gorm:"not null,default:0"`
	EndTime      int64   `gorm:"not null,default:0"`
	Total        int     `gorm:"not null,default:0"`
	Status       int     `gorm:"not null,default:0"`
	SecSpeed     int     `gorm:"not null,default:0"`
	BuyLimit     int     `gorm:"not null,default:0"`
	BuyRate      float64 `gorm:"not null,default:0.00"`
}
