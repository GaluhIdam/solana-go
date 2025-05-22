package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"note-api/app/models"
	"note-api/core/config"
)

var DB *gorm.DB

func ConnectDB(configPath string) (*gorm.DB, error) {
	err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}
	dbConfig := config.GlobalConfig.Database

	if dbConfig.Host == "" || dbConfig.User == "" || dbConfig.Password == "" || dbConfig.Name == "" || dbConfig.Port == "" {
		log.Println("⏩ Skipping database connection (no database config found in YAML)")
		return nil, nil
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbConfig.Host, dbConfig.User, dbConfig.Password,
		dbConfig.Name, dbConfig.Port, dbConfig.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to connect to DB: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to get raw DB connection: %v", err)
	}
	if dbConfig.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	}
	if dbConfig.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	}
	if dbConfig.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)
	}

	DB = db

	if err := db.AutoMigrate(models.ModelList()...); err != nil {
		return nil, fmt.Errorf("❌ Auto migration failed: %v", err)
	}

	switch dbConfig.DDLMode {
	case "create":
		log.Println("🧹 Dropping all tables by models...")
		if err := DropAll(DB); err != nil {
			return nil, err
		}
		fallthrough
	case "create-drop":
		log.Println("🧹 Dropping all tables...")
		if err := DropAllTablesForce(DB); err != nil {
			return nil, err
		}
		fallthrough
	case "update":
		log.Println("📄 Running AutoMigrate...")
		if err := MigrateAll(DB); err != nil {
			return nil, err
		}
		if dbConfig.Reset {
			log.Println("⚠️  Reset enabled: deleting all records in tables...")
			for _, model := range models.ModelList() {
				if err := DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(model).Error; err != nil {
					log.Printf("❌ Failed to reset table for %T: %v\n", model, err)
				}
			}
			log.Println("✅ All tables reset (data deleted)")
		}
	case "none":
		log.Println("⏩ Skipping DB migration (ddl-mode: none)")
	default:
		log.Fatalf("❌ Unknown ddl-mode: %s", dbConfig.DDLMode)
	}

	if dbConfig.GCInterval > 0 {
		go runGarbageCollector(db, dbConfig.GCInterval)
	}

	log.Println("📦 Database connected successfully")
	log.Println("✅ Database ready")
	return db, nil
}

func runGarbageCollector(_ *gorm.DB, interval time.Duration) {
	log.Println("🧹 Garbage Collector started with interval:", interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("🧹 Running periodic garbage collection tasks...")
	}
}
