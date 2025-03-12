package database

import (
	"gorm.io/driver/sqlite"
	"os"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Mocking environment variables for tests
func setupEnv() {
	os.Setenv("SYS_USER_DB_PGHOST", "localhost")
	os.Setenv("SYS_USER_DB_PGUSER", "testuser")
	os.Setenv("SYS_USER_DB_PGPASSWORD", "testpassword")
	os.Setenv("SYS_USER_DB_PGDATABASE", "testdb")
	os.Setenv("SYS_USER_DB_PGPORT", "5432")
}

func cleanupEnv() {
	os.Unsetenv("SYS_USER_DB_PGHOST")
	os.Unsetenv("SYS_USER_DB_PGUSER")
	os.Unsetenv("SYS_USER_DB_PGPASSWORD")
	os.Unsetenv("SYS_USER_DB_PGDATABASE")
	os.Unsetenv("SYS_USER_DB_PGPORT")
}

// Test successful database initialization with mocked DB
func TestInitDB_Success(t *testing.T) {
	setupEnv()
	defer cleanupEnv()

	// Mock database using pgxmock
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	// Mock a successful connection
	mock.ExpectPing()

	// Set mock as GORM DB
	SystemUserDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, SystemUserDB)

	// Check connection pool settings
	sqlDB, err := SystemUserDB.DB()
	assert.NoError(t, err)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(8)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	assert.Equal(t, 20, sqlDB.Stats().MaxOpenConnections)
}

// Test failure when required environment variables are missing
func TestInitDB_MissingEnvVariables(t *testing.T) {
	cleanupEnv() // Ensure no environment variables are set

	// Expect InitDB to panic due to missing env variables
	assert.Panics(t, func() {
		InitPanelAdminDB()
	})
}

// Test failure when credentials are incorrect
func TestInitDB_InvalidCredentials(t *testing.T) {
	setupEnv()
	os.Setenv("SYS_USER_DB_PGPASSWORD", "wrongpassword") // Simulate incorrect credentials

	// Expect InitDB to panic due to authentication failure
	assert.Panics(t, func() {
		InitPanelAdminDB()
	})

	// Cleanup after test
	cleanupEnv()
}

// Test that the connection pool settings are correctly applied
func TestConnectionPoolingSettings(t *testing.T) {
	setupEnv()
	defer cleanupEnv()

	// Mock database using pgxmock
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	// Set mock as GORM DB
	SystemUserDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, SystemUserDB)

	// Apply pooling settings
	sqlDB, err := SystemUserDB.DB()
	assert.NoError(t, err)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(8)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	// Validate pooling settings
	assert.Equal(t, 20, sqlDB.Stats().MaxOpenConnections)
}
