package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAd(t *testing.T) {
	client := getTestClient()

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(0))
	assert.False(t, response.Data.Published)
}

func TestChangeAdStatus(t *testing.T) {
	client := getTestClient()

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	response, err = client.changeAdStatus(0, response.Data.ID, true)

	assert.NoError(t, err)
	assert.True(t, response.Data.Published)

	response, err = client.changeAdStatus(0, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)

	response, err = client.changeAdStatus(0, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)
}

func TestUpdateAd(t *testing.T) {
	client := getTestClient()

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	response, err = client.updateAd(0, response.Data.ID, "привет", "мир")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "привет")
	assert.Equal(t, response.Data.Text, "мир")
}

func TestGetAds(t *testing.T) {
	client := getTestClient()

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(0, "best cat", "not for sale")
	assert.NoError(t, err)
	ads, err := client.getAds()
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}
func TestGetAd(t *testing.T) {
	client := getTestClient()
	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	publishedAd, err := client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)
	ad, err := client.getAd(response.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, ad.Data.ID, publishedAd.Data.ID)
	assert.Equal(t, ad.Data.Title, publishedAd.Data.Title)
	assert.Equal(t, ad.Data.Text, publishedAd.Data.Text)
	assert.Equal(t, ad.Data.AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ad.Data.Published)
}

func TestGetAdsByTitle(t *testing.T) {
	client := getTestClient()
	response1, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	response2, err := client.getAdByTitle("hello")
	assert.NoError(t, err)
	assert.Equal(t, response1, response2)
}

func TestAdsByParamsFilter(t *testing.T) {
	client := getTestClient()
	_, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	_, err = client.changeAdStatus(0, 0, true)
	assert.NoError(t, err)
	response2, err := client.createAd(0, "hello", "world1")
	assert.NoError(t, err)
	_, err = client.changeAdStatus(0, 1, true)
	assert.NoError(t, err)
	_, err = client.createAd(0, "hello", "world2")
	assert.NoError(t, err)
	_, err = client.changeAdStatus(0, 2, true)
	assert.NoError(t, err)

	ads, err := client.getAdsByParams(map[string]any{"author_id": 0})
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 3)
	assert.Equal(t, ads.Data[1].ID, response2.Data.ID)
	assert.Equal(t, ads.Data[1].Title, response2.Data.Title)
	assert.Equal(t, ads.Data[1].AuthorID, response2.Data.AuthorID)
}

func TestDeleteAd(t *testing.T) {
	client := getTestClient()
	_, _ = client.createUser("helli", "world")
	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	_, err = client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)
	ads, err := client.getAds()
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	response, err = client.deleteAd(0, response.Data.ID)
	assert.NoError(t, err)
	ads, err = client.getAds()
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 0)
}
