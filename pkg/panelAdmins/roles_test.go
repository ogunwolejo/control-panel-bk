package panelAdmins

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RoleTestSuite struct {
	suite.Suite
	db     *mongo.Database
	client *mongo.Client
	ctx    context.Context
}

func (suite *RoleTestSuite) SetupSuite() {
	// Connect to test MongoDB
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.client = client
	suite.db = client.Database("test_roles_db")
	suite.ctx = context.Background()
}

func (suite *RoleTestSuite) SetupTest() {
	// Clear roles collection before each test
	_, err := suite.db.Collection("roles").DeleteMany(suite.ctx, bson.M{})
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *RoleTestSuite) TearDownSuite() {
	// Cleanup after all tests
	suite.db.Drop(suite.ctx)
	suite.client.Disconnect(suite.ctx)
}

func TestRoleTestSuite(t *testing.T) {
	suite.Run(t, new(RoleTestSuite))
}

func (suite *RoleTestSuite) TestCreateRole_Success() {
	cRole := CRole{
		Name:        "test-role",
		Description: "Test description",
		Permission:  Permission{}, // Add proper permission data
		CreatedBy:   "tester",
		UpdatedBy:   "tester",
	}

	resp, err := CreateRole(cRole, suite.ctx, suite.db)

	suite.NoError(err)
	suite.True(resp.Status)
	suite.Equal("Role has been created", resp.Message)

	// Verify database entry
	var role Role
	err = suite.db.Collection("roles").FindOne(suite.ctx, bson.M{"name": "test-role"}).Decode(&role)
	suite.NoError(err)
	suite.Equal(cRole.Name, role.Name)
}

func (suite *RoleTestSuite) TestCreateRole_Duplicate() {
	// Create initial role
	cRole := CRole{Name: "duplicate-role"}
	_, err := CreateRole(cRole, suite.ctx, suite.db)
	suite.NoError(err)

	// Try to create duplicate
	_, err = CreateRole(cRole, suite.ctx, suite.db)
	suite.Error(err)
	suite.Equal("a role having the same name already exists", err.Error())
}

func (suite *RoleTestSuite) TestArchiveRole_Success() {
	// Insert test role
	role := Role{
		ID:            bson.NewObjectID().Hex(),
		Name:          "archivable-role",
		ArchiveStatus: false,
		UpdatedBy:     "tester",
	}
	objID, _ := bson.ObjectIDFromHex(role.ID)
	_, err := suite.db.Collection("roles").InsertOne(suite.ctx, bson.M{
		"_id":            objID,
		"name":           role.Name,
		"archive_status": role.ArchiveStatus,
		"updated_by":     role.UpdatedBy,
	})
	suite.NoError(err)

	// Archive the role
	updatedRole, err, code := role.ArchiveRole(suite.ctx, suite.db)

	suite.NoError(err)
	suite.Equal(http.StatusAccepted, code)
	suite.True(updatedRole.ArchiveStatus)
}

func (suite *RoleTestSuite) TestFetchRoles_Pagination() {
	// Insert multiple test roles
	for i := 0; i < 15; i++ {
		role := CRole{
			Name:      fmt.Sprintf("role-%d", i),
			CreatedBy: "tester",
			UpdatedBy: "tester",
		}
		_, err := CreateRole(role, suite.ctx, suite.db)
		suite.NoError(err)
	}

	// Fetch first page
	roles, err, code := FetchRoles(1, 10, suite.ctx, suite.db)
	suite.NoError(err)
	suite.Equal(http.StatusOK, code)
	suite.Len(roles, 10)

	// Fetch second page
	roles, err, code = FetchRoles(2, 10, suite.ctx, suite.db)
	suite.NoError(err)
	suite.Equal(http.StatusOK, code)
	suite.Len(roles, 5)
}

func (suite *RoleTestSuite) TestGeneralizedUpdate_Success() {
	// Create test role
	role := Role{
		ID:        bson.NewObjectID().Hex(),
		Name:      "original-name",
		UpdatedBy: "tester",
	}
	objID, _ := bson.ObjectIDFromHex(role.ID)
	_, err := suite.db.Collection("roles").InsertOne(suite.ctx, bson.M{
		"_id":        objID,
		"name":       role.Name,
		"updated_by": role.UpdatedBy,
	})
	suite.NoError(err)

	// Update the role
	role.Name = "updated-name"
	role.Description = "new description"
	updatedRole, err, code := role.GeneralizedUpdate(suite.ctx, suite.db)

	suite.NoError(err)
	suite.Equal(http.StatusOK, code)
	suite.Equal("updated-name", updatedRole.Name)
	suite.Equal("new description", updatedRole.Description)
}

