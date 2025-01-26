package tests

import (
	"testing"

	"github.com/hugocarreira/easycache/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// LRUTestSuite defines the test structure
type LRUTestSuite struct {
	suite.Suite
	c *cache.Cache
}

// Setup before each test
func (suite *LRUTestSuite) SetupTest() {
	suite.c = cache.New(&cache.Config{
		EvictionPolicy: cache.LRU,
		MaxSize:        2,
	})
}

// Test LRU (Least Recently Used)
func (suite *LRUTestSuite) TestLRUEviction() {
	suite.c.Set("A", "Item A")
	suite.c.Set("B", "Item B")

	suite.c.Get("A")

	suite.c.Set("C", "Item C")

	assert.False(suite.T(), suite.c.Has("B"))
	assert.True(suite.T(), suite.c.Has("A"))
	assert.True(suite.T(), suite.c.Has("C"))
}

// Run the test suite
func TestLRUTestSuite(t *testing.T) {
	suite.Run(t, new(LRUTestSuite))
}
