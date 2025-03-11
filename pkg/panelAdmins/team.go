package panelAdmins

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"
	"sync"
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
	t.SortTeamMembers()
	return nil
}

func (t *Team) RemoveTeamMember(mates []TeamMate) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, mate := range mates {
		wg.Add(1)

		go func(mate TeamMate) {
			defer wg.Done()

			mutex.Lock() // we lock the state, so that while we are performing some action nothing can be added or removed from the team members
			defer mutex.Unlock()

			for i, member := range t.TeamMember {
				if member.ID == mate.ID && !member.IsLead {
					t.TeamMember = append(t.TeamMember[:i], t.TeamMember[i+1:]...)
					break
				}
			}

		}(mate)

		wg.Wait()
	}
}

func (t *Team) ChangeTeamLead(mate TeamMate, currentLead TeamMate) error {
	if t.TeamLead.ID != currentLead.ID {
		return errors.New(fmt.Sprintf("the current lead of id %s is not accurate with the lead id sent", currentLead.ID))
	}

	if !t.IsMember(mate) {
		nm := make([]TeamMate, 1)
		nm = append(nm, mate)
		if err := t.AddNewTeamMember(nm); err != nil {
			return err // errors.New("unable to add proposed team lead to the team as the user is not a team member")
		}
	}

	// Change the Team Member Status
	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, member := range t.TeamMember {
		wg.Add(1)

		go func(mt TeamMate) {
			defer wg.Done()

			mutex.Lock()
			defer mutex.Unlock()

			// We set previous team lead to false
			if member.ID != mt.ID && member.IsLead {
				member.IsLead = false
			}

			// We change the team lead
			if member.ID == mt.ID && !member.IsLead {
				member.IsLead = true
				t.TeamLead = member
			}

		}(mate)
		wg.Wait()
	}

	return nil
}

func (t *Team) SortTeamMembers() {
	slices.SortStableFunc(t.TeamMember, func(a, b TeamMate) int {
		return strings.Compare(a.ID, b.ID)
	})
	log.Println(t.TeamMember)
}

func (t *Team) IsMember(mate TeamMate) bool {
	if isSorted := slices.IsSortedFunc(t.TeamMember, func(a, b TeamMate) int {
		return strings.Compare(a.ID, b.ID)
	}); !isSorted {
		t.SortTeamMembers()
	}

	_, isFound := slices.BinarySearchFunc(t.TeamMember, mate, func(a, b TeamMate) int {
		return strings.Compare(a.ID, b.ID)
	})

	return isFound
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
