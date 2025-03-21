package panelAdmins

import (
	"testing"
)

func TestFetchRoleById(t *testing.T) {
	//mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	//
	//mt.Run("Test fetch role by id -> success", func(mt *mtest.T) {
	//	mockRoleID := bson.NewObjectID()
	//	expectedRole := Role{
	//		ID:              mockRoleID.Hex(),
	//		Name:            "Test Role",
	//		Description:     "Test Description",
	//		Permission:      Permission{},
	//		CreatedBy:       "tester",
	//		UpdatedBy:       "tester",
	//		ArchiveStatus:   false,
	//		IsDeletedStatus: false,
	//		CreatedAt:       time.Now(),
	//		UpdatedAt:       time.Now(),
	//	}
	//
	//	mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.roles", mtest.FirstBatch, bson.D{
	//		{"_id", mockRoleID},
	//		{"name", expectedRole.Name},
	//		{"description", expectedRole.Description},
	//		{"permission", expectedRole.Permission},
	//		{"created_by", expectedRole.CreatedBy},
	//		{"updated_by", expectedRole.UpdatedBy},
	//		{"archive_status", expectedRole.ArchiveStatus},
	//		{"is_deleted_status", expectedRole.IsDeletedStatus},
	//		{"created_at", expectedRole.CreatedAt},
	//		{"updated_at", expectedRole.UpdatedAt},
	//	}))
	//
	//	req, err := http.NewRequest("GET", fmt.Sprintf("/roles/%s", mockRoleID.Hex()), nil)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//
	//	rr := httptest.NewRecorder()
	//	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	//		roleResp, roleErr, roleCde := FetchRoleById(mockRoleID.Hex(), request.Context(), mt.Client)
	//		if roleErr != nil {
	//			util.ErrorException(writer, roleErr, roleCde)
	//			return
	//		}
	//
	//		rpn := Response{
	//			Status: roleCde,
	//			error:  nil,
	//			Data:   roleResp,
	//		}
	//
	//		rpnBytes, _ := json.Marshal(&rpn)
	//		writer.Header().Set("Content-Type", "application/json")
	//		writer.WriteHeader(roleCde)
	//		writer.Write(rpnBytes)
	//	})
	//
	//	handler.ServeHTTP(rr, req)
	//	if status := rr.Code; status != http.StatusOK {
	//		t.Errorf("expected status %d, got %d", http.StatusOK, status)
	//	}
	//
	//	var testResp Response
	//	err = json.Unmarshal(rr.Body.Bytes(), &testResp)
	//	if err != nil {
	//		require.Error(t, err)
	//	}
	//
	//	require.NoError(t, err)
	//	assert.Equal(t, http.StatusOK, testResp.Status)
	//	assert.Equal(t, expectedRole.ID, testResp.Data.(map[string]interface{})["_id"])
	//	assert.Equal(t, expectedRole.Name, testResp.Data.(map[string]interface{})["name"])
	//})

	//mt.Run("not found", func(mt *mtest.T) {
	//	mockRoleID := primitive.NewObjectID()
	//	mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.roles", mtest.FirstBatch))
	//
	//	_, err, status := FetchRoleById(mockRoleID.Hex(), mt.Context())
	//	require.Error(t, err)
	//	assert.Equal(t, http.StatusOK, status)
	//	assert.Contains(t, err.Error(), "record regarding this role was not found")
	//})
}

