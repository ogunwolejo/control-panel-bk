package panelAdmins

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"time"
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
			Name: "James",
			Role: Role{
				Name:          "Admin",
				createdAt:     time.Date(2021, time.April, 5, 3, 0, 0, 1, time.UTC),
				CreatedBy:     "123",
				UpdatedBy:     "123",
				ArchiveStatus: false,
				Description:   "",
				ID:            "role 1",
				lastModified:  time.Date(2021, time.April, 5, 3, 30, 15, 2, time.UTC),
				Permission:    prm,
			},
			ID:      "Mem 1",
			Email:   "member1@gmail.com",
			Phone:   "09031846448",
			Gender:  "Male",
			IsLead:  false,
			Profile: "pic",
		},
		{
			Name: "Peter",
			Role: Role{
				Name:          "Sales",
				createdAt:     time.Date(2021, time.November, 5, 3, 0, 0, 1, time.UTC),
				CreatedBy:     "123",
				UpdatedBy:     "123",
				ArchiveStatus: false,
				Description:   "",
				ID:            "role 2",
				lastModified:  time.Date(2021, time.November, 5, 3, 30, 15, 2, time.UTC),
				Permission:    prm,
			},
			ID:      "Mem 2",
			Email:   "member2@gmail.com",
			Phone:   "09031846447",
			Gender:  "Male",
			IsLead:  true,
			Profile: "pic",
		},
		{
			Name: "Lizzy",
			Role: Role{
				Name:          "HR",
				createdAt:     time.Date(2021, time.May, 5, 3, 0, 0, 1, time.UTC),
				CreatedBy:     "123",
				UpdatedBy:     "123",
				ArchiveStatus: false,
				Description:   "",
				ID:            "role 3",
				lastModified:  time.Date(2021, time.May, 5, 3, 30, 15, 2, time.UTC),
				Permission:    prm,
			},
			ID:      "Mem 3",
			Email:   "member3@gmail.com",
			Phone:   "09031846448",
			Gender:  "Male",
			IsLead:  false,
			Profile: "pic",
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
					Profile: "new profile 1",
					IsLead:  false,
					Gender:  "female",
					Phone:   "0810678908",
					ID:      "new 123",
					Name:    "Younger1",
					Email:   "newmem1@gmail.com",
				},
				{
					Profile: "new profile 2",
					IsLead:  false,
					Gender:  "female",
					Phone:   "0810678908",
					ID:      "new 123",
					Name:    "Younger2",
					Email:   "newmem2@gmail.com",
				},
				{
					Profile: "new profile 3",
					IsLead:  false,
					Gender:  "female",
					Phone:   "0810678908",
					ID:      "new 123",
					Name:    "Younger3",
					Email:   "newmem3@gmail.com",
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
					Profile: "new profile 1",
					IsLead:  false,
					Gender:  "female",
					Phone:   "0810678908",
					ID:      "new 123",
					Name:    "Younger1",
					Email:   "newmem1@gmail.com",
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

func TestTeam_ChangeTeamLead(t *testing.T) {}

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
