package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"sync"
	"time"
)

type DbConfig struct {
	Dsn        string
	DbInstance *gorm.DB
}

func connectDb(dns string) (*gorm.DB, error) {
	// Define DB connection parameters
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require", os.Getenv("SYS_USER_DB_PGHOST"), os.Getenv("SYS_USER_DB_PGUSER"), os.Getenv("SYS_USER_DB_PGPASSWORD"), os.Getenv("SYS_USER_DB_PGDATABASE"), os.Getenv("SYS_USER_DB_PGPORT"))

	// Open GORM connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Get underlying sql.DB instance
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(20)                 // Max open connections (adjust based on workload)
	sqlDB.SetMaxIdleConns(8)                  // Max idle connections (reduce idle resource usage)
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Max lifetime of a connection (prevents stale connections)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute) // Max idle time before closing connection

	// Ping DB to ensure connectivity
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %v", err)
	}

	fmt.Println("Database connection initialized successfully with connection pooling!")
	return db, nil
}

func InitializeAllDbs(dbConfigs []DbConfig) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	wg.Add(len(dbConfigs))
	for _, dbConfig := range dbConfigs {

		go func(dsn string, dbPointer *gorm.DB) {
			defer wg.Done()

			db, err := connectDb(dsn)
			if err == nil {
				mutex.Lock()
				dbPointer = db
				mutex.Unlock()
			} else {
				log.Fatalf("Failed to connect to dsn %s, dbInstance %#v", dsn, dbPointer)
			}
		}(dbConfig.Dsn, dbConfig.DbInstance)
	}

	wg.Wait()
}
