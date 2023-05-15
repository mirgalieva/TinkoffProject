package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type adData struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	AuthorID  int64  `json:"author_id"`
	Published bool   `json:"published"`
}

type adResponse struct {
	Data adData `json:"data"`
}

type adsResponse struct {
	Data []adData `json:"data"`
}

type paramsAdRequest struct {
	Params map[string]any `json:"params"`
}

func (tc *testClient) createAd(userID int64, title string, text string) (adResponse, error) {
	body := map[string]any{
		"user_id": userID,
		"title":   title,
		"text":    text,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, tc.baseURL+"/api/v1/ads", bytes.NewReader(data))
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

func (tc *testClient) changeAdStatus(userID int64, adID int64, published bool) (adResponse, error) {
	body := map[string]any{
		"user_id":   userID,
		"published": published,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d/status", adID), bytes.NewReader(data))
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

func (tc *testClient) updateAd(userID int64, adID int64, title string, text string) (adResponse, error) {
	body := map[string]any{
		"user_id": userID,
		"title":   title,
		"text":    text,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d", adID), bytes.NewReader(data))
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

func (tc *testClient) getAds() (adsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, tc.baseURL+"/api/v1/ads", nil)
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	var response adsResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adsResponse{}, err
	}

	return response, nil
}

func (tc *testClient) getAd(id int64) (adResponse, error) {
	body := map[string]any{
		"user_id": id,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d", id), bytes.NewReader(data))
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}
	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}
	return response, nil
}

func (tc *testClient) getAdByTitle(title string) (adResponse, error) {
	body := map[string]any{
		"title": title,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(tc.baseURL+"/api/v1/ads/title/%s", title), bytes.NewReader(data))
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}
	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}
	return response, nil
}

func (tc *testClient) getAdsByParams(params map[string]interface{}) (adsResponse, error) {
	data, err := json.Marshal(paramsAdRequest{params})
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodGet, tc.baseURL+"/api/v1/ads/params", bytes.NewReader(data))
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	var response adsResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adsResponse{}, err
	}
	return response, nil
}
