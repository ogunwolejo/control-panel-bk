package panelAdmins

import "time"

type Role struct {
	ID            string
	Name          string
	Description   string
	Permission    Permission
	CreatedBy     string
	UpdatedBy     string
	ArchiveStatus bool
	lastModified  time.Time
	createdAt     time.Time
}

type CRole struct {
	Name         string
	Description  string
	Permission   Permission
	CreatedBy    string
	UpdatedBy    string
	lastModified time.Time
	createdAt    time.Time
}

func (rl *Role) UnArchiveRole(updatorId string) {}

func (rl *Role) ArchiveRole(updatorId string) {}

func (rl *Role) DeleteRole(updatorId string) {}

func (rl *Role) UpdateRoleDataToDB() {}

func CreateRole(crl CRole) Role {
	return Role{}
}

func FetchRoles(startNo int, endNo int, skip int) []Role {
	return []Role{}
}

func FetchRole(id string) Role {
	return Role{}
}
