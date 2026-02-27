package db

import (
	"fmt"
	"strings"
	"time"

	"release-platform/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(databaseDriver, databaseDSN string) (*gorm.DB, error) {
	dialector, err := buildDialector(databaseDriver, databaseDSN)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}
	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := migrate(db); err != nil {
		return nil, err
	}

	if err := seedDefaultPolicy(db); err != nil {
		return nil, err
	}

	return db, nil
}

func buildDialector(driver, dsn string) (gorm.Dialector, error) {
	switch strings.ToLower(strings.TrimSpace(driver)) {
	case "postgres", "postgresql":
		return postgres.Open(dsn), nil
	case "sqlite":
		return sqlite.Open(dsn), nil
	default:
		return nil, fmt.Errorf("unsupported database driver %q", driver)
	}
}

func migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.Application{},
		&model.ReleaseTicket{},
		&model.Artifact{},
		&model.ScanReport{},
		&model.SQLIssue{},
		&model.ConfigSnapshot{},
		&model.ConfigDiff{},
		&model.ApprovalRecord{},
		&model.PolicyVersion{},
		&model.AuditLog{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return nil
}

func seedDefaultPolicy(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.PolicyVersion{}).Where("version = ?", "v1-default").Count(&count).Error; err != nil {
		return fmt.Errorf("count default policy: %w", err)
	}
	if count > 0 {
		return nil
	}

	rules := []byte(`{"high_sql_action":"block","critical_config_action":"secondary_approval","medium_sql_action":"warn","low_sql_action":"pass"}`)
	policy := model.PolicyVersion{
		Environment: "default",
		Version:     "v1-default",
		Rules:       rules,
		IsActive:    true,
		CreatedBy:   "system",
	}
	if err := db.Create(&policy).Error; err != nil {
		return fmt.Errorf("seed default policy: %w", err)
	}
	return nil
}
