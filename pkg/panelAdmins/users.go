package panelAdmins

import (
	"context"
	"control-panel-bk/config"
	"control-panel-bk/internal/aws"
	"control-panel-bk/util"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"net/http"
	"strconv"
	"time"
)

const (
	MAX_LIMIT = 50
)

type Personal struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Gender    string    `json:"gender"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Dob       time.Time `json:"dob,omitempty"`
	Profile   string    `json:"profile,omitempty"`
}

type User struct {
	ID            string    `json:"id,omitempty"`
	Personal      Personal  `json:"personal"`
	RoleId        string    `json:"role_id"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedAt     time.Time `json:"created_at"`
	IsActive      bool      `json:"is_active"`
	ArchiveStatus bool      `json:"archive_status"`
	DeletedStatus bool      `json:"is_deleted_status"`
	CreatedBy     string    `json:"created_by"`
	UpdatedBy     string    `json:"updated_by"`

	// The id of the user from the cognito user pool
	UpId string `json:"up_id,omitempty"`
}

type NewUser struct {
	Personal
	RoleId     string `json:"role_id,omitempty"`
	teamId     string `json:"team_id,omitempty"`
	IsTeamLead bool   `json:"is_team_lead"`
	Role       CRole  `json:"role,omitempty"`
	CreatedBy  string `json:"created_by"`
	UpdatedBy  string `json:"updated_by"`
}

func getSession(client *mongo.Client) (session *mongo.Session, err error) {
	return client.StartSession()
}

func (u *User) DeActiveUser(db *mongo.Database) (*User, error) {
	if !u.IsActive {
		return nil, errors.New("user is currently deactivated")
	}

	userID, userIDErr := util.GetPrimitiveID(u.ID)
	if userIDErr != nil {
		return nil, userIDErr
	}

	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"is_active":      false,
			"archive_status": true,
			"updated_at":     time.Now(),
			"updated_by":     u.UpdatedBy,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := db.Collection("users").FindOneAndUpdate(context.TODO(), filter, update, opt).Decode(&u); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) ActiveUser(db *mongo.Database) (*User, error) {
	if u.IsActive {
		return nil, errors.New("user is currently active")
	}

	userID, userIDErr := util.GetPrimitiveID(u.ID)
	if userIDErr != nil {
		return nil, userIDErr
	}

	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"is_active":      true,
			"archive_status": false,
			"updated_at":     time.Now(),
			"updated_by":     u.UpdatedBy,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := db.Collection("users").FindOneAndUpdate(context.TODO(), filter, update, opt).Decode(&u); err != nil {
		return nil, err
	}

	return u, nil
}

