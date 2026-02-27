package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/angolovin/yougile-cli/pkg/client"
)

// Login obtains an API key using email and password.
// It calls GetCompanies to resolve company ID (uses first company if multiple), then creates an auth key.
// baseURL is the YouGile API base (e.g. https://ru.yougile.com).
func Login(ctx context.Context, baseURL, email, password string) (apiKey string, err error) {
	baseURL = strings.TrimRight(baseURL, "/")
	api, err := client.NewClientWithResponses(baseURL)
	if err != nil {
		return "", fmt.Errorf("create client: %w", err)
	}

	companiesResp, err := api.GetCompaniesWithResponse(ctx, nil, client.GetCompaniesJSONRequestBody{
		Login:    email,
		Password: password,
	})
	if err != nil {
		return "", fmt.Errorf("get companies: %w", err)
	}

	if companiesResp.HTTPResponse.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get companies: HTTP %s", companiesResp.HTTPResponse.Status)
	}
	if companiesResp.JSON200 == nil || len(companiesResp.JSON200.Content) == 0 {
		return "", fmt.Errorf("no companies found for this account")
	}

	companyID := companiesResp.JSON200.Content[0].Id

	createResp, err := api.AuthKeyControllerCreateWithResponse(ctx, client.AuthKeyControllerCreateJSONRequestBody{
		Login:     email,
		Password:  password,
		CompanyId: companyID,
	})
	if err != nil {
		return "", fmt.Errorf("create key: %w", err)
	}

	if createResp.HTTPResponse.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("create key: HTTP %s", createResp.HTTPResponse.Status)
	}
	if createResp.JSON201 == nil || createResp.JSON201.Key == "" {
		return "", fmt.Errorf("create key: empty key in response")
	}

	return createResp.JSON201.Key, nil
}
