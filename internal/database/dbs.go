// Contain all the various db configurations our application will be connecting to

package database

import (
	"fmt"
	"gorm.io/gorm"
	"os"
)

var (
	SystemUserDB *gorm.DB
)

var (
	Dbs []DbConfig = []DbConfig{
		{
			Dsn:        fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require", os.Getenv("SYS_USER_DB_PGHOST"), os.Getenv("SYS_USER_DB_PGUSER"), os.Getenv("SYS_USER_DB_PGPASSWORD"), os.Getenv("SYS_USER_DB_PGDATABASE"), os.Getenv("SYS_USER_DB_PGPORT")),
			DbInstance: SystemUserDB,
		},
	}
)
