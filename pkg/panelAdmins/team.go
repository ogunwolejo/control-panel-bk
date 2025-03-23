package panelAdmins

import (
	"context"
	"control-panel-bk/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Team struct {
	ID            string    `json:"_id,omitempty"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	TeamLead      string    `json:"team_lead"`
	TeamMember    []string  `json:"team_member"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	CreatedBy     string    `json:"created_by"`
	UpdatedBy     string    `json:"updated_by"`
	ArchiveStatus bool      `json:"archive_status"`
	DeletedStatus bool      `json:"is_deleted_status"`
}

type CTeam struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	TeamLead    string   `json:"team_lead"`
	TeamMember  []string `json:"team_member"`
	CreatedBy   string   `json:"created_by"`
	UpdatedBy   string   `json:"updated_by"`
}

type CBody struct {
	TeamMembers []string `json:"team_members"`
	Team        Team     `json:"team"`
}

type CLead struct {
	Team    Team   `json:"team"`
	NewLead string `json:"new_lead"`
}

func CreateTeam(nt CTeam, ctx context.Context, client *mongo.Database) (*mongo.InsertOneResult, error, int) {
	db := client

	insert := bson.M{
		"updated_at":        time.Now(),
		"created_at":        time.Now(),
		"created_by":        nt.CreatedBy,
		"updated_by":        nt.UpdatedBy,
		"name":              nt.Name,
		"description":       nt.Description,
		"team_lead":         nt.TeamLead,
		"team_member":       nt.TeamMember,
		"archive_status":    false,
		"is_deleted_status": false,
	}

	create, err := db.Collection("teams").InsertOne(ctx, insert)
	if err != nil {
		return nil, err, http.StatusNotFound
	}

	return create, nil, http.StatusCreated
}

func (t *Team) AddNewTeamMember(member []string, teamId *bson.ObjectID, client *mongo.Database, ctx context.Context) (*Team, error, int) {
	db := client

	tm := append(t.TeamMember, member...)
	slices.SortStableFunc(tm, func(a, b string) int {
		return strings.Compare(a, b)
	})

	filter := bson.M{"_id": teamId}
	update := bson.M{
		"$set": bson.M{
			"updated_at":  time.Now(),
			"updated_by":  t.UpdatedBy,
			"team_member": tm,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := db.Collection("teams").FindOneAndUpdate(ctx, filter, update, opts).Decode(&t); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("no teams was found"), http.StatusOK
		}

		return nil, err, http.StatusNotFound
	}

	return t, nil, http.StatusAccepted
}

func (t *Team) RemoveTeamMember(mates []string, teamId bson.ObjectID, client *mongo.Database, ctx context.Context) (*Team, error, int) {
	db := client

	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, mate := range mates {
		wg.Add(1)

		go func(mate string) {
			defer wg.Done()

			mutex.Lock() // we lock the state, so that while we are performing some action nothing can be added or removed from the team members
			defer mutex.Unlock()

			for i, member := range t.TeamMember {
				if member == mate && t.TeamLead != member {
					t.TeamMember = append(t.TeamMember[:i], t.TeamMember[i+1:]...)
					break
				}
			}

		}(mate)

		wg.Wait()
	}

	filter := bson.D{{"_id", teamId}}
	update := bson.M{
		"$set": bson.M{
			"updated_by":  t.UpdatedBy,
			"updated_at":  time.Now(),
			"team_member": t.TeamMember,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := db.Collection("teams").FindOneAndUpdate(ctx, filter, update, opt).Decode(&t); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("team %s was not found, hence its data could not be updated", t.Name), http.StatusOK
		}
		return nil, err, http.StatusNotFound
	}

	return t, nil, http.StatusAccepted
}

