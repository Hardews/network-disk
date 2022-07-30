package model

type User struct {
	Uid      int    `gorm:"primaryKey;AUTO_INCREMENT=1;not null"`
	Username string `gorm:"not null;unique;type:varchar(20)"`
	Password string `gorm:"type:varchar(100)"`
}

type UserResources struct {
	Folder       string
	Filename     string
	ResourceName string
	Permission   string
	CreateAt     string
}
