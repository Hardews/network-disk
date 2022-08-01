package model

type User struct {
	Uid      int    `gorm:"primaryKey;AUTO_INCREMENT=1;not null"`
	Username string `gorm:"not null;unique;type:varchar(20)"`
	Password string `gorm:"type:varchar(100)"`
}

type AdminUser struct {
	Uid      int    `gorm:"primaryKey;AUTO_INCREMENT=1;not null"`
	Username string `gorm:"not null;unique;type:varchar(20)"`
}

type Url struct {
	Uid int    `gorm:"primaryKey;AUTO_INCREMENT=1;not null"`
	Url string `gorm:"not null;unique;type:varchar(200)"`
}

type UserResources struct {
	Folder       string
	Path         string
	Filename     string
	ResourceName string
	Permission   string
	DownloadAddr string
	CreateAt     string
}
