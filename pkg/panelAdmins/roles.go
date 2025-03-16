package panelAdmins

import (
	"context"
	"control-panel-bk/internal/aws"
	"errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log"
	"net/http"
	"time"
)

type CreateRoleResponse struct {
	Status  bool
	Message string
	Data    *mongo.InsertOneResult
}

type Role struct {
	ID              string     `json:"_id,omitempty" bson:"_id,omitempty"`
	Name            string     `json:"name" bson:"name"`
	Description     string     `json:"description,omitempty" bson:"description,omitempty"`
	Permission      Permission `json:"permission" bson:"permission"`
	CreatedBy       string     `json:"createdBy" bson:"created_by"`
	UpdatedBy       string     `json:"updatedBy" bson:"updated_by"`
	ArchiveStatus   bool       `json:"archiveStatus" bson:"archive_status"`
	IsDeletedStatus bool       `json:"is_deleted_status" bson:"is_deleted_status"`
	UpdatedAt       time.Time  `json:"UpdatedAt,omitempty" bson:"updated_at"`
	CreatedAt       time.Time  `json:"createdAt,omitempty" bson:"created_at"`
}

type CRole struct {
	Name        string     `json:"name" bson:"name"`
	Description string     `json:"description,omitempty" bson:"description,omitempty"`
	Permission  Permission `json:"permission" bson:"permission"`
	CreatedBy   string     `json:"createdBy" bson:"created_by"`
	UpdatedBy   string     `json:"updatedBy" bson:"updated_by"`
	UpdatedAt   time.Time  `json:"UpdatedAt,omitempty" bson:"updated_at"`
	CreatedAt   time.Time  `json:"createdAt,omitempty" bson:"created_at"`
}

func getDB() *mongo.Database {
	if aws.MongoFlowCxDBClient == nil {
		panic("MongoDB database not initialized")
	}
	return aws.MongoFlowCxDBClient
}

func (rl *Role) UnArchiveRole(updatorId string) {}

func (rl *Role) ArchiveRole(updatorId string) {}

func (rl *Role) DeleteRole(updatorId string) {}

func (rl *Role) UpdateRoleDataToDB() {}

func FetchRoles(startNo int, endNo int, skip int) []Role {
	return []Role{}
}

func FetchRole(id string, ctx context.Context) (error, int) {
	db := getDB()
	var role bson.M
	err := db.Collection("roles").FindOne(ctx, bson.D{{"_id", id}}).Decode(&role)

	if err != nil {
		if errors.Is(err, mongo.ErrClientDisconnected) {
			return err, http.StatusInternalServerError
		}

		if errors.Is(err, mongo.ErrNoDocuments) {
			return err, http.StatusOK
		}

		return err, http.StatusBadRequest
	}

	log.Printf("role by id %v was found to be %#v", id, role)
	return nil, http.StatusFound
}

func CreateRole(crl CRole, ctx context.Context) (*CreateRoleResponse, error) {
	db := getDB()
	doc, err := db.Collection("roles").InsertOne(ctx, bson.D{
		{"name", crl.Name},
		{"description", crl.Description},
		{"permission", crl.Permission},
		{"created_by", crl.CreatedBy},
		{"update_by", crl.UpdatedBy},
		{"archive_status", false},
		{"archive_status", false},
		{"is_deleted_status", false},
	})

	if err != nil {
		return nil, err
	}

	response := CreateRoleResponse{
		Data:    doc,
		Message: "Role has been created",
		Status:  true,
	}

	return &response, nil
}

//update := bson.M{
//"$set": bson.M{
//"permission": updatedPermission,
//"updated_by": updatedBy,
//},
//"$currentDate": bson.M{
//"last_modified": true, // Automatically updates lastModified with the current timestamp
//},
//}
