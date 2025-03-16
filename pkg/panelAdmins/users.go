package panelAdmins

import (
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
	ID string `json:"id;omitempty"`
	Personal
	Role          Role   `json:"role"`
	Team          []Team `json:"team;omitempty"`
	LastModified  int64  `json:"lastModified;omitempty"`
	CreatedAt     int64  `json:"createdAt;omitempty"`
	ArchiveStatus bool   `json:"archiveStatus;omitempty"`
	DeletedStatus bool   `json:"deletedStatus;omitempty"`
	CreatedBy     string `json:"createdBy"`
	ModifiedBy    string `json:"modifiedBy"`
}
