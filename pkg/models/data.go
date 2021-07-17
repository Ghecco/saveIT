package models

type User struct {
	ID       uint `gorm:"primaryKey;autoIncrement:true"`
	Name     string
	Password string
}

type Idea struct {
	ID      uint `gorm:"primaryKey;autoIncrement:true"`
	Content string
	User    User `gorm:"foreignKey:ID"`
}
