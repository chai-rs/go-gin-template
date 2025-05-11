package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestBookSuite(t *testing.T) {
	suite.Run(t, new(BookCreateSuite))
}

type BookCreateSuite struct{ BaseSuite }
