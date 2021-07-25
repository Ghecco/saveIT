package models

type User struct {
	ID       uint64 `gorm:"primaryKey"`
	Name     string
	Password string
}

type Idea struct {
	ID      uint64 `gorm:"primaryKey"`
	UserID  uint64
	Title   string
	Content string
	Time    string
}
