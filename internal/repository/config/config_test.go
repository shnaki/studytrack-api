package config_test

import (
	"testing"

	"github.com/shnaki/studytrack-api/internal/repository/config"
)

func TestLoad_Defaults(t *testing.T) {
	// Ensure relevant env vars are unset so defaults apply.
	t.Setenv("PORT", "")
	t.Setenv("DB_URL", "")
	t.Setenv("CORS_ORIGINS", "")
	t.Setenv("LOG_LEVEL", "")

	cfg := config.Load()

	if cfg.Port != "8080" {
		t.Errorf("expected default Port '8080', got '%s'", cfg.Port)
	}
	if cfg.DBURL != "postgres://studytrack:studytrack@localhost:5432/studytrack?sslmode=disable" {
		t.Errorf("expected default DBURL, got '%s'", cfg.DBURL)
	}
	if len(cfg.CORSOrigins) != 1 || cfg.CORSOrigins[0] != "http://localhost:3000" {
		t.Errorf("expected default CORSOrigins [http://localhost:3000], got %v", cfg.CORSOrigins)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected default LogLevel 'debug', got '%s'", cfg.LogLevel)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	t.Setenv("PORT", "3000")
	t.Setenv("DB_URL", "postgres://custom:pass@db:5432/mydb?sslmode=require")
	t.Setenv("CORS_ORIGINS", "https://example.com,https://app.example.com")
	t.Setenv("LOG_LEVEL", "info")

	cfg := config.Load()

	if cfg.Port != "3000" {
		t.Errorf("expected Port '3000', got '%s'", cfg.Port)
	}
	if cfg.DBURL != "postgres://custom:pass@db:5432/mydb?sslmode=require" {
		t.Errorf("expected custom DBURL, got '%s'", cfg.DBURL)
	}
	if len(cfg.CORSOrigins) != 2 {
		t.Fatalf("expected 2 CORS origins, got %d: %v", len(cfg.CORSOrigins), cfg.CORSOrigins)
	}
	if cfg.CORSOrigins[0] != "https://example.com" {
		t.Errorf("expected first CORS origin 'https://example.com', got '%s'", cfg.CORSOrigins[0])
	}
	if cfg.CORSOrigins[1] != "https://app.example.com" {
		t.Errorf("expected second CORS origin 'https://app.example.com', got '%s'", cfg.CORSOrigins[1])
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected LogLevel 'info', got '%s'", cfg.LogLevel)
	}
}

func TestLoad_SingleCORSOrigin(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DB_URL", "")
	t.Setenv("CORS_ORIGINS", "https://mysite.com")
	t.Setenv("LOG_LEVEL", "")

	cfg := config.Load()

	if len(cfg.CORSOrigins) != 1 {
		t.Fatalf("expected 1 CORS origin, got %d: %v", len(cfg.CORSOrigins), cfg.CORSOrigins)
	}
	if cfg.CORSOrigins[0] != "https://mysite.com" {
		t.Errorf("expected CORS origin 'https://mysite.com', got '%s'", cfg.CORSOrigins[0])
	}
}
