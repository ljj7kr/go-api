package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"go-api/internal/config"
)

func Open(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.MySQLDSN)
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	// 애플리케이션 레벨 기본 풀 설정
	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.DBConnMaxIdleTime)

	return db, nil
}

func Ping(ctx context.Context, db *sql.DB) error {
	// 시작 시점에 연결 가능 여부를 즉시 확인
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping mysql: %w", err)
	}

	return nil
}
