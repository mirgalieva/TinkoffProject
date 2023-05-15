package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/adapters/userrepo"
	"homework9/internal/app"
	"homework9/internal/ports/httpgin"
	"io"
	"net/http"
	"net/http/httptest"
)

var (
	ErrBadRequest = fmt.Errorf("bad request")
	ErrForbidden  = fmt.Errorf("forbidden")
)

type testClient struct {
	client  *http.Client
	baseURL string
}

func getTestClient() *testClient {
	server := httpgin.NewHTTPServer(":18080", app.NewApp(adrepo.New(), userrepo.New()))
	testServer := httptest.NewServer(server.Handler())
	client := &testClient{
		client:  testServer.Client(),
		baseURL: testServer.URL,
	}
	//client.createUser("Admin", "adminEmail")   //nolint:errcheck
	//client.createUser("Admin2", "adminEmail2") //nolint:errcheck
	return client
}

func (tc *testClient) getResponse(req *http.Request, out any) error {
	resp, err := tc.client.Do(req)
	if err != nil {
		return fmt.Errorf("unexpected error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return ErrBadRequest
		}
		if resp.StatusCode == http.StatusForbidden {
			return ErrForbidden
		}
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read response: %w", err)
	}

	err = json.Unmarshal(respBody, out)
	if err != nil {
		return fmt.Errorf("unable to unmarshal: %w", err)
	}

	return nil
}

func (tc *testClient) deleteAd(userID int64, adID int64) (adResponse, error) {
	body := map[string]any{
		"user_id": userID,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d", adID), bytes.NewReader(data))
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	//req.Header.Add("Content-Type", "application/json")

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}
