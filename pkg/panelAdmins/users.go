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
	Profile   string    `json:"profile,omitempty"`
}

type User struct {
	ID        string `json:"id,omitempty"`
	Personal  `json:"personal"`
	Role      Role      `json:"role"`
	Team      []Team    `json:"team,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	//ArchiveStatus bool      `json:"archive_status"`
	DeletedStatus bool   `json:"is_deleted_status"`
	CreatedBy     string `json:"created_by"`
	UpdatedBy     string `json:"updated_by"`
}