func (suite *RoleTestSuite) TestHandleCreateRole_HTTP() {
	handler := HandleCreateRole(suite.db)

	// Create request
	cRole := CRole{
		Name:      "http-test-role",
		CreatedBy: "tester",
		UpdatedBy: "tester",
	}
	body, _ := json.Marshal(cRole)
	req := httptest.NewRequest("POST", "/roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	suite.Equal(http.StatusCreated, w.Code)

	// Verify response
	var response CreateRoleResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.True(response.Status)

	// Verify database
	var role Role
	err = suite.db.Collection("roles").FindOne(suite.ctx, bson.M{"name": "http-test-role"}).Decode(&role)
	suite.NoError(err)
	suite.Equal(cRole.Name, role.Name)
}

func (suite *RoleTestSuite) TestHardDeleteRole_Success() {
	// Insert test role
	role := Role{ID: bson.NewObjectID().Hex()}
	objID, _ := bson.ObjectIDFromHex(role.ID)
	_, err := suite.db.Collection("roles").InsertOne(suite.ctx, bson.M{"_id": objID})
	suite.NoError(err)

	// Delete the role
	deletedID, err, code := role.HardDeleteRole(suite.ctx, suite.db)

	suite.NoError(err)
	suite.Equal(http.StatusOK, code)
	suite.Equal(role.ID, *deletedID)

	// Verify deletion
	err = suite.db.Collection("roles").FindOne(suite.ctx, bson.M{"_id": objID}).Err()
	suite.Equal(mongo.ErrNoDocuments, err)
}

func (suite *RoleTestSuite) TestFetchRoleById_NotFound() {
	nonExistentID := bson.NewObjectID().Hex()
	_, err, code := FetchRoleById(nonExistentID, suite.ctx, suite.db)

	suite.Error(err)
	suite.Contains(err.Error(), "record regarding this role was not found")
	suite.Equal(http.StatusOK, code)
}

func (suite *RoleTestSuite) TestUnArchiveRole_Success() {
	// Create archived role
	role := Role{
		ID:            primitive.NewObjectID().Hex(),
		Name:          "archived-role",
		ArchiveStatus: true,
		UpdatedBy:     "tester",
	}
	objID, _ := primitive.ObjectIDFromHex(role.ID)
	_, err := suite.db.Collection("roles").InsertOne(suite.ctx, bson.M{
		"_id":            objID,
		"name":           role.Name,
		"archive_status": true,
		"updated_by":     role.UpdatedBy,
	})
	suite.NoError(err)

	// Unarchive the role
	updatedRole, err, code := role.UnArchiveRole(suite.ctx, suite.db)

	suite.NoError(err)
	suite.Equal(http.StatusOK, code)
	suite.False(updatedRole.ArchiveStatus)

	// Verify database
	var dbRole Role
	err = suite.db.Collection("roles").FindOne(suite.ctx, bson.M{"_id": objID}).Decode(&dbRole)
	suite.NoError(err)
	suite.False(dbRole.ArchiveStatus)
}

func (suite *RoleTestSuite) TestUnArchiveRole_NotArchived() {
	role := Role{
		ID:            primitive.NewObjectID().Hex(),
		Name:          "active-role",
		ArchiveStatus: false,
	}
	objID, _ := primitive.ObjectIDFromHex(role.ID)
	_, err := suite.db.Collection("roles").InsertOne(suite.ctx, bson.M{
		"_id":            objID,
		"name":           role.Name,
		"archive_status": false,
	})
	suite.NoError(err)

	_, err, code := role.UnArchiveRole(suite.ctx, suite.db)

	suite.Error(err)
	suite.Contains(err.Error(), "is not archived")
	suite.Equal(http.StatusBadRequest, code)
}

func (suite *RoleTestSuite) TestArchiveRole_AlreadyArchived() {
	role := Role{
		ID:            primitive.NewObjectID().Hex(),
		Name:          "already-archived",
		ArchiveStatus: true,
	}
	objID, _ := primitive.ObjectIDFromHex(role.ID)
	_, err := suite.db.Collection("roles").InsertOne(suite.ctx, bson.M{
		"_id":            objID,
		"name":           role.Name,
		"archive_status": true,
	})
	suite.NoError(err)

	_, err, code := role.ArchiveRole(suite.ctx, suite.db)

	suite.Error(err)
	suite.Equal("role is already archived", err.Error())
	suite.Equal(http.StatusInternalServerError, code)
}

func (suite *RoleTestSuite) TestPushRoleToBin_Success() {
	role := Role{
		ID:              primitive.NewObjectID().Hex(),
		Name:            "to-delete",
		IsDeletedStatus: false,
		UpdatedBy:       "tester",
	}
	objID, _ := primitive.ObjectIDFromHex(role.ID)
	_, err := suite.db.Collection("roles").InsertOne(suite.ctx, bson.M{
		"_id":               objID,
		"name":              role.Name,
		"is_deleted_status": false,
	})
	suite.NoError(err)

	updatedRole, err, code := role.PushRoleToBin(suite.ctx, suite.db)

	suite.NoError(err)
	suite.Equal(http.StatusOK, code)
	suite.True(updatedRole.IsDeletedStatus)

	// Verify database
	var dbRole Role
	err = suite.db.Collection("roles").FindOne(suite.ctx, bson.M{"_id": objID}).Decode(&dbRole)
	suite.NoError(err)
	suite.True(dbRole.IsDeletedStatus)
}

func (suite *RoleTestSuite) TestPushRoleToBin_AlreadyDeleted() {
	role := Role{
		ID:              primitive.NewObjectID().Hex(),
		IsDeletedStatus: true,
	}
	objID, _ := primitive.ObjectIDFromHex(role.ID)
	_, err := suite.db.Collection("roles").InsertOne(suite.ctx, bson.M{
		"_id":               objID,
		"is_deleted_status": true,
	})
	suite.NoError(err)

	_, err, code := role.PushRoleToBin(suite.ctx, suite.db)

	suite.Error(err)
	suite.Contains(err.Error(), "has been sent to the bin")
	suite.Equal(http.StatusOK, code)
}

func (suite *RoleTestSuite) TestRestoreRoleFromBin_Success() {
	role := Role{
		ID:              primitive.NewObjectID().Hex(),
		IsDeletedStatus: true,
		UpdatedBy:       "tester",
	}
	objID, _ := primitive.ObjectIDFromHex(role.ID)
	_, err := suite.db.Collection("roles").InsertOne(suite.ctx, bson.M{
		"_id":               objID,
		"is_deleted_status": true,
	})
	suite.NoError(err)

	updatedRole, err, code := role.RestoreRoleFromBin(suite.ctx, suite.db)

	suite.NoError(err)
	suite.Equal(http.StatusOK, code)
	suite.False(updatedRole.IsDeletedStatus)

	// Verify database
	var dbRole Role
	err = suite.db.Collection("roles").FindOne(suite.ctx, bson.M{"_id": objID}).Decode(&dbRole)
	suite.NoError(err)
	suite.False(dbRole.IsDeletedStatus)
}

func (suite *RoleTestSuite) TestRestoreRoleFromBin_NotInBin() {
	role := Role{
		ID:              primitive.NewObjectID().Hex(),
		IsDeletedStatus: false,
	}
	objID, _ := primitive.ObjectIDFromHex(role.ID)
	_, err := suite.db.Collection("roles").InsertOne(suite.ctx, bson.M{
		"_id":               objID,
		"is_deleted_status": false,
	})
	suite.NoError(err)

	_, err, code := role.RestoreRoleFromBin(suite.ctx, suite.db)

	suite.Error(err)
	suite.Contains(err.Error(), "not in the bin catalogue")
	suite.Equal(http.StatusOK, code)
}

func (suite *RoleTestSuite) TestHandleArchiveRole_HTTP() {
	// Setup test role
	role := Role{
		ID:            primitive.NewObjectID().Hex(),
		Name:          "http-archive-test",
		ArchiveStatus: false,
		UpdatedBy:     "tester",
	}
	objID, _ := primitive.ObjectIDFromHex(role.ID)
	_, err := suite.db.Collection("roles").InsertOne(suite.ctx, bson.M{
		"_id":            objID,
		"name":           role.Name,
		"archive_status": false,
	})
	suite.NoError(err)

	// Create request
	handler := HandleArchiveRole(suite.db)
	body, _ := json.Marshal(role)
	req := httptest.NewRequest("POST", "/roles/archive", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	suite.Equal(http.StatusAccepted, w.Code)

	// Verify response
	var response Role
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.True(response.ArchiveStatus)
}

func (suite *RoleTestSuite) TestHandleRestoreFromBin_HTTP() {
	// Setup test role in bin
	role := Role{
		ID:              primitive.NewObjectID().Hex(),
		IsDeletedStatus: true,
		UpdatedBy:       "tester",
	}
	objID, _ := primitive.ObjectIDFromHex(role.ID)
	_, err := suite.db.Collection("roles").InsertOne(suite.ctx, bson.M{
		"_id":               objID,
		"is_deleted_status": true,
	})
	suite.NoError(err)

	handler := HandleRestoreRoleFromBin(suite.db)
	body, _ := json.Marshal(role)
	req := httptest.NewRequest("POST", "/roles/restore", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	// Verify response
	var response Role
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.False(response.IsDeletedStatus)
}

func (suite *RoleTestSuite) TestHandleCreateRole_Success() {
	handler := HandleCreateRole(suite.db)

	cRole := CRole{
		Name:      "test-handler-role",
		CreatedBy: "tester",
		UpdatedBy: "tester",
	}
	body, _ := json.Marshal(cRole)

	req := httptest.NewRequest("POST", "/roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	suite.Equal(http.StatusCreated, w.Code)

	var response CreateRoleResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.True(response.Status)
}

func (suite *RoleTestSuite) TestHandleCreateRole_Duplicate() {
	// Create initial role
	cRole := CRole{Name: "duplicate-handler-role"}
	_, err := CreateRole(cRole, suite.ctx, suite.db)
	suite.NoError(err)

	handler := HandleCreateRole(suite.db)
	body, _ := json.Marshal(cRole)
	req := httptest.NewRequest("POST", "/roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.False(response["status"].(bool))
}

func (suite *RoleTestSuite) TestHandleFetchRoleByName_Success() {
	// Create test roles
	roles := []CRole{
		{Name: "fetch-test-role"},
		{Name: "fetch-test-role"},
	}
	for _, r := range roles {
		_, err := CreateRole(r, suite.ctx, suite.db)
		suite.NoError(err)
	}

	handler := HandleFetchRoleByName(suite.db)
	req := httptest.NewRequest("GET", "/roles?name=fetch-test-role", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var result []Role
	err := json.Unmarshal(w.Body.Bytes(), &result)
	suite.NoError(err)
	suite.Len(result, 2)
}

func (suite *RoleTestSuite) TestHandleFetchRoleById_Success() {
	// Create test role
	role := CRole{Name: "fetch-by-id-role"}
	created, err := CreateRole(role, suite.ctx, suite.db)
	suite.NoError(err)

	handler := HandleFetchRoleById(suite.db)
	req := httptest.NewRequest("GET", "/roles/"+created.Data.InsertedID.(primitive.ObjectID).Hex(), nil)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", created.Data.InsertedID.(primitive.ObjectID).Hex())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var result Role
	err = json.Unmarshal(w.Body.Bytes(), &result)
	suite.NoError(err)
	suite.Equal("fetch-by-id-role", result.Name)
}

func (suite *RoleTestSuite) TestHandleFetchRoles_Pagination() {
	// Create 15 test roles
	for i := 0; i < 15; i++ {
		role := CRole{Name: fmt.Sprintf("page-role-%d", i)}
		_, err := CreateRole(role, suite.ctx, suite.db)
		suite.NoError(err)
	}

	handler := HandleFetchRoles(suite.db)
	req := httptest.NewRequest("GET", "/roles?page=2&limit=10", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var result []Role
	err := json.Unmarshal(w.Body.Bytes(), &result)
	suite.NoError(err)
	suite.Len(result, 5)
}

func (suite *RoleTestSuite) TestHandleHardDeleteOfRole_Success() {
	// Create test role
	role := CRole{Name: "hard-delete-role"}
	created, err := CreateRole(role, suite.ctx, suite.db)
	suite.NoError(err)

	handler := HandleHardDeleteOfRole(suite.db)
	body, _ := json.Marshal(map[string]string{"_id": created.Data.InsertedID.(bson.ObjectID).Hex()})
	req := httptest.NewRequest("DELETE", "/roles", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	// Verify deletion
	_, err, _ = FetchRoleById(created.Data.InsertedID.(bson.ObjectID).Hex(), suite.ctx, suite.db)
	suite.Error(err)
}

func (suite *RoleTestSuite) TestHandleGeneralUpdate_Success() {
	// Create test role
	role := CRole{Name: "update-test-role"}
	created, err := CreateRole(role, suite.ctx, suite.db)
	suite.NoError(err)

	handler := HandleGeneralUpdate(suite.db)
	updateData := Role{
		ID:        created.Data.InsertedID.(bson.ObjectID).Hex(),
		Name:      "updated-name",
		UpdatedBy: "tester",
	}
	body, _ := json.Marshal(updateData)

	req := httptest.NewRequest("PUT", "/roles", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	suite.Equal(http.StatusAccepted, w.Code)

	// Verify update
	updated, err, _ := FetchRoleById(updateData.ID, suite.ctx, suite.db)
	suite.NoError(err)
	suite.Equal("updated-name", updated.Name)
}

func (suite *RoleTestSuite) TestHandlePushRoleToBin_HTTP() {
	// Create test role
	role := CRole{Name: "bin-test-role"}
	created, err := CreateRole(role, suite.ctx, suite.db)
	suite.NoError(err)

	handler := HandlePushRoleToBin(suite.db)
	body, _ := json.Marshal(map[string]string{"_id": created.Data.InsertedID.(bson.ObjectID).Hex()})
	req := httptest.NewRequest("POST", "/roles/bin", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	// Verify bin status
	inBin, err, _ := FetchRoleById(created.Data.InsertedID.(bson.ObjectID).Hex(), suite.ctx, suite.db)
	suite.NoError(err)
	suite.True(inBin.IsDeletedStatus)
}
