package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TestBookSuite runs the book suite.
func TestBookSuite(t *testing.T) {
	suite.Run(t, new(BookCreateSuite))
}

// BookCreateSuite runs the book create suite.
type BookCreateSuite struct{ BaseSuite }
