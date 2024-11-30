package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

var GlobalDB *gorm.DB

func InitDatabase() (err error) {
	// TODO: use secure pattern
	// Setup Vault and use vault-injector
	os.Setenv("DB_USERNAME", "aichatworkspace")
	os.Setenv("DB_PASSWORD", "aichatworkspace")
	os.Setenv("DATABASE_HOST", "aichat-workspace-operator-mysql.aichat-workspace-operator-system.svc.cluster.local")
	os.Setenv("DB_DATABASE", "aichatworkspace")

	dsn := fmt.Sprintf(
		"%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DB_DATABASE"),
	)

	GlobalDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	return
}
