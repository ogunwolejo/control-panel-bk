package util

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

// UtilTestSuite defines the test suite for the util package
type UtilTestSuite struct {
	suite.Suite
}

// TestGetBytesResponse_Success tests JSON serialization with valid input
func (suite *UtilTestSuite) TestGetBytesResponse_Success() {
	data := map[string]string{"message": "Success"}
	expectedJSON := `{"Status":200,"Error":null,"Data":{"message":"Success"}}`

	bytes, err := GetBytesResponse(200, data)

	suite.Require().NoError(err, "JSON encoding should not return an error")
	suite.JSONEq(expectedJSON, string(bytes), "JSON output should match expected format")
}

// TestGetBytesResponse_NilData tests JSON serialization with nil data
func (suite *UtilTestSuite) TestGetBytesResponse_NilData() {
	expectedJSON := `{"Status":200,"Error":null,"Data":null}`

	bytes, err := GetBytesResponse(200, nil)

	suite.Require().NoError(err, "JSON encoding should not return an error")
	suite.JSONEq(expectedJSON, string(bytes), "JSON output should match expected format")
}

// TestGetBytesResponse_ComplexData tests JSON serialization with complex data structures
func (suite *UtilTestSuite) TestGetBytesResponse_ComplexData() {
	data := struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Valid bool   `json:"valid"`
	}{
		Name:  "John Doe",
		Age:   30,
		Valid: true,
	}

	bytes, err := GetBytesResponse(200, data)
	suite.Require().NoError(err, "JSON encoding should not return an error")

	var response Response
	err = json.Unmarshal(bytes, &response)
	suite.Require().NoError(err, "Should be able to decode the JSON response")

	suite.Equal(200, response.Status, "Status should match")
	suite.Nil(response.Error, "Error should be nil")
	suite.NotNil(response.Data, "Data should not be nil")
}

// Run the test suite
func TestUtilSuite(t *testing.T) {
	suite.Run(t, new(UtilTestSuite))
}
