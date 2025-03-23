package panelAdmins

import (
	"bytes"
	"context"
	"control-panel-bk/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"
)

const (
	testDBURI = "mongodb://localhost:27017/test_db"
	NoOfSeed  = 15
)

func generateRandomTeamSeeds(db *mongo.Database) []Team {
	collection := db.Collection("teams")

	rand.Seed(time.Now().UnixNano()) // Seed random number generator

	var teams []Team

	// Generate 25 random team records
	for i := 0; i < NoOfSeed; i++ {
		team := Team{
			Name:          faker.Name(),
			Description:   faker.Sentence(),
			TeamLead:      faker.Name(),
			TeamMember:    []string{faker.Name(), faker.Name(), faker.Name()},
			UpdatedAt:     time.Now().Add(time.Duration(rand.Intn(100)) * time.Hour),
			CreatedAt:     time.Now().Add(-time.Duration(rand.Intn(1000)) * time.Hour),
			CreatedBy:     faker.Name(),
			UpdatedBy:     faker.Name(),
			ArchiveStatus: rand.Intn(2) == 1,
			DeletedStatus: rand.Intn(2) == 1,
		}

		teams = append(teams, team)
	}

	// Insert the generated records into MongoDB
	result, err := collection.InsertMany(context.Background(), teams)
	if err != nil {
		log.Fatalf("Failed to insert seed data: %v", err)
	}

	// Fetch every thing from id
	results, _ := collection.Find(context.Background(), bson.M{})
	results.All(context.Background(), &teams)

	fmt.Printf("Inserted %d teams: %+v\n", len(result.InsertedIDs), result.InsertedIDs)
	return teams
}

func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	t.Helper()

	bsonOpts := &options.BSONOptions{
		UseJSONStructTags:   true, // Replaces the bson struct tags with json struct tags
		ObjectIDAsHexString: true, // Allows the ObjectID to be marshalled as a string
		NilSliceAsEmpty:     true,
		UseLocalTimeZone:    false,
	}

	client, err := mongo.Connect(options.Client().ApplyURI(testDBURI).SetBSONOptions(bsonOpts))
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	db := client.Database("test_db")

	// Cleanup function to run after each test
	cleanup := func() {
		if err := db.Drop(context.TODO()); err != nil {
			t.Fatalf("Failed to clean up test database: %v", err)
		}
		client.Disconnect(context.TODO())
	}

	return db, cleanup
}

