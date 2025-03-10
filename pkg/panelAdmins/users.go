package panelAdmins

import "time"

type Personal struct {
	FirstName string
	LastName  string
	Gender    string
	Email     string
	Phone     string
	Dob       time.Duration
}

type User struct {
	ID            string
	Personal      Personal
	Role          Role
	Team          Team
	LastModified  time.Time
	CreatedAt     time.Time
	ArchiveStatus bool
	DeletedStatus bool
	CreatedBy     string
	ModifiedBy    string
}
