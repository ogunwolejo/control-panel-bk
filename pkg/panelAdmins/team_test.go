package panelAdmins

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"
)

var prm = Permission{
	Onboarding: ReadWrite{true, true},
	Team:       ReadWrite{true, true},
	Role:       ReadWrite{true, true},
	Billing:    ReadWrite{true, true},
	Tenant:     ReadWrite{true, false},
}

var mockTeam = Team{
	TeamMember: []TeamMate{
		{
			ID: "Mem 1",
		},
		{
			ID: "Mem 2",
		},
		{
			ID: "Mem 3",
		},
	},
}

func TestTeam_AddNewTeamMember(t *testing.T) {
	type Testkit struct {
		Title         string
		Tm            Team
		MemberCount   int
		ExpectedCount int
		newMember     []TeamMate
	}

	table := []Testkit{
		{
			"Add a new user",
			mockTeam,
			len(mockTeam.TeamMember),
			len(mockTeam.TeamMember) + 3,
			[]TeamMate{
				{
					IsLead: false,
					ID:     "new 123",
				},
				{
					IsLead: false,
					ID:     "new 123",
				},
				{
					ID: "new 123",
				},
			},
		},
		{
			"Add new user",
			mockTeam,
			len(mockTeam.TeamMember),
			len(mockTeam.TeamMember) + 1,
			[]TeamMate{
				{
					ID: "new 123",
				},
			},
		},
	}

	for _, tt := range table {
		tm := tt.Tm
		e := tm.AddNewTeamMember(tt.newMember)

		if e != nil {
			t.Errorf("%s: Expected %d but got error %v", tt.Title, tt.ExpectedCount, e)
		}

		if tt.ExpectedCount == len(tm.TeamMember) {
			t.Logf("%s: Expected %d but got %d", tt.Title, tt.ExpectedCount, len(tm.TeamMember))
			t.Logf("%s: Expected %d to be greater than %d", tt.Title, len(tm.TeamMember), tt.MemberCount)
		}

		if tt.ExpectedCount < len(tm.TeamMember) {
			t.Errorf("%s: Expected %d but got %d", tt.Title, tt.ExpectedCount, len(tm.TeamMember))
		}
	}
}

func TestTeam_ArchiveTeam(t *testing.T) {
	type Table struct {
		Name                  string
		team                  Team
		ExpectedArchiveStatus bool
		ExpectedErrMessage    error
	}

	tbl := []Table{
		{
			"Test team archive positively",
			Team{
				ArchiveStatus: false,
			},
			true,
			nil,
		},
		{
			"Test team doesn't archive",
			Team{
				ArchiveStatus: true,
			},
			true,
			errors.New("team has already been archived"),
		},
	}

	for _, tt := range tbl {
		tm := tt.team

		err := tm.ArchiveTeam()

		if err != nil {
			t.Logf("%s: Expect error message -> %v to be equals to -> %v", tt.Name, err, tt.ExpectedErrMessage)
			t.Logf("%s: Expect Status %v to be equals to %v", tt.Name, tt.ExpectedArchiveStatus, tm.ArchiveStatus)
		}

		if err == nil {
			t.Logf("%s: Expect team status to be %v and got %v", tt.Name, tm.ArchiveStatus, tt.ExpectedArchiveStatus)
		}
	}

}

func TestTeam_ChangeTeamLead(t *testing.T) {
	type TestOutCome struct {
		Title           string
		Tm              Team
		CurrentLead     TeamMate
		ProposedLead    TeamMate
		ExpectedOutcome bool
		ExpectedError   error
	}

	tbl := []TestOutCome{
		{
			Title: "Test change team lead",
			Tm: Team{
				TeamLead: TeamMate{
					ID: "alex",
				},
				TeamMember: []TeamMate{
					{
						ID:     "temi",
						IsLead: false,
					},
					{
						ID:     "alex",
						IsLead: true,
					},
				},
			},
			CurrentLead:     TeamMate{ID: "alex"},
			ProposedLead:    TeamMate{ID: "temi"},
			ExpectedOutcome: true,
			ExpectedError:   nil,
		},
		{
			Title: "Test change team lead even though proposed lead is not a team member",
			Tm: Team{
				TeamLead: TeamMate{
					ID:     "grace",
					IsLead: true,
				},
				TeamMember: []TeamMate{
					{
						ID:     "grace",
						IsLead: true,
					},
				},
			},
			CurrentLead:     TeamMate{ID: "grace"},
			ExpectedOutcome: true,
			ExpectedError:   nil,
			ProposedLead: TeamMate{
				ID:     "ife",
				IsLead: false,
			},
		},
		{
			Title: "Test cannot change team lead",
			Tm: Team{
				TeamLead: TeamMate{
					ID:     "lola",
					IsLead: true,
				},
				TeamMember: []TeamMate{
					{
						ID:     "lola",
						IsLead: true,
					},
				},
			},
			CurrentLead:     TeamMate{ID: "lola", IsLead: true},
			ExpectedOutcome: true,
			ExpectedError:   errors.New(fmt.Sprintf("the current lead of id %s is not accurate with the lead id sent", "lola")),
			ProposedLead: TeamMate{
				ID:     "ife",
				IsLead: false,
			},
		},
	}

	for _, tt := range tbl {
		tm := tt.Tm

		// We will check to see if the proposed lead is a team member
		isFound := slices.IndexFunc(tm.TeamMember, func(mate TeamMate) bool {
			return mate.ID == tt.ProposedLead.ID
		})

		err := tm.ChangeTeamLead(tt.ProposedLead, tt.CurrentLead)

		if err != nil && isFound == -1 {
			if len(tm.TeamMember) == len(tt.Tm.TeamMember)+1 {
				t.Errorf("%s: Expected the length of team member of %d, to increase by 1, resulting in %d", tt.Title, len(tt.Tm.TeamMember), len(tm.TeamMember))
			}
		}

		if err != nil && !errors.Is(err, tt.ExpectedError) {
			t.Errorf("%s: Expected error is %s got %s", tt.Title, tt.ExpectedError.Error(), err.Error())
		}

		if err != nil && errors.Is(err, tt.ExpectedError) {
			t.Logf("%s: Expected error is %s got %s", tt.Title, tt.ExpectedError.Error(), err.Error())
		}

		if err == nil && tm.TeamLead.ID != tt.ProposedLead.ID {
			t.Errorf("%s: Expected team lead to be %s, got %s", tt.Title, tt.ProposedLead.ID, tm.TeamLead.ID)
		}
	}
}