func TestTeam_CreateTeam(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.TODO()
	team := CTeam{
		Name:        "Test Team",
		Description: "This is a test team",
		TeamLead:    "lead123",
		TeamMember:  []string{"member1", "member2"},
		CreatedBy:   "admin",
		UpdatedBy:   "admin",
	}

	result, err, status := CreateTeam(team, ctx, db)
	if err != nil || status != http.StatusCreated {
		t.Fatalf("Expected successful creation, got error: %v, status: %d", err, status)
	}

	var storedTeam Team
	err = db.Collection("teams").FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&storedTeam)
	if err != nil {
		t.Fatalf("Failed to find inserted team: %v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, "Test Team", storedTeam.Name)
}

func TestTeam_AddTeamMember(t *testing.T) {
	type AddMemberTest struct {
		Title          string
		NewMembers     []string
		CurrentTeam    Team
		ExpectedResult Team
		ExpectError    error
		ExpectedStatus int
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	teams := generateRandomTeamSeeds(db)
	team := teams[rand.Intn(len(teams))]

	table := []AddMemberTest{
		{
			Title:       "Test adding a new member to the team",
			CurrentTeam: team,
			NewMembers:  []string{"member3"},
			ExpectedResult: Team{
				TeamMember: append(team.TeamMember, "member3"),
			},
			ExpectError:    nil,
			ExpectedStatus: http.StatusAccepted,
		},
		{
			Title:       "Test adding a four new member to the team",
			CurrentTeam: team,
			NewMembers:  []string{"member3", "member4", "member5", "member6"},
			ExpectedResult: Team{
				TeamMember: append(team.TeamMember, "member3", "member4", "member5", "member6"),
			},
			ExpectError:    nil,
			ExpectedStatus: http.StatusAccepted,
		},
	}

	for _, tt := range table {
		objId, _ := util.GetPrimitiveID(tt.CurrentTeam.ID)

		_, err, code := tt.CurrentTeam.AddNewTeamMember(tt.NewMembers, objId, db, context.Background())

		assert.Equal(t, len(tt.ExpectedResult.TeamMember), len(tt.CurrentTeam.TeamMember))
		assert.Equal(t, tt.ExpectError, err)
		assert.Equal(t, tt.ExpectedStatus, code)
	}

}

func TestTeam_RemoveTeamMember(t *testing.T) {
	type RemoveMemberTest struct {
		Title          string
		RemoveMembers  []string
		CurrentTeam    Team
		ExpectedResult Team
		ExpectError    error
		ExpectedStatus int
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	teams := generateRandomTeamSeeds(db)

	team := teams[rand.Intn(len(teams))]

	teamMembers := team.TeamMember

	table := []RemoveMemberTest{
		{
			Title:         "Test remove member from the team",
			CurrentTeam:   team,
			RemoveMembers: []string{teamMembers[0]},
			ExpectedResult: Team{
				TeamMember: append(teamMembers[:0], teamMembers[1:]...),
			},
			ExpectError:    nil,
			ExpectedStatus: http.StatusAccepted,
		},
	}

	for _, tt := range table {
		objId, _ := bson.ObjectIDFromHex(tt.CurrentTeam.ID)

		tm, err, code := tt.CurrentTeam.RemoveTeamMember(tt.RemoveMembers, objId, db, context.Background())

		assert.Equal(t, tt.ExpectedStatus, code)
		assert.Equal(t, tt.ExpectError, err)
		assert.Equal(t, len(tt.ExpectedResult.TeamMember), len(tm.TeamMember))
		assert.NotEqual(t, tm, nil, "The team cannot is nil")

	}

}

func TestTeam_ChangeTeamLead(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	teams := generateRandomTeamSeeds(db)

	team := teams[rand.Intn(len(teams))]

	l := len(team.TeamMember)
	var newMember string

	for isOk := true; isOk; {
		m := team.TeamMember[rand.Intn(l)]
		if m != team.TeamLead {
			isOk = false
			newMember = m
			break
		}
	}

	type ChangeLead struct {
		Title          string
		NewLead        string
		TeamLeadIdSent string
		CurrentTeam    Team
		ExpectedResult *Team
		ExpectError    error
		ExpectedStatus int
	}

	table := []ChangeLead{
		{
			Title:          "Test change team lead",
			CurrentTeam:    team,
			NewLead:        newMember,
			TeamLeadIdSent: team.TeamLead,
			ExpectError:    nil,
			ExpectedStatus: http.StatusOK,
			ExpectedResult: &Team{
				TeamLead:   newMember,
				TeamMember: team.TeamMember,
			},
		},
		{
			Title:          "Test failed -> sent wrong team lead data",
			CurrentTeam:    team,
			NewLead:        newMember,
			TeamLeadIdSent: "wrongLead",
			ExpectError:    errors.New(fmt.Sprintf("the current lead of id %s is not accurate with the lead id sent", "wrongLead")),
			ExpectedStatus: http.StatusOK,
			ExpectedResult: nil,
		},
		{
			Title:          "Test add a completely new user who is not on team to be lead",
			CurrentTeam:    team,
			NewLead:        "newLead",
			TeamLeadIdSent: team.TeamLead,
			ExpectError:    nil,
			ExpectedStatus: http.StatusOK,
			ExpectedResult: &Team{
				TeamLead:   "newLead",
				TeamMember: append(team.TeamMember, "newLead"),
			},
		},
	}

	for _, tt := range table {
		objId, _ := util.GetPrimitiveID(tt.CurrentTeam.ID)
		tm, err, code := tt.CurrentTeam.ChangeTeamLead(tt.NewLead, tt.TeamLeadIdSent, objId, db, context.Background())

		assert.Equal(t, tt.ExpectError, err)
		assert.Equal(t, tt.ExpectedStatus, code)
		if tm != nil {
			assert.Equal(t, tt.ExpectedResult.TeamLead, tm.TeamLead)
			assert.Equal(t, len(tt.ExpectedResult.TeamMember), len(tm.TeamMember))
		}
	}

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

func TestTeam_ArchiveTeam(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	type ExpectedOutCome struct {
		ExpectedStatus int
		ExpectedError  error
		ExpectResult   *Team
	}

	expectedOutCome := ExpectedOutCome{
		ExpectedStatus: http.StatusOK,
		ExpectedError:  nil,
		ExpectResult: &Team{
			ArchiveStatus: true,
		},
	}

	var nonArchivedTeams []Team
	teams := generateRandomTeamSeeds(db)

	for _, team := range teams {
		if !team.ArchiveStatus {
			nonArchivedTeams = append(nonArchivedTeams, team)
		}
	}

	for _, nonArchivedTeam := range nonArchivedTeams {
		objId, _ := util.GetPrimitiveID(nonArchivedTeam.ID)
		tm, err, code := nonArchivedTeam.ArchiveTeam(context.TODO(), objId, db)

		assert.Equal(t, expectedOutCome.ExpectedStatus, code)
		assert.Equal(t, expectedOutCome.ExpectedError, err)
		if err == nil {
			assert.Equal(t, expectedOutCome.ExpectResult.ArchiveStatus, tm.ArchiveStatus)
		}
	}

}

func TestTeam_UnArchiveTeam(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	type ExpectedOutCome struct {
		ExpectedStatus int
		ExpectedError  error
		ExpectResult   *Team
	}

	expectedOutCome := ExpectedOutCome{
		ExpectedStatus: http.StatusOK,
		ExpectedError:  nil,
		ExpectResult: &Team{
			ArchiveStatus: false,
		},
	}

	var archivedTeams []Team
	teams := generateRandomTeamSeeds(db)

	for _, team := range teams {
		if team.ArchiveStatus {
			archivedTeams = append(archivedTeams, team)
		}
	}

	for _, archivedTeam := range archivedTeams {
		objId, _ := util.GetPrimitiveID(archivedTeam.ID)
		tm, err, code := archivedTeam.UnArchiveTeam(context.TODO(), objId, db)

		assert.Equal(t, expectedOutCome.ExpectedStatus, code)
		assert.Equal(t, expectedOutCome.ExpectedError, err)
		if err == nil {
			assert.Equal(t, expectedOutCome.ExpectResult.ArchiveStatus, tm.ArchiveStatus)
		}
	}
}

func TestTeam_HandleCreatTeam(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	type TestKit struct {
		Title            string
		Body             CTeam
		ExpectedStatus   int
		ExpectedResponse util.Response
	}

	tests := []TestKit{
		{
			Title: "Test create team",
			Body: CTeam{
				Name:        "Test Team",
				Description: "This is a test team",
				TeamLead:    "member1",
				TeamMember:  []string{"member1", "member2"},
				CreatedBy:   "admin",
				UpdatedBy:   "admin",
			},
			ExpectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		body, _ := json.Marshal(tt.Body)

		req := httptest.NewRequest("POST", "/api/v1/teams/create", strings.NewReader(string(body)))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		HandleCreateTeam(db).ServeHTTP(rr, req)

		assert.Equal(t, tt.ExpectedStatus, rr.Code)
	}
}

func TestTeam_GetTeams(t *testing.T) {
	type TestKit struct {
		ExpectedStatus  int
		ExpectedError   error
		ExpectedDataLen int
	}

	t.Run("Test get all teams", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		teams := generateRandomTeamSeeds(db)

		tb := []TestKit{
			{
				ExpectedStatus:  http.StatusOK,
				ExpectedError:   nil,
				ExpectedDataLen: len(teams),
			},
		}

		for _, tt := range tb {
			req := httptest.NewRequest("GET", "/api/v1/teams", nil)
			rr := httptest.NewRecorder()
			rr.Header().Set("Content-Type", "application/json")

			// Stimulate the response
			b, _ := util.GetBytesResponse(http.StatusOK, teams)
			rr.Write(b)

			GetTeams(db).ServeHTTP(rr, req)

			var resp util.Response
			json.NewDecoder(rr.Body).Decode(&resp)

			assert.Equal(t, rr.Code, tt.ExpectedStatus)
			assert.Equal(t, tt.ExpectedError, resp.Error)
		}
	})

	t.Run("Test failure, no records in db", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		tb := []TestKit{
			{
				ExpectedStatus:  http.StatusOK,
				ExpectedError:   mongo.ErrNoDocuments,
				ExpectedDataLen: 0,
			},
		}

		for _, tt := range tb {
			req := httptest.NewRequest("GET", "/api/v1/teams", nil)
			rr := httptest.NewRecorder()

			rr.Header().Set("Content-Type", "application/json")
			rr.Code = http.StatusOK

			// Stimulate the response
			util.ErrorException(rr, tt.ExpectedError, tt.ExpectedStatus)

			GetTeams(db).ServeHTTP(rr, req)

			var resp map[string]string
			json.NewDecoder(rr.Body).Decode(&resp)

			assert.Equal(t, rr.Code, tt.ExpectedStatus)
			assert.Equal(t, tt.ExpectedError.Error(), resp["error"])
		}

	})

}

func TestTeam_GetTeam(t *testing.T) {
	type TestKit struct {
		requestId      string
		ExpectedStatus int
		ExpectedError  error
		ExpectedData   *Team
	}

	t.Run("Test get team -> no record", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		tb := []TestKit{
			{
				"123",
				http.StatusNotFound,
				mongo.ErrNoDocuments,
				nil,
			},
			{
				"doc1",
				http.StatusNotFound,
				mongo.ErrNoDocuments,
				nil,
			},
		}

		for _, tt := range tb {
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/teams/%s", tt.requestId), nil)
			rr := httptest.NewRecorder()
			rr.Header().Set("Content-Type", "application/json")

			rr.Code = tt.ExpectedStatus
			util.ErrorException(rr, mongo.ErrNoDocuments, http.StatusNotFound)

			GetTeam(db).ServeHTTP(rr, req)

			var resp map[string]string
			json.NewDecoder(rr.Body).Decode(&resp)

			assert.Equal(t, rr.Code, tt.ExpectedStatus)
			assert.Equal(t, tt.ExpectedError.Error(), resp["error"])
		}

	})

	t.Run("Test get team by id", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		teams := generateRandomTeamSeeds(db)

		tb := []TestKit{
			{
				teams[0].ID,
				http.StatusOK,
				nil,
				&teams[0],
			},
			{
				teams[1].ID,
				http.StatusOK,
				nil,
				&teams[1],
			},
		}

		for _, tt := range tb {
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/teams/%s", tt.requestId), nil)
			rr := httptest.NewRecorder()
			rr.Header().Set("Content-Type", "application/json")

			rr.Code = tt.ExpectedStatus
			b, _ := util.GetBytesResponse(http.StatusOK, tt.ExpectedData)
			rr.Header().Set("Content-Type", "application/json")
			rr.Write(b)

			GetTeam(db).ServeHTTP(rr, req)

			var resp util.Response
			json.NewDecoder(rr.Body).Decode(&resp)

			assert.Equal(t, rr.Code, tt.ExpectedStatus)
			assert.Equal(t, tt.ExpectedError, resp.Error)
			assert.Equal(t, tt.ExpectedData.Name, resp.Data.(map[string]interface{})["name"])
		}
	})

	t.Run("Test get team by searching name", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		teams := generateRandomTeamSeeds(db)

		tb := []TestKit{
			{
				teams[0].Name,
				http.StatusOK,
				nil,
				&teams[0],
			},
			{
				teams[1].Name,
				http.StatusOK,
				nil,
				&teams[1],
			},
		}

		for _, tt := range tb {
			log.Println("Name", tt.requestId)
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/teams/%c", tt.requestId[0]), nil)
			rr := httptest.NewRecorder()
			rr.Header().Set("Content-Type", "application/json")

			rr.Code = tt.ExpectedStatus
			b, _ := util.GetBytesResponse(http.StatusOK, tt.ExpectedData)
			rr.Header().Set("Content-Type", "application/json")
			rr.Write(b)

			GetTeam(db).ServeHTTP(rr, req)

			log.Printf("resp: %+v", rr.Body)

			var resp util.Response
			json.NewDecoder(rr.Body).Decode(&resp)

			assert.Equal(t, rr.Code, tt.ExpectedStatus)
			assert.Equal(t, tt.ExpectedError, resp.Error)
			assert.Equal(t, tt.ExpectedData.Name, resp.Data.(map[string]interface{})["name"])
		}
	})

}

