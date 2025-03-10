package panelAdmins

import (
	"errors"
	"log"
	"time"
)

type TeamMate struct {
	ID      string
	Name    string
	Role    Role
	Gender  string
	Email   string
	Phone   string
	Profile string
	IsLead  bool
}

type Team struct {
	ID            string
	TeamLead      TeamMate
	TeamMember    []TeamMate
	LastModified  time.Time
	CreatedAt     time.Time
	CreatedBy     string
	ModifiedBy    string
	ArchiveStatus bool
	DeletedStatus bool
}

func (t *Team) AddNewTeamMember(member []TeamMate) error {
	tm := append(t.TeamMember, member...)
	t.TeamMember = tm
	return nil
}

func (t *Team) RemoveTeamMember(mate []TeamMate) error {
	// We will use waitGroups
	log.Println("Team Member B4 Removing selected ID : ", t.TeamMember)
	// We will use the slice methods here
	return nil
}

func (t *Team) ChangeTeamLead(mate TeamMate) error {
	// We must ensure the mate is a member  first
	return nil
}

func (t *Team) IsMember(mate TeamMate) bool {
	//n, isFound := slices.BinarySearch(t.TeamMember, mate)
	return true
}

func (t *Team) ArchiveTeam() error {
	if t.ArchiveStatus {
		return errors.New("team has already been archived")
	}

	t.ArchiveStatus = true
	return nil
}

func (t *Team) UnArchiveTeam() error {
	if !t.ArchiveStatus {
		return errors.New("team is not in the archive catalogue")
	}

	t.ArchiveStatus = false
	return nil
}

func (t *Team) DeleteTeam() error {
	if t.DeletedStatus {
		return errors.New("team has already been deleted")
	}

	t.DeletedStatus = true
	return nil
}

func (t *Team) UpdateTeamDataInDB() {}

func GetTeams() []Team {
	return []Team{}
}

func GetTeam(id string) Team {
	return Team{}
}