func (t *Team) ChangeTeamLead(mate string, currentLead string, teamId *bson.ObjectID, client *mongo.Database, ctx context.Context) (*Team, error, int) {
	if t.TeamLead != currentLead {
		return nil, fmt.Errorf("the current lead of id %s is not accurate with the lead id sent", currentLead), http.StatusOK
	}

	cloneTeam := t

	if !t.IsMember(mate) {
		nm := make([]string, 0)
		nm = append(nm, mate)
		tm, err, cde := t.AddNewTeamMember(nm, teamId, client, ctx)

		if err != nil {
			return nil, err, cde // errors.New("unable to add proposed team lead to the team as the user is not a team member")
		}

		cloneTeam = tm
	}

	// Change the Team Member Status
	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, member := range cloneTeam.TeamMember {
		wg.Add(1)

		go func(mt string) {
			defer wg.Done()

			mutex.Lock()
			defer mutex.Unlock()

			// We set previous team lead to false
			if member != mt && member == t.TeamLead {
				t.TeamLead = mt // We assign the team lead to the new leader
			}

			// We change the team lead
			if member == mt && member != t.TeamLead {
				t.TeamLead = mt
			}

		}(mate)
		wg.Wait()
	}

	db := client

	objID, objErr := util.GetPrimitiveID(t.ID)
	if objErr != nil {
		return nil, objErr, http.StatusNotFound
	}

	filter := bson.M{"_id": objID}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	update := bson.M{
		"$set": bson.M{
			"updated_by":  t.UpdatedBy,
			"updated_at":  time.Now(),
			"team_member": t.TeamMember,
			"team_lead":   t.TeamLead,
		},
	}

	err := db.Collection("teams").FindOneAndUpdate(ctx, filter, update, opts).Decode(&t)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("no team record was found"), http.StatusOK
		}
		return nil, err, http.StatusNotFound
	}

	return t, nil, http.StatusOK
}

func (t *Team) SortTeamMembers() {
	slices.SortStableFunc(t.TeamMember, func(a, b string) int {
		return strings.Compare(a, b)
	})
}

func (t *Team) IsMember(mate string) bool {
	if isSorted := slices.IsSortedFunc(t.TeamMember, func(a, b string) int {
		return strings.Compare(a, b)
	}); !isSorted {
		t.SortTeamMembers()
	}

	_, isFound := slices.BinarySearchFunc(t.TeamMember, mate, func(a, b string) int {
		return strings.Compare(a, b)
	})

	return isFound
}

func (t *Team) ArchiveTeam(ctx context.Context, teamId *bson.ObjectID, client *mongo.Database) (*Team, error, int) {
	if t.ArchiveStatus {
		return nil, errors.New("team has already been archived"), http.StatusOK
	}

	db := client

	flt := bson.M{"_id": teamId}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	update := bson.D{{
		"$set", bson.M{
			"updated_at":     time.Now(),
			"updated_by":     t.UpdatedBy,
			"archive_status": true,
		},
	}}

	if err := db.Collection("teams").FindOneAndUpdate(ctx, flt, update, opt).Decode(t); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("team %s was not found", t.Name), http.StatusOK
		}
		return nil, err, http.StatusNotFound
	}

	return t, nil, http.StatusOK
}

func (t *Team) UnArchiveTeam(ctx context.Context, teamId *bson.ObjectID, client *mongo.Database) (*Team, error, int) {
	if !t.ArchiveStatus {
		return nil, errors.New("team is not in the archive catalogue"), http.StatusOK
	}

	db := client

	flt := bson.M{"_id": teamId}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	update := bson.D{{
		"$set", bson.M{
			"updated_at":     time.Now(),
			"updated_by":     t.UpdatedBy,
			"archive_status": false,
		},
	}}

	if err := db.Collection("teams").FindOneAndUpdate(ctx, flt, update, opt).Decode(t); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("team %s was not found", t.Name), http.StatusOK
		}
		return nil, err, http.StatusNotFound
	}

	return t, nil, http.StatusOK
}

// Handlers

func HandleCreateTeam(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var body CTeam
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

		result, e, code := CreateTeam(body, r.Context(), db)
		if e != nil {
			util.ErrorException(w, e, code)
			return
		}

		respBytes, respErr := util.GetBytesResponse(code, result)
		if respErr != nil {
			util.ErrorException(w, respErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		if _, err := w.Write(respBytes); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
		}
	}
}

