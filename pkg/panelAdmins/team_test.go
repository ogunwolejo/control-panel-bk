package panelAdmins

import (
	"control-panel-bk/util"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

type MockDatabase struct {
	mock.Mock
	db *mongo.Database
}

func (m *MockDatabase) Db(name string) {
	client := &mongo.Client{}
	m.db = client.Database(name)
}

type MockSingleResult struct {
	mock.Mock
}

func (m *MockSingleResult) Decode(v interface{}) error {
	args := m.Called(v)
	if args.Get(0) != nil {
		data, _ := json.Marshal(args.Get(0))
		return json.Unmarshal(data, v)
	}
	return args.Error(1)
}

var prm = Permission{
	Onboarding: ReadWrite{true, true},
	Team:       ReadWrite{true, true},
	Role:       ReadWrite{true, true},
	Billing:    ReadWrite{true, true},
	Tenant:     ReadWrite{true, false},
}

func TestTeam_SortTeamMembers(t *testing.T) {
	type TestKit struct {
		Title           string
		Team            Team
		ExpectedMembers []string
	}

	table := []TestKit{
		{
			Title: "Testing sorted team members and result is equals to the expected outcome",
			Team: Team{
				TeamMember: []string{
					"Mary",
					"Cyan",
					"Lizzy",
				},
			},
			ExpectedMembers: []string{
				"Cyan",
				"Lizzy",
				"Mary",
			},
		},
		{
			Title: "Testing sorted team members and result is equals to the expected outcome",
			Team: Team{
				TeamMember: []string{
					"abraham",
					"ruth",
					"daniel",
					"alexander",
				},
			},
			ExpectedMembers: []string{
				"abraham",
				"alexander",
				"daniel",
				"ruth",
			},
		},
	}

	for _, tt := range table {
		tm := tt.Team

		tm.SortTeamMembers()

		n := slices.CompareFunc(tm.TeamMember, tt.ExpectedMembers, func(mate string, mate2 string) int {
			return strings.Compare(mate, mate2)
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

func TestTeam_IsMember(t *testing.T) {
	type Outcome struct {
		Title           string
		Team            Team
		SearchMember    string
		ExpectedOutcome bool
	}

	tst := []Outcome{
		{
			"Test is member should return true",
			Team{
				TeamMember: []string{
					"Josh",
					"lizzy",
					"james",
				},
			},
			"lizzy",
			true,
		},
		{
			"Test is member should return false",
			Team{
				TeamMember: []string{
					"Josh",
					"lizzy",
					"james",
				},
			},
			"bt",
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

func TestGetTeam(t *testing.T) {
	mockDB := new(MockDatabase)
	handler := GetTeam(mockDB.db)

	id := "65a6b1c3e3c2a1b9e4f3d2a1"
	objId, _ := util.GetPrimitiveID(id)

	mockTeam := Team{ID: id, Name: "Alpha"}

	// Mock decode
	singleResult := new(MockSingleResult)
	singleResult.On("Decode", mock.Anything).Return(mockTeam, nil)

	mockDB.On("FindOne", mock.Anything, bson.M{"_id": objId}).Return(singleResult)

	req := httptest.NewRequest("GET", "/teams/65a6b1c3e3c2a1b9e4f3d2a1", nil)
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/teams/{id}", handler)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var response []Team
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Len(t, response, 1)
	require.Equal(t, "Alpha", response[0].Name)
}

func TestGetTeams(t *testing.T) {}
