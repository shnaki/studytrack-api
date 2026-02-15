// Package controller contains HTTP handlers and routing logic.
package controller

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/shnaki/studytrack-api/internal/usecase"
)

// Usecases holds all application usecases.
type Usecases struct {
	User     *usecase.UserUsecase
	Project  *usecase.ProjectUsecase
	StudyLog *usecase.StudyLogUsecase
	Goal     *usecase.GoalUsecase
	Stats    *usecase.StatsUsecase
	Note     *usecase.NoteUsecase
}

// NewRouter creates and configures the main HTTP router.
func NewRouter(usecases *Usecases, corsOrigins []string, logger *slog.Logger) http.Handler {
	router := chi.NewMux()

	// Middleware
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(requestLogger(logger))
	router.Use(chimiddleware.Recoverer)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   corsOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Versioned API sub-router
	v1 := chi.NewRouter()
	router.Mount("/v1", v1)

	// Huma API on /v1
	config := huma.DefaultConfig("StudyTrack API", "1.0.0")
	config.Info.Description = "Learning progress tracking REST API"
	config.Servers = []*huma.Server{{URL: "/v1"}}
	api := humachi.New(v1, config)

	// Register routes (paths are relative to /v1)
	RegisterUserRoutes(api, usecases.User)
	RegisterProjectRoutes(api, usecases.Project)
	RegisterStudyLogRoutes(api, usecases.StudyLog)
	RegisterGoalRoutes(api, usecases.Goal)
	RegisterStatsRoutes(api, usecases.Stats)
	RegisterNoteRoutes(api, usecases.Note)

	return router
}

func requestLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			logger.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"duration_ms", time.Since(start).Milliseconds(),
				"bytes", ww.BytesWritten(),
			)
		})
	}
}
