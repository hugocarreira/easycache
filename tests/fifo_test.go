package tests

import (
	"testing"

	"github.com/hugocarreira/easycache/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// FIFOTestSuite defines the test structure
type FIFOTestSuite struct {
	suite.Suite
	c *cache.Cache
}

// Setup before each test
func (suite *FIFOTestSuite) SetupTest() {
	suite.c = cache.New(&cache.Config{
		EvictionPolicy: cache.FIFO,
		MaxSize:        2,
	})
}

// Test FIFO (First-In, First-Out)
func (suite *FIFOTestSuite) TestFIFOEviction() {
	suite.c.Set("A", "Item A")
	suite.c.Set("B", "Item B")
	suite.c.Set("C", "Item C")

	assert.False(suite.T(), suite.c.Has("A"))
	assert.True(suite.T(), suite.c.Has("B"))
	assert.True(suite.T(), suite.c.Has("C"))
}

// Run the test suite
func TestFIFOTestSuite(t *testing.T) {
	suite.Run(t, new(FIFOTestSuite))
}
