package panelAdmins

type TeamMate struct {
	Name   string
	Role   Role
	Gender string
	Email string
	Phone string
	Profile string
	IsLead bool
}

type Team struct {
	ID         string
	TeamLead   TeamMate
	TeamMember []TeamMate
}

func (t *Team) AddNewTeamMember(member TeamMate) {}

func (t *Team) RemoveTeamMember() {}

func (t *Team) ChangeTeamLead() {}

func (t *Team) ArchiveTeam() {}

func (t *Team) UnArchiveTeam() {}

func (t *Team) DeleteTeam() {}

func (t *Team) UpdateTeamDataInDB() {}

func GetTeams() []Team {
	return []Team{}
}

func GetTeam(id string) Team {
	return Team{}
}