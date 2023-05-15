package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Nickname, "hello")
	assert.Equal(t, response.Data.Email, "world")
}

func TestGetUser(t *testing.T) {
	client := getTestClient()
	response, err := client.createUser("hello", "world")
	assert.NoError(t, err)
	response2, err := client.getUser(response.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, response2.Data, response.Data)
}

func TestDeleteUser(t *testing.T) {
	client := getTestClient()
	response, err := client.createUser("hello", "world")
	assert.NoError(t, err)
	_, err = client.getUser(response.Data.ID)
	assert.NoError(t, err)
	_, err = client.deleteUser(response.Data.ID)
	assert.NoError(t, err)
	_, err = client.getUser(response.Data.ID)
	assert.Error(t, err)
}
