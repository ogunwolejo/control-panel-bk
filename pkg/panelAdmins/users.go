package panelAdmins

import (
	"errors"
	uuid2 "github.com/gofrs/uuid"
	"gorm.io/gorm"
	"time"
)

type Personal struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Gender    string    `json:"gender"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Dob       time.Time `json:"dob"`
	Profile   string    `json:"profile;omitempty"`
}

type User struct {
	gorm.Model
	ID string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();autoIncrement:false"json:"id;omitempty"`
	Personal
	Role          Role           `json:"role"`
	Team          []Team         `json:"team;omitempty" gorm:"many2many:user_teams;"`
	LastModified  int64          `json:"lastModified;omitempty" gorm:"autoUpdateTime"`
	CreatedAt     int64          `json:"createdAt;omitempty" gorm:"autoCreateTime"`
	ArchiveStatus bool           `json:"archiveStatus;omitempty" gorm:"default:false"`
	DeletedStatus bool           `json:"deletedStatus;omitempty" gorm:"default:false"`
	CreatedBy     string         `json:"createdBy"`
	ModifiedBy    string         `json:"modifiedBy"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (usr *User) BeforeCreate(tx *gorm.DB) (err error) {
	if uid, err := uuid2.NewV4(); err != nil {
		return err
	} else {
		usr.ID = uid.String()
		return nil
	}
}

func (usr *User) MigrateDB(tx *gorm.DB) (err error) {
	if err := tx.AutoMigrate(&User{}); err != nil {
		return err
	}

	return nil
}

func (usr *User) CreateUser(tx *gorm.DB, newUser User) (*User, error) {
	createdNewUser := tx.Create(&newUser)

	if createdNewUser.Error != nil {
		return nil, createdNewUser.Error
	}

	return &newUser, nil
}

func (usr *User) SetUserDeleteStatus() error {
	if !usr.DeletedStatus && (!usr.DeletedAt.Valid || usr.DeletedAt.Valid) {
		usr.DeletedStatus = true
		return nil
	}

	return errors.New("user has been deleted from the database already")
}

func (usr *User) DeleteUser(tx *gorm.DB) error {
	if err := tx.Delete(&User{}, usr.ID).Error; err != nil {
		return err
	}

	if err := usr.SetUserDeleteStatus(); err != nil {
		return err
	}

	if err := tx.Model(&User{}).Update("DeletedStatus", usr.DeletedStatus).Error; err != nil {
		return err
	}

	return nil
}

func GetUsers(tx *gorm.DB, limit int, offset int) (records []User, err error) {
	result := tx.Limit(limit).Offset(offset).Order("name asc").Preload("Team").Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}

	return records, nil
}

func GetUser(tx *gorm.DB, Id string) (record *User, err error) {
	result := tx.First(&record, "id = ?", Id)

	if result.Error != nil {
		return nil, result.Error
	}

	return record, nil
}
