package panelAdmins

import (
	"errors"
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
	// We must ensure the mate is a member  first
	//n := slices.CompareFunc(currentLead, t.TeamLead, func(e1 TeamMate, e2 TeamMate) int {
	//	return strings.Compare(e1.ID, e2.ID)
	//})
	//
	//if t.IsMember(mate) && n == 0 {
	//	return nil
	//}
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
