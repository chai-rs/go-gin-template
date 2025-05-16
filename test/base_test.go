package test

import (
	"github.com/chai-rs/simple-bookstore/config"
	"github.com/chai-rs/simple-bookstore/infrastructure/docker"
	"github.com/ory/dockertest/v3"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// BaseSuite provides a reusable test suite setup for integration tests requiring a PostgreSQL database.
// It uses dockertest to spin up a PostgreSQL container, run migrations, and provide a *gorm.DB instance.
type BaseSuite struct {
	suite.Suite                      // Embeds testify's suite for test lifecycle management
	db          *gorm.DB             // Database connection used in tests
	pool        *dockertest.Pool     // Docker pool for managing containers
	resource    *dockertest.Resource // Reference to the running PostgreSQL container
}

// SetupSuite is called once before any tests in the suite are run.
// It starts a PostgreSQL Docker container, establishes a GORM DB connection,
// and runs database migrations to prepare the test environment.
func (s *BaseSuite) SetupSuite() {
	var err error
	s.pool, s.resource, err = docker.RunPostgresContainer(20) // Start PostgreSQL container
	if err != nil {
		log.Error().Err(err).Msg("🚨 failed to run postgres container")
		s.T().FailNow()
	}

	s.db, err = docker.OpenPostgresGormDB(s.pool, s.resource) // Open GORM DB connection
	if err != nil {
		log.Error().Err(err).Msg("🚨 failed to open postgres gorm db")
		s.T().FailNow()
	}

	if err := docker.MigratePostgreSQL(s.pool, s.resource, "../db/migrations"); err != nil {
		log.Error().Err(err).Msg("🚨 failed to migrate postgres")
		s.T().FailNow()
	}

	setupConfig()
	// log.Info().Msg("🚀 Before Suite executed")
}

// TearDownSuite is called once after all tests in the suite have run.
// It cleans up the Docker container and any resources used by the test suite.
func (s *BaseSuite) TearDownSuite() {
	if err := s.pool.Purge(s.resource); err != nil {
		log.Error().Err(err).Msg("🚨 failed to purge resource")
		s.T().FailNow()
	}

	// log.Info().Msg("👋 After Suite executed")
}

func setupConfig() {
	config.ACCESS_SECRET = "secret"
	config.REFRESH_SECRET = "secret"
}
