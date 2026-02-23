package dbrepo

import (
	"database/sql"

	"github.com/azwwz/bookingHotelTBMWAWG/internal/config"
	"github.com/azwwz/bookingHotelTBMWAWG/internal/repository"
)

type postgresDBRepo struct {
	DB  *sql.DB
	App *config.AppConfig
}

func NewPostgresDBRepo(d *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		DB:  d,
		App: a,
	}
}
