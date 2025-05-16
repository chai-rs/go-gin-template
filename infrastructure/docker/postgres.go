package docker

import (
	"fmt"
	syslog "log"
	"os"
	"time"

	"github.com/chai-rs/simple-bookstore/pkg/migration"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	PostgresUser     = "postgres"
	PostgresPassword = "postgres"
	PostgresDB       = "postgres"
)

// ! used for testing only:
// RunPostgresContainer runs a postgres container
func RunPostgresContainer(secondToKill uint) (*dockertest.Pool, *dockertest.Resource, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, err
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "17",
		Env: []string{
			"listen_addresses = '*'",
			fmt.Sprintf("POSTGRES_USER=%s", PostgresUser),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", PostgresPassword),
			fmt.Sprintf("POSTGRES_DB=%s", PostgresDB),
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432/tcp": {
				{
					HostIP:   "0.0.0.0",
					HostPort: "0", // 0 = random port
				},
			},
		},
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
		hc.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})

	if err != nil {
		return nil, nil, err
	}

	if err := resource.Expire(secondToKill); err != nil {
		return nil, nil, err
	}

	pool.MaxWait = PoolMaxWait
	if err := pool.Retry(func() error { return err }); err != nil {
		return nil, nil, err
	}

	return pool, resource, nil
}

// ! used for testing only:
// OpenPostgresGormDB opens a gorm db connection to the postgres container
func OpenPostgresGormDB(pool *dockertest.Pool, resource *dockertest.Resource) (*gorm.DB, error) {
	dns := "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
	retries := 10
	host := resource.GetBoundIP("5432/tcp")
	port := resource.GetPort("5432/tcp")
	url := fmt.Sprintf(dns, host, port, PostgresUser, PostgresPassword, PostgresDB)

	gormLogger := logger.New(
		syslog.New(os.Stdout, "\r\n", syslog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: gormLogger,
	})

	for err != nil {
		if retries > 1 {
			retries--
			time.Sleep(1 * time.Second)
			db, err = gorm.Open(postgres.Open(url), &gorm.Config{})
			continue
		}

		if err := pool.Purge(resource); err != nil {
			log.Error().Err(err).Msg("🚨 failed to purge resource")
		}

		return nil, err
	}

	return db, nil
}

// ! used for testing only:
// MigratePostgreSQL migrates the postgres container
func MigratePostgreSQL(pool *dockertest.Pool, resource *dockertest.Resource, migrationDir string) error {
	input := migration.PostgresMigrationInput{
		Username: PostgresUser,
		Password: PostgresPassword,
		Host:     resource.GetBoundIP("5432/tcp"),
		Port:     resource.GetPort("5432/tcp"),
		DBName:   PostgresDB,
		File:     migrationDir,
	}

	retries := 10
	migrator := migration.NewPostgresMigration(&input)
	err := migrator.Up()
	for err != nil {
		if retries > 1 {
			retries--
			time.Sleep(1 * time.Second)
			err = migrator.Up()
			continue
		}

		if err := pool.Purge(resource); err != nil {
			log.Error().Err(err).Msg("🚨 failed to purge resource")
		}

		return err
	}

	return nil
}