func TestTeam_HandleArchiveTeam(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	teams := generateRandomTeamSeeds(db)

	var nonArchiveTeams []Team
	var archiveTeams []Team

	for _, team := range teams {
		if team.ArchiveStatus {
			archiveTeams = append(archiveTeams, team)
		} else {
			nonArchiveTeams = append(nonArchiveTeams, team)
		}
	}

	t.Run("Test archive teams", func(t *testing.T) {
		for _, tt := range nonArchiveTeams {
			body, _ := json.Marshal(tt)

			// Mock expected team attribute change
			expectedTeam := tt
			expectedTeam.ArchiveStatus = true

			req := httptest.NewRequest("PATCH", "/api/teams/archive", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			// mock response
			respBytes, _ := util.GetBytesResponse(http.StatusAccepted, expectedTeam)

			rr.WriteHeader(http.StatusAccepted)
			rr.Write(respBytes)

			HandleArchiveTeam(db).ServeHTTP(rr, req)

			var resp util.Response
			json.NewDecoder(rr.Body).Decode(&resp)

			assert.Equal(t, resp.Error, nil)
			assert.Equal(t, rr.Code, http.StatusAccepted)
			assert.Equal(t, rr.Code, resp.Status)
			assert.Equal(t, expectedTeam.ArchiveStatus, resp.Data.(map[string]interface{})["archive_status"])
		}
	})

	t.Run("Test cannot archive teams as they are already archived", func(t *testing.T) {
		for _, tt := range archiveTeams {
			body, _ := json.Marshal(tt)

			// Mock Expected err and code
			expectedErr := errors.New("team has already been archived")
			expectedCode := http.StatusOK

			req := httptest.NewRequest("PATCH", "/api/teams/archive", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			HandleArchiveTeam(db).ServeHTTP(rr, req)

			var resp map[string]string
			json.NewDecoder(rr.Body).Decode(&resp)

			assert.Equal(t, rr.Code, expectedCode)
			assert.Equal(t, tt.ArchiveStatus, true)
			assert.Equal(t, expectedErr.Error(), resp["error"])
		}
	})
}

func TestTeam_HandleUnArchiveTeam(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	teams := generateRandomTeamSeeds(db)

	var nonArchiveTeams []Team
	var archiveTeams []Team

	for _, team := range teams {
		if team.ArchiveStatus {
			archiveTeams = append(archiveTeams, team)
		} else {
			nonArchiveTeams = append(nonArchiveTeams, team)
		}
	}

	t.Run("Test un-archive teams successfully", func(t *testing.T) {
		for _, tt := range archiveTeams {
			expectedCode := http.StatusOK
			expectedTeam := tt
			expectedTeam.ArchiveStatus = false

			body, _ := json.Marshal(tt)

			req := httptest.NewRequest("PATCH", "/api/v1/teams/unarchive", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			HandleUnArchiveTeam(db).ServeHTTP(rr, req)

			var response util.Response
			json.NewDecoder(rr.Body).Decode(&response)

			assert.Equal(t, rr.Code, expectedCode)
			assert.Equal(t, expectedCode, response.Status)
			assert.Equal(t, response.Error, nil)
			assert.Equal(t, expectedTeam.ArchiveStatus, response.Data.(map[string]interface{})["archive_status"])
		}
	})

	t.Run("", func(t *testing.T) {
		for _, tt := range nonArchiveTeams {
			// Expected outcomes
			expectedErr := errors.New("team is not in the archive catalogue")
			expectedCode := http.StatusOK

			body, _ := json.Marshal(tt)

			req := httptest.NewRequest("PATCH", "/api/v1/teams/unarchive", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			HandleUnArchiveTeam(db).ServeHTTP(rr, req)

			var response map[string]string
			json.NewDecoder(rr.Body).Decode(&response)

			assert.Equal(t, rr.Code, expectedCode)
			assert.Equal(t, expectedErr.Error(), response["error"])
		}
	})

}

func TestTeam_HandleAddNewMember(t *testing.T) {
	type MockTb struct {
		TeamId         string
		ExpectedStatus int
		ExpectedError  error
		ExpectedTeam   *Team
		ReqBody        *CBody
	}

	t.Run("Test add member successfully", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		teams := generateRandomTeamSeeds(db)

		for _, team := range teams {

			mt := MockTb{
				ExpectedError:  nil,
				ExpectedStatus: http.StatusAccepted,
				TeamId:         team.ID,
				ExpectedTeam: &Team{
					TeamMember: append(team.TeamMember, "member2"),
				},
				ReqBody: &CBody{
					Team:        team,
					TeamMembers: []string{"member2"},
				},
			}

			bdy, _ := json.Marshal(mt.ReqBody)
			req := httptest.NewRequest("PATCH", "/api/v1/teams/add-member", bytes.NewReader(bdy))
			rr := httptest.NewRecorder()

			HandleAddNewMembers(db).ServeHTTP(rr, req)

			var resp util.Response
			json.NewDecoder(rr.Body).Decode(&resp)

			r := resp.Data.(map[string]interface{})["team_member"].([]interface{})

			assert.Equal(t, rr.Code, mt.ExpectedStatus)
			assert.Equal(t, resp.Error, mt.ExpectedError)
			assert.Equal(t, len(mt.ExpectedTeam.TeamMember), len(r))
		}
	})

	t.Run("Test expect error when adding a members successfully", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		teams := generateRandomTeamSeeds(db)

		for _, team := range teams {

			mt := MockTb{
				ExpectedError:  nil,
				ExpectedStatus: http.StatusInternalServerError,
				TeamId:         team.ID,
				ExpectedTeam: &Team{
					TeamMember: team.TeamMember,
				},
				ReqBody: nil,
			}

			bdy, _ := json.Marshal(mt.ReqBody)
			req := httptest.NewRequest("PATCH", "/api/v1/teams/add-member", bytes.NewReader(bdy))

			rr := httptest.NewRecorder()

			HandleAddNewMembers(db).ServeHTTP(rr, req)

			var resp map[string]string
			json.NewDecoder(rr.Body).Decode(&resp)

			assert.Equal(t, rr.Code, mt.ExpectedStatus)
			assert.NotNil(t, resp["error"])
			assert.Equal(t, len(mt.ExpectedTeam.TeamMember), len(team.TeamMember))
		}
	})
}

func TestTeam_HandleRemoveNewMember(t *testing.T) {}

func TestTeam_HardDeleteTeam(t *testing.T) {
	t.Run("Test deleting docs from collection", func(t *testing.T) {
		type Del struct {
			ExpectedStatus int
			Tm             Team
		}

		db, cleanup := setupTestDB(t)
		defer cleanup()

		teams := generateRandomTeamSeeds(db)

		for _, team := range teams {

			tm := Del{
				ExpectedStatus: http.StatusAccepted,
				Tm:             team,
			}

			byt, _ := json.Marshal(tm.Tm)
			req := httptest.NewRequest("DELETE", "/api/v1/teams/delete", bytes.NewReader(byt))
			rr := httptest.NewRecorder()

			HardDeleteTeam(db).ServeHTTP(rr, req)

			var rep util.Response
			json.NewDecoder(rr.Body).Decode(&rep)

			assert.Equal(t, rr.Code, tm.ExpectedStatus)
			assert.Equal(t, rep.Status, tm.ExpectedStatus)
			assert.Nil(t, rep.Error)
		}
	})

	t.Run("Test deleting a document Failed -> no document in collection ", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		docId := bson.NewObjectID().Hex()

		tm := Team{
			Name: "name",
			ID:   docId,
		}

		byt, _ := json.Marshal(tm)
		req := httptest.NewRequest("DELETE", "/api/v1/teams/delete", bytes.NewReader(byt))
		rr := httptest.NewRecorder()

		HardDeleteTeam(db).ServeHTTP(rr, req)
		var rep map[string]string
		json.NewDecoder(rr.Body).Decode(&rep)

		assert.Equal(t, rr.Code, http.StatusOK)
		assert.NotNil(t, rep["error"])
		assert.Equal(t, errors.New(rep["error"]), errors.New("no team matching the record was found and hence it can't be deleted"))
	})

	t.Run("Test deleting Failed -> Inputing a wrong docId", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		generateRandomTeamSeeds(db)

		docId := bson.NewObjectID().String()

		tm := Team{
			Name: "name",
			ID:   docId,
		}

		byt, _ := json.Marshal(tm)
		req := httptest.NewRequest("DELETE", "/api/v1/teams/delete", bytes.NewReader(byt))
		rr := httptest.NewRecorder()

		HardDeleteTeam(db).ServeHTTP(rr, req)
		var rep map[string]string
		json.NewDecoder(rr.Body).Decode(&rep)

		assert.Equal(t, rr.Code, http.StatusInternalServerError)
		assert.NotNil(t, rep["error"])
		assert.Equal(t, errors.New(rep["error"]), errors.New("the provided hex string is not a valid ObjectID"))
	})
}

func TestTeam_HandleChangeTeamLead(t *testing.T) {}

func TestTeam_PushTeamToBin(t *testing.T) {}

func TestTeam_RestoreTeamFromBin(t *testing.T) {}
