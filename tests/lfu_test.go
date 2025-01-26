package tests

import (
	"testing"

	"github.com/hugocarreira/easycache/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// LFUTestSuite defines the test structure
type LFUTestSuite struct {
	suite.Suite
	c *cache.Cache
}

// Setup before each test
func (suite *LFUTestSuite) SetupTest() {
	suite.c = cache.New(&cache.Config{
		EvictionPolicy: cache.LFU,
		MaxSize:        2,
	})
}

// Test LFU (Least Frequently Used)
func (suite *LFUTestSuite) TestLFUEviction() {
	suite.c.Set("A", "Item A")
	suite.c.Set("B", "Item B")

	suite.c.Get("A")
	suite.c.Get("A")

	suite.c.Set("C", "Item C")

	assert.False(suite.T(), suite.c.Has("B"))
	assert.True(suite.T(), suite.c.Has("A"))
	assert.True(suite.T(), suite.c.Has("C"))
}

// Run the test suite
func TestLFUTestSuite(t *testing.T) {
	suite.Run(t, new(LFUTestSuite))
}
