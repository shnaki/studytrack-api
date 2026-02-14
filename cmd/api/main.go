package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shnaki/studytrack-api/internal/controller"
	"github.com/shnaki/studytrack-api/internal/repository/config"
	"github.com/shnaki/studytrack-api/internal/repository/postgres"
	"github.com/shnaki/studytrack-api/internal/usecase"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	cfg := config.Load()

	// Run migrations
	logger.Info("running migrations")
	if err := postgres.RunMigrations(cfg.DBURL, "db/migrations"); err != nil {
		logger.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Connect to DB
	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, cfg.DBURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	logger.Info("connected to database")

	// Repositories
	userRepo := postgres.NewUserRepository(pool)
	subjectRepo := postgres.NewSubjectRepository(pool)
	studyLogRepo := postgres.NewStudyLogRepository(pool)
	goalRepo := postgres.NewGoalRepository(pool)

	// Usecases
	usecases := &controller.Usecases{
		User:     usecase.NewUserUsecase(userRepo),
		Subject:  usecase.NewSubjectUsecase(subjectRepo, userRepo),
		StudyLog: usecase.NewStudyLogUsecase(studyLogRepo, userRepo, subjectRepo),
		Goal:     usecase.NewGoalUsecase(goalRepo, userRepo, subjectRepo),
		Stats:    usecase.NewStatsUsecase(studyLogRepo, goalRepo, subjectRepo),
	}

	// Router
	router := controller.NewRouter(usecases, cfg.CORSOrigins, logger)

	// Server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		logger.Info("shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("server shutdown error", "error", err)
		}
	}()

	logger.Info(fmt.Sprintf("server starting on :%s", cfg.Port))
	logger.Info("OpenAPI docs available at /docs")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
