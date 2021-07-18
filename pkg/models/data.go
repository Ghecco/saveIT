package models

type User struct {
	ID       uint64 `gorm:"primaryKey"`
	Name     string
	Password string
}

type Idea struct {
	ID      uint64 `gorm:"primaryKey"`
	Content string
	UserID  uint64
}