func TestTeam_DeleteTeam(t *testing.T) {
	type Table struct {
		Name           string
		team           Team
		ExpectedStatus bool
		ExpectedError  error
	}

	tbl := []Table{
		{
			"Test Delete Team",
			Team{
				DeletedStatus: false,
			},
			true,
			nil,
		},
		{
			"Test unable to delete team",
			Team{
				DeletedStatus: true,
			},
			true,
			errors.New("team has already been deleted"),
		},
	}

	for _, tt := range tbl {
		tm := tt.team

		e := tm.DeleteTeam()
		if e != nil {
			t.Logf("%s: Expect error message -> %v, got %v", tt.Name, tt.ExpectedError, e)
		}

		if e == nil {
			t.Logf("%s: Expected delete status %v and got %v", tt.Name, tt.ExpectedStatus, tm.DeletedStatus)
			t.Logf("%s: Expected error message %v to be equals to %v", tt.Name, tt.ExpectedError, e)
		}

		if tt.ExpectedStatus == tm.DeletedStatus {
			t.Logf("%s: Expected delete status %v and got %v", tt.Name, tt.ExpectedStatus, tm.DeletedStatus)
		}

		if tt.ExpectedStatus != tm.DeletedStatus && e == nil {
			t.Errorf("%s: Expected %v but got %v", tt.Name, tt.ExpectedStatus, tm.DeletedStatus)
		}
	}

}

func TestTeam_RemoveTeamMember(t *testing.T) {
	type TestOutCome struct {
		Title           string
		Tm              Team
		RemoveMembers   []TeamMate
		ExpectedMembers []TeamMate
	}

	tb := []TestOutCome{
		{
			Title: "Test some members to be removed",
			Tm: Team{
				TeamMember: []TeamMate{
					{
						ID:     "james",
						IsLead: false,
					},
					{
						ID:     "john",
						IsLead: false,
					},
					{
						ID:     "isaac",
						IsLead: false,
					},
					{
						ID:     "grace",
						IsLead: false,
					},
				},
			},
			RemoveMembers: []TeamMate{
				{
					ID: "isaac",
				},
				{
					ID: "grace",
				},
			},
			ExpectedMembers: []TeamMate{
				{
					ID: "james",
				},
				{
					ID: "john",
				},
			},
		},

		{
			Title: "Test only non team lead members are removed",
			Tm: Team{
				TeamMember: []TeamMate{
					{
						ID:     "james",
						IsLead: false,
					},
					{
						ID:     "john",
						IsLead: false,
					},
					{
						ID:     "bale",
						IsLead: true,
					},
					{
						ID:     "grace",
						IsLead: false,
					},
				},
			},
			RemoveMembers: []TeamMate{
				{
					ID: "bale",
				},
				{
					ID: "grace",
				},
			},
			ExpectedMembers: []TeamMate{
				{
					ID: "james",
				},
				{
					ID: "john",
				},
				{
					ID: "bale",
				},
			},
		},
	}

	for _, tt := range tb {
		tm := tt.Tm
		tm.RemoveTeamMember(tt.RemoveMembers)

		if len(tt.ExpectedMembers) != len(tm.TeamMember) {
			t.Errorf("%s: The expected length of the team members (%d) is not equal to the outcome %d", tt.Title, len(tt.ExpectedMembers), len(tm.TeamMember))
		}

		n := slices.CompareFunc(tt.ExpectedMembers, tm.TeamMember, func(mate TeamMate, mate2 TeamMate) int {
			return strings.Compare(mate.ID, mate2.ID)
		})

		if n == 0 {
			t.Logf("%s: the expected outcome length is %d, got remaining team members length is %d", tt.Title, len(tt.ExpectedMembers), len(tm.TeamMember))
		} else if n == -1 {
			t.Errorf("%s: the expected outcome in terms of length which is %d is less than result %d", tt.Title, len(tt.ExpectedMembers), len(tm.TeamMember))
		} else if n == 1 {
			t.Errorf("%s: the expected outcome in terms of length which is %d is greater than result %d", tt.Title, len(tt.ExpectedMembers), len(tm.TeamMember))
		}
	}
}

