package helpers

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractAndValidateContextUserId_Success(t *testing.T) {
	expectedID := uuid.New()
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("UserId", expectedID)

	id, err := ExtractAndValidateContextUserId(c)
	require.NoError(t, err)
	assert.Equal(t, expectedID, id)
}

func TestExtractAndValidateContextUserId_NotFound(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	_, err := ExtractAndValidateContextUserId(c)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "UserId not found")
}

func TestExtractAndValidateContextUserId_InvalidType(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("UserId", "not-a-uuid")

	_, err := ExtractAndValidateContextUserId(c)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid user ID type")
}