func GetTeams(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		page, pgErr := strconv.Atoi(params.Get("page"))
		limit, lmtErr := strconv.Atoi(params.Get("limit"))

		if pgErr != nil {
			util.ErrorException(w, pgErr, http.StatusInternalServerError)
			return
		}

		if lmtErr != nil {
			util.ErrorException(w, lmtErr, http.StatusInternalServerError)
			return
		}

		var teams []Team

		// Calculating amount of docs to skip
		var skip int64 = 0
		if page > 1 {
			skip = int64((page - 1) * limit)
		}

		opt := options.Find().SetSkip(skip).SetLimit(int64(limit)).SetSort(bson.M{"created_by": -1})
		result, err := db.Collection("teams").Find(r.Context(), bson.M{}, opt)

		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				util.ErrorException(w, errors.New("no team record was found"), http.StatusOK)
				return
			}

			util.ErrorException(w, err, http.StatusNotFound)
			return
		}

		if err = result.All(r.Context(), &teams); err != nil {
			util.ErrorException(w, err, http.StatusNotImplemented)
			return
		}

		// The response nature
		respBytes, respBytesErr := util.GetBytesResponse(http.StatusOK, teams)
		if respBytesErr != nil {
			util.ErrorException(w, respBytesErr, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(respBytes); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}
	}
}

func GetTeam(client *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var teams []Team
		db := client //getDB(client)

		id := chi.URLParam(r, "id")

		// The id that be either the team by name or by id, hence if it matches it is an ID else it's a name
		doesMatch, err := regexp.Match("^[a-f0-9]{24}$", []byte(id))

		if err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

		// When the id matches th OBJECT ID type
		if doesMatch {
			objID, objErr := util.GetPrimitiveID(id)

			if objErr != nil {
				util.ErrorException(w, objErr, http.StatusInternalServerError)
				return
			}

			var team Team

			fil := bson.M{"_id": objID}
			opt := options.FindOne().SetSort(bson.M{"created_at": -1})
			err = db.Collection("teams").FindOne(r.Context(), fil, opt).Decode(&team)

			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					util.ErrorException(w, errors.New("no team record was found"), http.StatusOK)
					return
				}

				util.ErrorException(w, err, http.StatusNotFound)
				return
			}

			teams = append(teams, team)
		}

		// When the id does not match the OBJECT ID type
		if !doesMatch {
			fil := bson.M{
				"$text": bson.M{
					"$search": id,
				},
			}

			opt := options.Find().SetLimit(50).SetAllowPartialResults(true).SetSort(bson.M{"name": 1})
			results, resultsErr := db.Collection("teams").Find(r.Context(), fil, opt)
			if resultsErr != nil {
				if errors.Is(resultsErr, mongo.ErrNoDocuments) {
					util.ErrorException(w, resultsErr, http.StatusOK)
					return
				}

				util.ErrorException(w, resultsErr, http.StatusNotFound)
				return
			}

			if err = results.All(r.Context(), &teams); err != nil {
				util.ErrorException(w, err, http.StatusNotFound)
				return
			}
		}

		respByt, respErr := util.GetBytesResponse(http.StatusOK, teams)
		if respErr != nil {
			util.ErrorException(w, respErr, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(respByt)
	}
}

