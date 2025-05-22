package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"{{.ModuleName}}/app/models"
)

func MigrateAll(db *gorm.DB) error {
	if err := db.AutoMigrate(models.ModelList()...); err != nil {
		return fmt.Errorf("❌ Failed to migrate models: %v", err)
	}
	log.Println("All models migrated successfully.")
	return nil
}

func DropAll(db *gorm.DB) error {
	if err := db.Migrator().DropTable(models.ModelList()...); err != nil {
		return fmt.Errorf("❌ Failed to drop tables: %v", err)
	}
	log.Println("All tables dropped successfully.")
	return nil
}

func DropAllTablesForce(db *gorm.DB) error {
	var tableNames []string
	err := db.Raw(`
		SELECT tablename 
		FROM pg_tables 
		WHERE schemaname = 'public';
	`).Scan(&tableNames).Error
	if err != nil {
		return fmt.Errorf("failed to get table names: %v", err)
	}

	if len(tableNames) == 0 {
		log.Println("No tables found to drop.")
		return nil
	}

	// Disable foreign key constraints supaya bisa drop tanpa error
	if err := db.Exec("SET session_replication_role = 'replica';").Error; err != nil {
		return fmt.Errorf("failed to disable foreign keys: %v", err)
	}

	for _, table := range tableNames {
		if err := db.Migrator().DropTable(table); err != nil {
			log.Printf("Failed to drop table %s: %v", table, err)
		} else {
			log.Printf("Dropped table %s", table)
		}
	}

	// Enable foreign key constraints kembali
	if err := db.Exec("SET session_replication_role = 'origin';").Error; err != nil {
		return fmt.Errorf("failed to enable foreign keys: %v", err)
	}

	return nil
}
