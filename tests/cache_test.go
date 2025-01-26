package tests

import (
	"testing"
	"time"

	"github.com/hugocarreira/easycache/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// CacheTestSuite defines the test structure
type CacheTestSuite struct {
	suite.Suite
	c *cache.Cache
}

// Setup before each test
func (suite *CacheTestSuite) SetupTest() {
	suite.c = cache.New(&cache.Config{
		EvictionPolicy: cache.Basic,
		MaxSize:        3,
		TTL:            60 * time.Second,
	})
}

// Test `Set()` e `Get()`
func (suite *CacheTestSuite) TestSetGet() {
	suite.c.Set("A", "Item A")
	suite.c.Set("B", "Item B")

	val, found := suite.c.Get("A")
	assert.True(suite.T(), found)
	assert.Equal(suite.T(), "Item A", val)

	val, found = suite.c.Get("B")
	assert.True(suite.T(), found)
	assert.Equal(suite.T(), "Item B", val)

	_, found = suite.c.Get("X")
	assert.False(suite.T(), found)
}

// Test `Delete()`
func (suite *CacheTestSuite) TestDelete() {
	suite.c.Set("A", "Item A")
	suite.c.Delete("A")

	_, found := suite.c.Get("A")
	assert.False(suite.T(), found)
}

// Test `Has()`
func (suite *CacheTestSuite) TestHas() {
	suite.c.Set("A", "Item A")

	assert.True(suite.T(), suite.c.Has("A"))
	assert.False(suite.T(), suite.c.Has("B"))
}

// Test `Len()`
func (suite *CacheTestSuite) TestLen() {
	suite.c.Set("A", "Item A")
	suite.c.Set("B", "Item B")

	assert.Equal(suite.T(), 2, suite.c.Len())
}

// Run the test suite
func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}