func HandleArchiveTeam(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var body Team
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

		objId, objErr := util.GetPrimitiveID(body.ID)
		if objErr != nil {
			util.ErrorException(w, objErr, http.StatusInternalServerError)
			return
		}

		result, err, code := body.ArchiveTeam(r.Context(), objId, db)
		if err != nil {
			util.ErrorException(w, err, code)
			return
		}

		respByt, respErr := util.GetBytesResponse(code, result)
		if respErr != nil {
			util.ErrorException(w, respErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(respByt)
	}
}

func HandleUnArchiveTeam(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var body Team
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

		objId, objErr := util.GetPrimitiveID(body.ID)
		if objErr != nil {
			util.ErrorException(w, objErr, http.StatusInternalServerError)
			return
		}

		result, err, code := body.UnArchiveTeam(r.Context(), objId, db)
		if err != nil {
			util.ErrorException(w, err, code)
			return
		}

		respByt, respErr := util.GetBytesResponse(code, result)
		if respErr != nil {
			util.ErrorException(w, respErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respByt)
	}
}

func HandleAddNewMembers(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var body CBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

		objId, objErr := util.GetPrimitiveID(body.Team.ID)
		if objErr != nil {
			util.ErrorException(w, objErr, http.StatusInternalServerError)
			return
		}

		result, e, code := body.Team.AddNewTeamMember(body.TeamMembers, objId, db, r.Context())
		if e != nil {
			util.ErrorException(w, e, code)
			return
		}

		respByte, respErr := util.GetBytesResponse(code, result)
		if respErr != nil {
			util.ErrorException(w, respErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(respByte)
	}
}

func HandleRemoveNewMembers(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var body CBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

		objId, objErr := util.GetPrimitiveID(body.Team.ID)
		if objErr != nil {
			util.ErrorException(w, objErr, http.StatusInternalServerError)
			return
		}

		result, e, code := body.Team.RemoveTeamMember(body.TeamMembers, *objId, db, r.Context())
		if e != nil {
			util.ErrorException(w, e, code)
			return
		}

		respByte, respErr := util.GetBytesResponse(code, result)
		if respErr != nil {
			util.ErrorException(w, respErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(respByte)
	}
}

func HardDeleteTeam(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var t Team
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

		objID, objErr := util.GetPrimitiveID(t.ID)
		if objErr != nil {
			util.ErrorException(w, objErr, http.StatusInternalServerError)
			return
		}

		if deleted := db.Collection("teams").FindOneAndDelete(r.Context(), bson.M{"_id": objID}).Decode(&t); deleted != nil {
			if errors.Is(deleted, mongo.ErrNoDocuments) {
				util.ErrorException(w, errors.New("no team matching the record was found and hence it can't be deleted"), http.StatusOK)
				return
			}

			util.ErrorException(w, deleted, http.StatusNotFound)
			return
		}

		respBy, respErr := util.GetBytesResponse(http.StatusAccepted, t.ID)
		if respErr != nil {
			util.ErrorException(w, respErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write(respBy)
	}
}

func HandleChangeTeamLead(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var body CLead
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

		objId, objErr := util.GetPrimitiveID(body.Team.ID)
		if objErr != nil {
			util.ErrorException(w, objErr, http.StatusInternalServerError)
			return
		}

		result, e, code := body.Team.ChangeTeamLead(body.NewLead, body.Team.TeamLead, objId, db, r.Context())
		if e != nil {
			util.ErrorException(w, e, code)
			return
		}

		respBy, respErr := util.GetBytesResponse(code, result)
		if respErr != nil {
			util.ErrorException(w, respErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(respBy)
	}
}

func PushTeamToBin(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var t Team
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

		if t.DeletedStatus {
			util.ErrorException(w, errors.New("team has already sent to the bin"), http.StatusOK)
			return
		}

		objId, objErr := util.GetPrimitiveID(t.ID)
		if objErr != nil {
			util.ErrorException(w, objErr, http.StatusInternalServerError)
			return
		}

		flt := bson.M{"_id": objId}
		opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
		update := bson.D{{
			"$set", bson.M{
				"updated_at":        time.Now(),
				"updated_by":        t.UpdatedBy,
				"is_deleted_status": true,
			},
		}}

		if err := db.Collection("teams").FindOneAndUpdate(r.Context(), flt, update, opt).Decode(&t); err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				util.ErrorException(w, fmt.Errorf("team %s was not found", t.Name), http.StatusOK)
				return
			}

			util.ErrorException(w, err, http.StatusNotFound)
			return
		}

		respBy, respErr := util.GetBytesResponse(http.StatusAccepted, t)
		if respErr != nil {
			util.ErrorException(w, respErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write(respBy)
	}
}

func RestoreTeamFromBin(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var t Team

		if !t.DeletedStatus {
			util.ErrorException(w, errors.New("team cannot be restored as it is not in the bin"), http.StatusOK)
			return
		}

		objId, objErr := util.GetPrimitiveID(t.ID)
		if objErr != nil {
			util.ErrorException(w, objErr, http.StatusInternalServerError)
			return
		}

		flt := bson.M{"_id": objId}
		opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
		update := bson.D{{
			"$set", bson.M{
				"updated_at":        time.Now(),
				"updated_by":        t.UpdatedBy,
				"is_deleted_status": false,
			},
		}}

		if err := db.Collection("teams").FindOneAndUpdate(r.Context(), flt, update, opt).Decode(&t); err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				util.ErrorException(w, fmt.Errorf("team %s was not found", t.Name), http.StatusOK)
				return
			}

			util.ErrorException(w, err, http.StatusNotFound)
			return
		}

		respBy, respErr := util.GetBytesResponse(http.StatusAccepted, t)
		if respErr != nil {
			util.ErrorException(w, respErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respBy)
	}
}
