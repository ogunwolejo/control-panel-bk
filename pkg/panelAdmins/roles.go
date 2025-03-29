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
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CreateRoleResponse struct {
	Status  bool
	Message string
	Data    *mongo.InsertOneResult
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

func (rl *Role) GeneralizedUpdate(ctx context.Context, client *mongo.Database) (*Role, error, int) {
	db := client

	objId, objErr := util.GetPrimitiveID(rl.ID)
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

func (rl *Role) UnArchiveRole(ctx context.Context, client *mongo.Database) (*Role, error, int) {
	db := client

	if !rl.ArchiveStatus {
		return nil, fmt.Errorf("role %s is not archived, hence this task cannot be performed", rl.Name), http.StatusBadRequest
	}

	objId, objErr := util.GetPrimitiveID(rl.ID)
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

func (rl *Role) ArchiveRole(ctx context.Context, client *mongo.Database) (*Role, error, int) {
	if rl.ArchiveStatus {
		return nil, errors.New("role is already archived"), http.StatusInternalServerError
	}

	db := client

	objId, objErr := util.GetPrimitiveID(rl.ID)
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

func (rl *Role) PushRoleToBin(ctx context.Context, client *mongo.Database) (*Role, error, int) {
	db := client

	if rl.IsDeletedStatus {
		return nil, fmt.Errorf("role has been sent to the bin"), http.StatusOK
	}

	objId, objErr := util.GetPrimitiveID(rl.ID)
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

func (rl *Role) RestoreRoleFromBin(ctx context.Context, client *mongo.Database) (*Role, error, int) {
	db := client

	if !rl.IsDeletedStatus {
		return nil, fmt.Errorf("role is not in the bin catalogue"), http.StatusOK
	}

	objId, objErr := util.GetPrimitiveID(rl.ID)
	if objErr != nil {
		return nil, objErr, http.StatusInternalServerError
	}

	filter := bson.D{{"_id", objId}}
	update := bson.M{
		"$set": bson.M{
			"updated_at":        time.Now().UTC(),
			"is_deleted_status": false,
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

func (rl *Role) HardDeleteRole(ctx context.Context, client *mongo.Database) (*string, error, int) {
	db := client

	objId, objErr := util.GetPrimitiveID(rl.ID)
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

func FetchRoles(page int, limit int, ctx context.Context, client *mongo.Database) ([]Role, error, int) {
	var roles []Role

	db := client

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

func FetchRoleById(roleId string, ctx context.Context, client *mongo.Database) (*Role, error, int) {
	db := client

	objID, err := util.GetPrimitiveID(roleId)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	var role Role

	flt := bson.M{"_id": objID}
	err = db.Collection("roles").FindOne(ctx, flt).Decode(&role)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("record regarding this role was not found"), http.StatusOK
		}

		return nil, err, http.StatusNotFound
	}

	return &role, nil, http.StatusOK
}

func FetchRoleByName(roleName string, ctx context.Context, client *mongo.Database) ([]Role, error) {
	db := client

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

func CreateRole(crl CRole, ctx context.Context, client *mongo.Database) (*CreateRoleResponse, error) {
	result, err := FetchRoleByName(crl.Name, ctx, client)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		return nil, errors.New("a role having the same name already exists")
	}

	db := client
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

// Handlers

func HandleCreateRole(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var body CRole
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

		output, outputErr := CreateRole(body, r.Context(), db)
		if outputErr != nil {
			if errors.Is(outputErr, errors.New("a role having the same name already exists")) {
				util.ErrorException(w, outputErr, http.StatusOK)
				return
			}

			util.ErrorException(w, outputErr, http.StatusInternalServerError)
			return
		}

		respBytes, e := util.GetBytesResponse(http.StatusCreated, output)
		if e != nil {
			util.ErrorException(w, e, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(respBytes); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

	}
}

func HandleFetchRoleByName(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roleName := r.URL.Query().Get("name")
		roles, err := FetchRoleByName(roleName, r.Context(), db)

		if err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

		respBytes, e := util.GetBytesResponse(http.StatusOK, roles)
		if e != nil {
			util.ErrorException(w, e, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(respBytes); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

	}
}

func HandleFetchRoleById(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roleId := chi.URLParam(r, "id")
		result, err, cde := FetchRoleById(roleId, r.Context(), db)

		if err != nil {
			util.ErrorException(w, err, cde)
			return
		}

		resp, respErr := json.Marshal(&result)
		if respErr != nil {
			util.ErrorException(w, respErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(resp); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
		}
	}
}

func HandleFetchRoles(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		var limit int
		var page int

		if len(query.Get("page")) > 0 {
			pg, pgErr := strconv.Atoi(query.Get("page"))
			if pgErr != nil {
				util.ErrorException(w, pgErr, http.StatusInternalServerError)
				return
			}
			page = pg
		} else {
			page = 1
		}

		if len(query.Get("limit")) > 0 {
			lmt, lmtErr := strconv.Atoi(query.Get("limit"))
			if lmtErr != nil {
				util.ErrorException(w, lmtErr, http.StatusInternalServerError)
				return
			}
			limit = lmt
		} else {
			limit = MAX_LIMIT
		}

		result, err, code := FetchRoles(page, limit, r.Context(), db)

		if err != nil {
			util.ErrorException(w, err, code)
			return
		}

		if respBytes, err := util.GetBytesResponse(code, result); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(respBytes); err != nil {
				util.ErrorException(w, err, http.StatusInternalServerError)
			}
		}

	}
}

func HandleHardDeleteOfRole(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var role Role

		if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

		if id, err, code := role.HardDeleteRole(r.Context(), db); err != nil {
			util.ErrorException(w, err, code)
			return
		} else {
			deleteBytes, delErr := util.GetBytesResponse(code, id)
			if delErr != nil {
				util.ErrorException(w, delErr, http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			if _, err := w.Write(deleteBytes); err != nil {
				util.ErrorException(w, err, http.StatusInternalServerError)
			}
		}
	}
}

func HandleGeneralUpdate(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var body Role

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

		updateDoc, updateError, code := body.GeneralizedUpdate(r.Context(), db)
		if updateError != nil {
			util.ErrorException(w, updateError, code)
			return
		}

		respBytes, respErr := util.GetBytesResponse(http.StatusAccepted, updateDoc)
		if respErr != nil {
			util.ErrorException(w, respErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		if _, err := w.Write(respBytes); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
		}
	}
}

func HandleArchiveRole(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var role Role
		if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

		doc, docErr, code := role.ArchiveRole(r.Context(), db)
		if docErr != nil {
			if errors.Is(docErr, errors.New("no document was found")) {
				util.ErrorException(w, docErr, code)
				return
			}

			util.ErrorException(w, docErr, code)
			return
		}

		archBytes, archErr := json.Marshal(&doc)
		if archErr != nil {
			util.ErrorException(w, archErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		if _, err := w.Write(archBytes); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
		}

	}
}

func HandleUnArchiveRole(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var role Role
		if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

		doc, docErr, code := role.UnArchiveRole(r.Context(), db)
		if docErr != nil {
			if errors.Is(docErr, errors.New("no document was found")) {
				util.ErrorException(w, docErr, code)
				return
			}

			util.ErrorException(w, docErr, code)
			return
		}

		unArchBytes, unArchErr := util.GetBytesResponse(code, doc)
		if unArchErr != nil {
			util.ErrorException(w, unArchErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		if _, err := w.Write(unArchBytes); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
		}
	}
}

func HandlePushRoleToBin(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var role Role
		if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

		bin, binErr, code := role.PushRoleToBin(r.Context(), db)
		if binErr != nil {
			if errors.Is(binErr, errors.New("no document was found")) {
				util.ErrorException(w, binErr, code)
				return
			}

			util.ErrorException(w, binErr, code)
			return
		}

		binByte, bbErr := util.GetBytesResponse(code, bin)
		if bbErr != nil {
			util.ErrorException(w, bbErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		if _, err := w.Write(binByte); err != nil {
			util.ErrorException(w, bbErr, http.StatusInternalServerError)
		}
	}
}

func HandleRestoreRoleFromBin(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var role Role
		if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
			return
		}

		bin, binErr, code := role.RestoreRoleFromBin(r.Context(), db)
		if binErr != nil {
			if errors.Is(binErr, errors.New("no document was found")) {
				util.ErrorException(w, binErr, code)
				return
			}

			util.ErrorException(w, binErr, code)
			return
		}

		binByte, bbErr := util.GetBytesResponse(code, bin)
		if bbErr != nil {
			util.ErrorException(w, bbErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(binByte)
	}
}