//func TestCreateRole(t *testing.T) {
//	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
//
//	mt.Run("success", func(mt *mtest.T) {
//		mockInsertResponse := mtest.CreateSuccessResponse(
//			bson.E{Key: "ok", Value: 1},
//			bson.E{Key: "insertedId", Value: primitive.NewObjectID()},
//		)
//
//		mt.AddMockResponses(mockInsertResponse)
//
//		cr := CRole{
//			Name:        "New Role",
//			Description: "Test Description",
//			Permission:  Permission{},
//			CreatedBy:   "tester",
//			UpdatedBy:   "tester",
//		}
//
//		response, err := CreateRole(cr, mt.Context())
//		require.NoError(t, err)
//		assert.True(t, response.Status)
//		assert.NotNil(t, response.Data)
//	})
//
//	mt.Run("duplicate role", func(mt *mtest.T) {
//		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.roles", mtest.FirstBatch, bson.D{
//			{"name", "existing role"},
//		}))
//
//		cr := CRole{
//			Name: "existing role",
//		}
//
//		_, err := CreateRole(cr, mt.Context())
//		require.Error(t, err)
//		assert.Contains(t, err.Error(), "already exists")
//	})
//}
//
//func TestGeneralizedUpdate(t *testing.T) {
//	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
//
//	mockRoleID := primitive.NewObjectID()
//	role := &Role{
//		ID:        mockRoleID.Hex(),
//		Name:      "Updated Role",
//		UpdatedBy: "tester",
//	}
//
//	mt.Run("success", func(mt *mtest.T) {
//		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.roles", mtest.FirstBatch, bson.D{
//			{"_id", mockRoleID},
//			{"name", role.Name},
//			{"updated_by", role.UpdatedBy},
//		}))
//
//		updatedRole, err, status := role.GeneralizedUpdate(mt.Context())
//		require.NoError(t, err)
//		assert.Equal(t, http.StatusOK, status)
//		assert.Equal(t, role.Name, updatedRole.Name)
//	})
//
//	mt.Run("not found", func(mt *mtest.T) {
//		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.roles", mtest.FirstBatch))
//
//		_, err, status := role.GeneralizedUpdate(mt.Context())
//		require.Error(t, err)
//		assert.Equal(t, http.StatusOK, status)
//		assert.Contains(t, err.Error(), "no role with the selected metrics were found")
//	})
//}
//
//func TestArchiveRole(t *testing.T) {
//	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
//
//	role := &Role{
//		ID:            primitive.NewObjectID().Hex(),
//		ArchiveStatus: false,
//		UpdatedBy:     "tester",
//	}
//
//	mt.Run("success", func(mt *mtest.T) {
//		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.roles", mtest.FirstBatch, bson.D{
//			{"_id", role.ID},
//			{"archive_status", true},
//		}))
//
//		result, err, status := role.ArchiveRole(mt.Context())
//		require.NoError(t, err)
//		assert.Equal(t, http.StatusAccepted, status)
//		assert.True(t, result.ArchiveStatus)
//	})
//
//	mt.Run("already archived", func(mt *mtest.T) {
//		role.ArchiveStatus = true
//		_, err, status := role.ArchiveRole(mt.Context())
//		require.Error(t, err)
//		assert.Equal(t, http.StatusInternalServerError, status)
//		assert.Contains(t, err.Error(), "already archived")
//	})
//}
//
//func TestFetchRoles(t *testing.T) {
//	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
//
//	mt.Run("success with pagination", func(mt *mtest.T) {
//		mockRoles := []bson.D{
//			{{"name", "Role 1"}},
//			{{"name", "Role 2"}},
//		}
//
//		mt.AddMockResponses(
//			mtest.CreateCursorResponse(2, "test.roles", mtest.FirstBatch, mockRoles...),
//			mtest.CreateCursorResponse(0, "test.roles", mtest.NextBatch),
//		)
//
//		roles, err, status := FetchRoles(1, 10, mt.Context())
//		require.NoError(t, err)
//		assert.Equal(t, http.StatusOK, status)
//		assert.Len(t, roles, 2)
//	})
//
//	mt.Run("no roles found", func(mt *mtest.T) {
//		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.roles", mtest.FirstBatch))
//
//		_, err, status := FetchRoles(1, 10, mt.Context())
//		require.Error(t, err)
//		assert.Equal(t, http.StatusOK, status)
//		assert.Contains(t, err.Error(), "no role has been created")
//	})
//}
//
//func TestHardDeleteRole(t *testing.T) {
//	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
//
//	role := &Role{
//		ID: primitive.NewObjectID().Hex(),
//	}
//
//	mt.Run("success", func(mt *mtest.T) {
//		mt.AddMockResponses(bson.D{
//			{"ok", 1},
//			{"n", 1},
//		})
//
//		id, err, status := role.HardDeleteRole(mt.Context())
//		require.NoError(t, err)
//		assert.Equal(t, http.StatusOK, status)
//		assert.Equal(t, role.ID, *id)
//	})
//
//	mt.Run("not found", func(mt *mtest.T) {
//		mt.AddMockResponses(bson.D{
//			{"ok", 1},
//			{"n", 0},
//		})
//
//		_, err, status := role.HardDeleteRole(mt.Context())
//		require.Error(t, err)
//		assert.Equal(t, http.StatusNotImplemented, status)
//	})
//}
