package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"morc/storage"
	"morc/storage/redis"
)

type Store struct {
	db     *pgxpool.Pool
	redis  storage.RedisRepoI
	user   storage.UserRepoI
	barrel storage.BarrelRepository
}

func NewPostgres(psqlConnString string, redisAddr, redisPassword string, redisDB int) storage.StorageRepoI {
	config, err := pgxpool.ParseConfig(psqlConnString)
	if err != nil {
		log.Panicf("Unable to parse connection string.: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Panicf("Unable to connect to the database: %v", err)
	}

	return &Store{
		db:    pool,
		redis: redis.NewRedisClient(redisAddr, redisPassword, redisDB),
	}
}

func (s *Store) CloseDB() {
	s.db.Close()
}

func (s *Store) User() storage.UserRepoI {
	if s.user == nil {
		s.user = &userRepo{
			db: s.db,
		}
	}
	return s.user
}

func (s *Store) Barrel() storage.BarrelRepository {
	if s.barrel == nil {
		s.barrel = &barrelRepo{
			db: s.db,
		}
	}
	return s.barrel
}


func (p *Store) Redis() storage.RedisRepoI {
	return p.redis
}