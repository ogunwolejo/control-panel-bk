package panelAdmins

import (
	"context"
	"control-panel-bk/internal/aws"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"net/http"
	"strings"
	"time"
)

type CreateRoleResponse struct {
	Status  bool
	Message string
	Data    *mongo.InsertOneResult
}

type Response struct {
	Status int
	error  error
	Data   interface{}
}

type Role struct {
	ID              string     `json:"_id"`
	Name            string     `json:"name"`
	Description     string     `json:"description,omitempty"`
	Permission      Permission `json:"permission"`
	CreatedBy       string     `json:"created_by"`
	UpdatedBy       string     `json:"updated_by"`
	ArchiveStatus   bool       `json:"archive_status"`
	IsDeletedStatus bool       `json:"is_deleted_status"`
	UpdatedAt       time.Time  `json:"updated_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at,omitempty"`
}

type CRole struct {
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Permission  Permission `json:"permission"`
	CreatedBy   string     `json:"created_by"`
	UpdatedBy   string     `json:"updated_by"`
}

func getDB() *mongo.Database {
	if aws.MongoFlowCxDBClient == nil {
		panic("MongoDB database not initialized")
	}
	return aws.MongoFlowCxDBClient
}

func getPrimitiveID(id string) (*bson.ObjectID, error) {
	objId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	return &objId, nil
}

func (rl *Role) GeneralizedUpdate(ctx context.Context) (*Role, error, int) {
	db := getDB()

	objId, objErr := getPrimitiveID(rl.ID)
	if objErr != nil {
		return nil, objErr, http.StatusInternalServerError
	}

	filter := bson.M{
		"_id": objId,
	}

	update := bson.M{
		"$set": bson.M{
			"name":        rl.Name,
			"description": rl.Description,
			"permission":  rl.Permission,
			"update_by":   rl.UpdatedBy,
			"updated_at":  time.Now().UTC(),
		},
	}

	fuOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := db.Collection("roles").FindOneAndUpdate(ctx, filter, update, fuOpts).Decode(&rl)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("no role with the selected metrics were found"), http.StatusOK
		}

		return nil, err, http.StatusNotFound
	}

	return rl, nil, http.StatusOK
}

func (rl *Role) UnArchiveRole(ctx context.Context) (*Role, error, int) {
	db := getDB()

	if !rl.ArchiveStatus {
		return nil, fmt.Errorf("role %s is not archived, hence this task cannot be performed", rl.Name), http.StatusBadRequest
	}

	objId, objErr := getPrimitiveID(rl.ID)
	if objErr != nil {
		return nil, objErr, http.StatusInternalServerError
	}

	filter := bson.D{{"_id", objId}}

	update := bson.M{
		"$set": bson.M{
			"updated_at":     time.Now().UTC(),
			"archive_status": false,
			"updated_by":     rl.UpdatedBy,
		},
	}

	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := db.Collection("roles").FindOneAndUpdate(ctx, filter, update, opt).Decode(rl)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("no role named %s, was found", rl.Name), http.StatusOK
		}

		return nil, err, http.StatusNotFound
	}

	return rl, nil, http.StatusOK
}

func (rl *Role) ArchiveRole(ctx context.Context) (*Role, error, int) {
	if rl.ArchiveStatus {
		return nil, errors.New("role is already archived"), http.StatusInternalServerError
	}

	db := getDB()

	objId, objErr := getPrimitiveID(rl.ID)
	if objErr != nil {
		return nil, objErr, http.StatusInternalServerError
	}

	filter := bson.D{{"_id", objId}}

	update := bson.M{
		"$set": bson.M{
			"updated_at":     time.Now().UTC(),
			"archive_status": true,
			"updated_by":     rl.UpdatedBy,
		},
	}

	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := db.Collection("roles").FindOneAndUpdate(ctx, filter, update, opt).Decode(&rl)
	if err != nil {
		log.Printf("Error : %s", err.Error())
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("no matching role found to archive"), http.StatusOK
		}
		return nil, err, http.StatusInternalServerError
	}

	return rl, nil, http.StatusAccepted
}

func (rl *Role) DeleteRole(ctx context.Context) (*Role, error, int) {
	db := getDB()

	if rl.IsDeletedStatus {
		return nil, fmt.Errorf("role has been sent to the bin"), http.StatusOK
	}

	objId, objErr := getPrimitiveID(rl.ID)
	if objErr != nil {
		return nil, objErr, http.StatusInternalServerError
	}

	filter := bson.D{{"_id", objId}}
	update := bson.M{
		"$set": bson.M{
			"updated_at":        time.Now().UTC(),
			"is_deleted_status": true,
			"updated_by":        rl.UpdatedBy,
		},
	}

	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := db.Collection("roles").FindOneAndUpdate(ctx, filter, update, opt).Decode(&rl)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("no document was found"), http.StatusOK
		}
		return nil, err, http.StatusNotFound
	}

	return rl, nil, http.StatusOK
}

func (rl *Role) HardDeleteRole(ctx context.Context) (*string, error, int) {
	db := getDB()

	objId, objErr := getPrimitiveID(rl.ID)
	if objErr != nil {
		return nil, objErr, http.StatusInternalServerError
	}

	del, err := db.Collection("roles").DeleteOne(ctx, bson.D{{"_id", objId}})
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
	docs, err := db.Collection("roles").Find(ctx, bson.M{}, opt)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("no role has been created"), http.StatusOK
		}
		return nil, err, http.StatusInternalServerError
	}

	defer docs.Close(ctx)

	if docs.Err() != nil {
		return nil, docs.Err(), http.StatusOK
	}

	if err := docs.All(ctx, &roles); err != nil {
		return nil, err, http.StatusNotFound
	}

	return roles, nil, http.StatusOK
}

func FetchRoleById(roleId string, ctx context.Context) (*Role, error, int) {
	db := getDB()

	objID, err := getPrimitiveID(roleId)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err, http.StatusInternalServerError
	}

	var role Role

	flt := bson.M{"_id": objID}
	err = db.Collection("roles").FindOne(ctx, flt).Decode(&role)

	log.Println("role: ", role)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("record regarding this role was not found"), http.StatusOK
		}

		return nil, err, http.StatusNotFound
	}

	log.Printf("ROLE: %#v", role)

	return &role, nil, http.StatusOK
}

func FetchRoleByName(roleName string, ctx context.Context) ([]Role, error) {
	db := getDB()

	var roles []Role

	doc, err := db.Collection("roles").Find(ctx, bson.M{"name": strings.ToLower(roleName)})
	if err != nil {
		return nil, err
	}

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
		{"name", strings.ToLower(crl.Name)},
		{"description", crl.Description},
		{"permission", crl.Permission},
		{"created_by", crl.CreatedBy},
		{"updated_by", crl.UpdatedBy},
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
