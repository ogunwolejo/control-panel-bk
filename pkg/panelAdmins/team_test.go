package panelAdmins

import (
	"errors"
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

func TestTeam_RemoveTeamMember(t *testing.T) {}

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

func TestTeam_UpdateTeamDataInDB(t *testing.T) {}

func TestTeam_IsMember(t *testing.T) {}

func TestGetTeam(t *testing.T) {}

func TestGetTeams(t *testing.T) {}
