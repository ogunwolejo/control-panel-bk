package panelAdmins

import (
	"context"
	"control-panel-bk/internal/aws"
	"errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
	ID              string     `json:"_id" bson:"_id"`
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

func (rl *Role) GeneralizedUpdate(ctx context.Context) (*Role, error) {
	db := getDB()

	filter := bson.M{
		"_id": rl.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"name":              rl.Name,
			"description":       rl.Description,
			"permission":        rl.Permission,
			"update_by":         rl.UpdatedBy,
			"archive_status":    rl.ArchiveStatus,
			"is_deleted_status": rl.IsDeletedStatus,
			"updated_at":        time.Now().UTC(),
		},
	}

	updated := db.Collection("roles").FindOneAndUpdate(ctx, filter, update)

	if updated.Err() != nil {
		if errors.Is(updated.Err(), mongo.ErrNoDocuments) {
			return nil, errors.New("no role with the selected metrics were found")
		}

		return nil, updated.Err()
	}

	if err := updated.Decode(rl); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *Role) UnArchiveRole(ctx context.Context) (*Role, error) {
	db := getDB()

	filter := bson.D{{"_id", rl.ID}}

	update := bson.M{
		"$set": bson.M{
			"updated_at":     time.Now().UTC(),
			"archive_status": false,
			"updated_by":     rl.UpdatedBy,
		},
	}

	updated, err := db.Collection("roles").UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("no document was found")
		}
		return nil, err
	}

	if updated.ModifiedCount == 0 {
		return nil, errors.New("could not be un-archive role, try again")
	} else {
		single := db.Collection("roles").FindOne(ctx, bson.M{"_id": rl.ID})

		if single.Err() != nil {
			return nil, single.Err()
		}

		if err := single.Decode(rl); err != nil {
			return nil, err
		}
	}

	return rl, nil
}

func (rl *Role) ArchiveRole(ctx context.Context) (*Role, error) {
	db := getDB()

	filter := bson.D{{"_id", rl.ID}}

	update := bson.M{
		"$set": bson.M{
			"updated_at":     time.Now().UTC(),
			"archive_status": true,
			"updated_by":     rl.UpdatedBy,
		},
	}

	updated, err := db.Collection("roles").UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("no document was found")
		}
		return nil, err
	}

	if updated.ModifiedCount == 0 {
		return nil, errors.New("could not archive role, try again")
	} else {
		single := db.Collection("roles").FindOne(ctx, bson.M{"_id": rl.ID})

		if single.Err() != nil {
			return nil, single.Err()
		}

		if err := single.Decode(rl); err != nil {
			return nil, err
		}
	}

	return rl, nil
}

func (rl *Role) DeleteRole(ctx context.Context) (*Role, error) {
	db := getDB()

	filter := bson.D{{"_id", rl.ID}}
	update := bson.M{
		"$set": bson.M{
			"updated_at":        time.Now().UTC(),
			"is_deleted_status": true,
			"updated_by":        rl.UpdatedBy,
		},
	}

	updated := db.Collection("roles").FindOneAndUpdate(ctx, filter, update)
	if updated.Err() != nil {
		if errors.Is(updated.Err(), mongo.ErrNoDocuments) {
			return nil, errors.New("no document was found")
		}
		return nil, updated.Err()
	}

	if err := updated.Decode(rl); err != nil {
		return nil, err
	}

	return rl, nil

}

func (rl *Role) HardDeleteRole(ctx context.Context) (*string, error, int) {
	db := getDB()

	del, err := db.Collection("roles").DeleteOne(ctx, bson.D{{"_id", rl.ID}})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err, http.StatusOK
		}

		return nil, err, http.StatusInternalServerError
	}

	if del.DeletedCount == 0 {
		return nil, errors.New("unable to delete role"), http.StatusNotImplemented
	}

	return &rl.ID, nil, http.StatusOK
}

func FetchRoles(page int, limit int, ctx context.Context) ([]Role, error, int) {
	var roles []Role

	db := getDB()

	lmt := int64(limit)
	var skip int64

	if page == 1 {
		skip = 0
	} else {
		skip = int64((page - 1) * limit)
	}

	opt := options.Find().SetSort(bson.D{{"created_at", -1}}).SetLimit(lmt).SetSkip(skip) //.SetSort(bson.D{{"created_at", 1}})
	docs, err := db.Collection("roles").Find(ctx, bson.D{}, opt)
	log.Println("DOCS: ", docs)

	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	defer docs.Close(ctx)

	if docs.Err() != nil {
		return nil, docs.Err(), http.StatusOK
	}

	for docs.Next(ctx) {
		var role Role

		if err := docs.Decode(&role); err == nil {
			roles = append(roles, role)
		}
	}

	return roles, nil, http.StatusOK
}

func FetchRoleById(id string, ctx context.Context) (*Role, error, int) {
	db := getDB()
	var role Role
	if err := db.Collection("roles").FindOne(ctx, bson.D{{"_id", id}}).Decode(&role); err != nil {
		if errors.Is(err, mongo.ErrClientDisconnected) {
			return nil, err, http.StatusInternalServerError
		}

		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err, http.StatusOK
		}

		return nil, err, http.StatusBadRequest
	}

	return &role, nil, http.StatusFound
}

func FetchRoleByName(roleName string, ctx context.Context) ([]Role, error) {
	db := getDB()

	var roles []Role

	doc, err := db.Collection("roles").Find(ctx, bson.D{{"name", roleName}})
	if err != nil {
		return nil, err
	}

	log.Println("DOCUMENTS ", doc)
	defer doc.Close(ctx)

	if doc.Err() != nil {
		return nil, doc.Err()
	}

	for doc.Next(ctx) {
		var role Role
		if err := doc.Decode(&role); err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func CreateRole(crl CRole, ctx context.Context) (*CreateRoleResponse, error) {
	result, err := FetchRoleByName(crl.Name, ctx)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		return nil, errors.New("a role having the same name already exists")
	}

	db := getDB()
	doc, err := db.Collection("roles").InsertOne(ctx, bson.D{
		{"name", crl.Name},
		{"description", crl.Description},
		{"permission", crl.Permission},
		{"created_by", crl.CreatedBy},
		{"update_by", crl.UpdatedBy},
		{"archive_status", false},
		{"is_deleted_status", false},
		{"created_at", time.Now().UTC()},
		{"updated_at", time.Now().UTC()},
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
