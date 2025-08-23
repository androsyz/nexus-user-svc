package initialize

import (
	"context"

	"github.com/androsyz/nexus-user-svc/config"
	"github.com/androsyz/nexus-user-svc/repository"
	"github.com/androsyz/nexus-user-svc/usecase"

	"github.com/rs/zerolog"
)

type App struct {
	UcUser *usecase.UcUser
}

func Bootstrap(ctx context.Context, cfg *config.Config, zlog zerolog.Logger) (App, error) {
	app := App{}

	// setup database
	zlog.Info().Msg("Initialize Database")
	dbConn, err := config.NewDatabase(cfg.Database)
	if err != nil {
		zlog.Error().Err(err).Msg("Failed initialize database")
		return app, err
	}

	// setup redis
	zlog.Info().Msg("Initialize Redis")
	rdsConn, err := config.NewRedis(cfg.Redis)
	if err != nil {
		zlog.Error().Err(err).Msg("Failed initialize redis")
		return app, err
	}

	// setup repository
	zlog.Info().Msg("Initialize Repository")
	repoUser := repository.NewUserRepository(dbConn, rdsConn)

	// setup usecase
	zlog.Info().Msg("Initialize Usecase")
	ucUser := usecase.NewUserUsecase(cfg, repoUser, zlog)

	return App{
		UcUser: ucUser,
	}, nil
}
