package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/spf13/viper"
)

type LoginRequest struct {
	Otp string `json:"otp"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CurrentUserResponse struct {
	Handle string `json:"handle"`
}

func FetchAccessToken() (*LoginResponse, error) {
	apiURL := APIBaseURL()
	client := &http.Client{}
	r, err := http.NewRequest("POST", apiURL+"/v1/auth/refresh", bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}
	r.Header.Add("X-Refresh-Token", viper.GetString("refresh_token"))
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("անվավեր refresh token")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var creds LoginResponse
	err = json.Unmarshal(body, &creds)
	return &creds, err
}

func LoginWithCode(code string) (*LoginResponse, error) {
	apiURL := APIBaseURL()
	req, err := json.Marshal(LoginRequest{Otp: code})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(apiURL+"/v1/auth/otp/login", "application/json", bytes.NewReader(req))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return nil, errors.New("սխալ login code. refresh արա browser-դ ու նորից փորձիր")
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var creds LoginResponse
	err = json.Unmarshal(body, &creds)
	if err != nil {
		return nil, err
	}

	return &creds, nil
}

func FetchCurrentUser() (*CurrentUserResponse, error) {
	body, err := fetchWithAuth("GET", "/v1/users/me")
	if err != nil {
		return nil, err
	}

	var user CurrentUserResponse
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func fetchWithAuth(method string, url string) ([]byte, error) {
	body, code, err := fetchWithAuthAndPayload(method, url, []byte{})
	if err != nil {
		return nil, err
	}
	if code == 402 {
		return nil, fmt.Errorf("էս դասի test-երը run և submit անելու համար պետք ա ունենաս ակտիվ Sovorem.am membership\nhttps://sovorem.am/pricing")
	}
	if code != 200 {
		return nil, fmt.Errorf("failed to %s to %s\nResponse: %d %s", method, url, code, string(body))
	}
	return body, err
}

func fetchWithAuthAndPayload(method string, url string, payload []byte) ([]byte, int, error) {
	apiURL := APIBaseURL()
	client := &http.Client{}
	r, err := http.NewRequest(method, apiURL+url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, 0, err
	}
	r.Header.Add("Authorization", "Bearer "+viper.GetString("access_token"))

	resp, err := client.Do(r)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}