func TestTeam_UnArchiveTeam(t *testing.T) {
	type Table struct {
		Name                  string
		team                  Team
		ExpectedArchiveStatus bool
		ExpectedErrMessage    error
	}

	tbl := []Table{
		{
			"Test team to un-archive",
			Team{
				ArchiveStatus: true,
			},
			false,
			nil,
		},
		{
			"Test team cannot be un-archive",
			Team{
				ArchiveStatus: false,
			},
			false,
			errors.New("team is not in the archive catalogue"),
		},
	}

	for _, tt := range tbl {
		tm := tt.team

		err := tm.UnArchiveTeam()

		if err != nil {
			t.Logf("%s: Expect error message -> %v got %v", tt.Name, tt.ExpectedErrMessage, err)
		}

		if err == nil {
			t.Logf("%s: Expect Archive status to be %v, got %v", tt.Name, tt.ExpectedArchiveStatus, tm.ArchiveStatus)
			t.Logf("%s: Expect Archive status to be %v, got %v", tt.Name, tt.ExpectedArchiveStatus, tm.ArchiveStatus)
		}
	}
}

func TestTeam_SortTeamMembers(t *testing.T) {
	type TestKit struct {
		Title           string
		Team            Team
		ExpectedMembers []TeamMate
	}

	table := []TestKit{
		{
			Title: "Testing sorted team members and result is equals to the expected outcome",
			Team: Team{
				TeamMember: []TeamMate{
					{
						ID: "Mary",
					},
					{
						ID: "Cyan",
					},
					{
						ID: "Lizzy",
					},
				},
			},
			ExpectedMembers: []TeamMate{
				{
					ID: "Cyan",
				},
				{
					ID: "Lizzy",
				},
				{
					ID: "Mary",
				},
			},
		},
		{
			Title: "Testing sorted team members and result is equals to the expected outcome",
			Team: Team{
				TeamMember: []TeamMate{
					{
						ID: "Abraham",
					},
					{
						ID: "Ruth",
					},
					{
						ID: "daniel",
					},
					{
						ID: "Alexander",
					},
				},
			},
			ExpectedMembers: []TeamMate{
				{
					ID: "Abraham",
				},
				{
					ID: "Alexander",
				},
				{
					ID: "Ruth",
				},
				{
					ID: "daniel",
				},
			},
		},
	}

	for _, tt := range table {
		tm := tt.Team

		tm.SortTeamMembers()

		n := slices.CompareFunc(tm.TeamMember, tt.ExpectedMembers, func(mate TeamMate, mate2 TeamMate) int {
			return strings.Compare(mate.ID, mate2.ID)
		})

		if n == 0 {
			t.Logf("%s: The two slice are equal", tt.Title)
		} else if n == -1 {
			t.Logf("%s: The two slice are not equal, hence outcome length %d is less than expected length of %d", tt.Title, len(tm.TeamMember), len(tt.ExpectedMembers))
		} else if n == 1 {
			t.Logf("%s: The two slice are not equal, hence outcome length %d is greater than expected length of %d", tt.Title, len(tm.TeamMember), len(tt.ExpectedMembers))
		}
	}
}

func TestTeam_UpdateTeamDataInDB(t *testing.T) {}

func TestTeam_IsMember(t *testing.T) {
	type Outcome struct {
		Title           string
		Team            Team
		SearchMember    TeamMate
		ExpectedOutcome bool
	}

	tst := []Outcome{
		{
			"Test is member should return true",
			Team{
				TeamMember: []TeamMate{
					{
						ID: "Josh",
					},
					{
						ID: "lizzy",
					},
					{
						ID: "james",
					},
				},
			},
			TeamMate{
				ID: "lizzy",
			},
			true,
		},
		{
			"Test is member should return false",
			Team{
				TeamMember: []TeamMate{
					{
						ID: "Josh",
					},
					{
						ID: "lizzy",
					},
					{
						ID: "james",
					},
				},
			},
			TeamMate{
				ID: "bt",
			},
			false,
		},
	}

	for _, tt := range tst {
		tm := tt.Team

		ok := tm.IsMember(tt.SearchMember)

		if ok == tt.ExpectedOutcome {
			t.Logf("%s: The expected outcome is %v, but got %v", tt.Title, tt.ExpectedOutcome, ok)
		}

		if ok != tt.ExpectedOutcome {
			t.Errorf("%s: The expected outcome is %v, but got %v", tt.Title, tt.ExpectedOutcome, ok)
		}
	}
}

func TestGetTeam(t *testing.T) {}

func TestGetTeams(t *testing.T) {}
