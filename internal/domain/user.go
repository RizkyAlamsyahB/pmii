package domain

// User domain model - Maps to tbl_user (legacy structure)
type User struct {
	ID       uint   `gorm:"column:user_id;primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"column:user_name;type:varchar(100)" json:"name"`
	Email    string `gorm:"column:user_email;type:varchar(60);uniqueIndex" json:"email"`
	Password string `gorm:"column:user_password;type:varchar(255)" json:"-"`             // bcrypt hash
	Level    string `gorm:"column:user_level;type:varchar(10)" json:"level"`             // 1=Admin, 2=Author
	Status   string `gorm:"column:user_status;type:varchar(10);default:1" json:"status"` // 1=Active, 0=Inactive
	Photo    string `gorm:"column:user_photo;type:varchar(40)" json:"-"`                 // Filename only, transform in getter
}

// TableName override nama tabel di database (legacy table name)
func (User) TableName() string {
	return "tbl_user"
}

// GetPhotoURL returns full URL for user photo
func (u *User) GetPhotoURL(baseURL string) string {
	if u.Photo == "" {
		return ""
	}
	return baseURL + "/public/uploads/" + u.Photo
}

// IsAdmin checks if user is admin
func (u *User) IsAdmin() bool {
	return u.Level == "1"
}

// IsActive checks if user is active
func (u *User) IsActive() bool {
	return u.Status == "1"
}