func CreateUser(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var newUser NewUser
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

		session, err := getSession(client)

		if err != nil {
			util.ErrorException(w, fmt.Errorf("failed to start session %w", err), http.StatusInternalServerError)
			return
		}

		defer session.EndSession(r.Context())

		db := session.Client().Database("flowCx")

		col := db.Collection("users")
		tmCol := db.Collection("teams")

		// Steps in creating a new user (transactional operation in mongo)
		err = mongo.WithSession(r.Context(), session, func(ctx context.Context) error {

			// STEP 1: CHECK IF THE ROLE ID (if the roleId is not provided then we need to create a new role for the user)
			if newUser.RoleId == "" {
				// Create a role and assign it to the newUser.RoleId
				crl := newUser.Role
				if rl, err := CreateRole(crl, r.Context(), db); err != nil {
					return fmt.Errorf("failed to create the user's role")
				} else {
					newUser.RoleId = rl.Data.InsertedID.(bson.ObjectID).Hex()
				}
			}

			// STEP 2: CREATE THE USER IN COGNITO USER POOL
			userId, outputErr := aws.CreateNewUser(config.AwsConfig, newUser.Email, newUser.RoleId, util.DefaultPassword)

			if outputErr != nil {
				return fmt.Errorf("failed to create a user in the userpool")
			}

			// STEP 3: CREATE THE USER IN A MONGO "users" COLLECTION WITH THE USER ID FROM THE USER POOL IN THE STUB
			doc, docErr := col.InsertOne(r.Context(), bson.M{
				"first_name":        newUser.FirstName,
				"last_name":         newUser.LastName,
				"email":             newUser.Email,
				"phone_num":         newUser.Phone,
				"gender":            newUser.Gender,
				"dob":               newUser.Dob,
				"created_at":        time.Now(),
				"updated_at":        time.Now(),
				"role_id":           newUser.RoleId,
				"up_id":             userId,
				"is_active":         false, // Will be set to true when user changes passwords
				"archive_status":    false,
				"is_deleted_status": false,
				"created_by":        newUser.CreatedBy,
				"updated_by":        newUser.UpdatedBy,
			})

			if docErr != nil {
				return fmt.Errorf("failed to insert the user document in the users collection %w", docErr)
			}

			userID := doc.InsertedID.(bson.ObjectID).Hex() //doc.InsertedID.(string)

			// STEP 4: ADD THE USER ID FROM THE "teams" COLLECTION into the team he was added to if such was provided
			var team Team
			isAssignToTeam := len(newUser.teamId) > 0 || false

			if isAssignToTeam {
				teamID, teamErr := util.GetPrimitiveID(newUser.teamId)

				if teamErr != nil {
					return teamErr
				}

				e := tmCol.FindOne(r.Context(), bson.M{"_id": teamID}).Decode(&team)

				if e != nil {
					if errors.Is(e, mongo.ErrNoDocuments) {
						return mongo.ErrNoDocuments
					}

					return e
				}

				// Add user to the team and update it
				if newUser.IsTeamLead {
					// The changeTeamLead will add the userId as a member of the team if he is not a member
					_, tlErr, _ := team.ChangeTeamLead(userID, team.TeamLead, teamID, db, r.Context())
					if tlErr != nil {
						return tlErr
					}
				} else {
					mbr := []string{userID}
					if _, e, _ := team.AddNewTeamMember(mbr, teamID, db, r.Context()); e != nil {
						return e
					}
				}
			}

			return nil

		})

		if err != nil {
			// cognito roll back
			if _, aErr := aws.DeleteUser(config.AwsConfig, newUser.Email); aErr != nil {
				util.ErrorException(w, aErr, http.StatusNotImplemented)
				return
			}

			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

		// We will need to find the user by email
		var response util.Response

		findErr := col.FindOne(r.Context(), bson.M{"email": newUser.Email}).Decode(&response.Data)

		if findErr != nil {
			util.ErrorException(w, findErr, http.StatusNotFound)
			return
		}

		response.Status = http.StatusCreated
		response.Error = nil

		respBy, respErr := util.GetBytesResponse(response.Status, response.Data)
		if respErr != nil {
			util.ErrorException(w, respErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(respBy)
	}
}

func GetUsers(db *mongo.Database) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var users []User
		var limit int64
		var skip int64

		var page int

		query := request.URL.Query()

		if len(query.Get("page")) > 0 {
			pg, pgErr := strconv.Atoi(query.Get("page"))
			if pgErr != nil {
				util.ErrorException(writer, pgErr, http.StatusInternalServerError)
				return
			}
			page = pg
		} else {
			page = 1
		}

		if len(query.Get("limit")) > 0 {
			lmt, lmtErr := strconv.Atoi(query.Get("limit"))
			if lmtErr != nil {
				util.ErrorException(writer, lmtErr, http.StatusInternalServerError)
				return
			}
			limit = int64(lmt)
		} else {
			limit = int64(MAX_LIMIT)
		}

		if page > 1 {
			skip = int64(page - 1)
		} else {
			skip = int64(0)
		}

		opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"created_at": -1})
		coll, err := db.Collection("users").Find(request.Context(), bson.M{}, opts)

		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				util.ErrorException(writer, mongo.ErrNoDocuments, http.StatusOK)
				return
			}

			util.ErrorException(writer, err, http.StatusNotImplemented)
			return
		}

		if err := coll.All(context.TODO(), &users); err != nil {
			util.ErrorException(writer, err, http.StatusNotImplemented)
			return
		}

		var response = util.Response{
			Status: http.StatusOK,
			Data:   users,
			Error:  nil,
		}

		respByt, respErr := util.GetBytesResponse(response.Status, response)

		if respErr != nil {
			util.ErrorException(writer, respErr, http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(response.Status)
		writer.Write(respByt)
	}
}
