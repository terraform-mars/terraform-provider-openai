package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the OpenAI Admin API client.
type Client struct {
	adminKey   string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new OpenAI Admin API client.
func NewClient(adminKey, baseURL string) *Client {
	return &Client{
		adminKey: adminKey,
		baseURL:  baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Error represents an API error response.
type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

type errorResponse struct {
	Error Error `json:"error"`
}

func (c *Client) doRequest(ctx context.Context, method, path string, body any) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.adminKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp errorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error.Message != "" {
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Error.Message)
		}
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Project represents an OpenAI organization project.
type Project struct {
	ID         string `json:"id"`
	Object     string `json:"object"`
	Name       string `json:"name"`
	CreatedAt  int64  `json:"created_at"`
	ArchivedAt *int64 `json:"archived_at,omitempty"`
	Status     string `json:"status"`
}

// ListProjectsResponse is the response from listing projects.
type ListProjectsResponse struct {
	Object  string    `json:"object"`
	Data    []Project `json:"data"`
	FirstID string    `json:"first_id"`
	LastID  string    `json:"last_id"`
	HasMore bool      `json:"has_more"`
}

// CreateProjectRequest is the request body for creating a project.
type CreateProjectRequest struct {
	Name string `json:"name"`
}

// UpdateProjectRequest is the request body for updating a project.
type UpdateProjectRequest struct {
	Name string `json:"name"`
}

// ListProjects lists all projects in the organization.
func (c *Client) ListProjects(ctx context.Context) (*ListProjectsResponse, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/organization/projects", nil)
	if err != nil {
		return nil, err
	}

	var result ListProjectsResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetProject retrieves a project by ID.
func (c *Client) GetProject(ctx context.Context, projectID string) (*Project, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/organization/projects/"+projectID, nil)
	if err != nil {
		return nil, err
	}

	var result Project
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateProject creates a new project.
func (c *Client) CreateProject(ctx context.Context, req *CreateProjectRequest) (*Project, error) {
	resp, err := c.doRequest(ctx, http.MethodPost, "/organization/projects", req)
	if err != nil {
		return nil, err
	}

	var result Project
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateProject updates a project's name.
func (c *Client) UpdateProject(ctx context.Context, projectID string, req *UpdateProjectRequest) (*Project, error) {
	resp, err := c.doRequest(ctx, http.MethodPost, "/organization/projects/"+projectID, req)
	if err != nil {
		return nil, err
	}

	var result Project
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ArchiveProject archives a project. Projects cannot be deleted, only archived.
func (c *Client) ArchiveProject(ctx context.Context, projectID string) (*Project, error) {
	resp, err := c.doRequest(ctx, http.MethodPost, "/organization/projects/"+projectID+"/archive", nil)
	if err != nil {
		return nil, err
	}

	var result Project
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ProjectAPIKey represents an API key for a project.
type ProjectAPIKey struct {
	Object       string            `json:"object"`
	RedactedValue string           `json:"redacted_value"`
	Name         string            `json:"name"`
	CreatedAt    int64             `json:"created_at"`
	ID           string            `json:"id"`
	Owner        ProjectAPIKeyOwner `json:"owner"`
}

// ProjectAPIKeyOwner represents the owner of an API key.
type ProjectAPIKeyOwner struct {
	Type           string          `json:"type"` // "user" or "service_account"
	User           *User           `json:"user,omitempty"`
	ServiceAccount *ServiceAccount `json:"service_account,omitempty"`
}

// User represents a user in the organization.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ListProjectAPIKeysResponse is the response from listing API keys.
type ListProjectAPIKeysResponse struct {
	Object  string          `json:"object"`
	Data    []ProjectAPIKey `json:"data"`
	FirstID string          `json:"first_id"`
	LastID  string          `json:"last_id"`
	HasMore bool            `json:"has_more"`
}

// ListProjectAPIKeys lists all API keys in a project.
func (c *Client) ListProjectAPIKeys(ctx context.Context, projectID string) (*ListProjectAPIKeysResponse, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/organization/projects/"+projectID+"/api_keys", nil)
	if err != nil {
		return nil, err
	}

	var result ListProjectAPIKeysResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetProjectAPIKey retrieves an API key by ID.
func (c *Client) GetProjectAPIKey(ctx context.Context, projectID, keyID string) (*ProjectAPIKey, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/organization/projects/"+projectID+"/api_keys/"+keyID, nil)
	if err != nil {
		return nil, err
	}

	var result ProjectAPIKey
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteProjectAPIKey deletes an API key.
func (c *Client) DeleteProjectAPIKey(ctx context.Context, projectID, keyID string) error {
	_, err := c.doRequest(ctx, http.MethodDelete, "/organization/projects/"+projectID+"/api_keys/"+keyID, nil)
	return err
}

// ServiceAccount represents a service account.
type ServiceAccount struct {
	Object    string `json:"object"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt int64  `json:"created_at"`
}

// ServiceAccountWithKey is the response when creating a service account (includes the API key).
type ServiceAccountWithKey struct {
	Object    string             `json:"object"`
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Role      string             `json:"role"`
	CreatedAt int64              `json:"created_at"`
	APIKey    ServiceAccountKey  `json:"api_key"`
}

// ServiceAccountKey is the API key returned when creating a service account.
type ServiceAccountKey struct {
	Object    string `json:"object"`
	Value     string `json:"value"` // Full key value, only available on creation
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	ID        string `json:"id"`
}

// ListServiceAccountsResponse is the response from listing service accounts.
type ListServiceAccountsResponse struct {
	Object  string           `json:"object"`
	Data    []ServiceAccount `json:"data"`
	FirstID string           `json:"first_id"`
	LastID  string           `json:"last_id"`
	HasMore bool             `json:"has_more"`
}

// CreateServiceAccountRequest is the request body for creating a service account.
type CreateServiceAccountRequest struct {
	Name string `json:"name"`
}

// ListServiceAccounts lists all service accounts in a project.
func (c *Client) ListServiceAccounts(ctx context.Context, projectID string) (*ListServiceAccountsResponse, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/organization/projects/"+projectID+"/service_accounts", nil)
	if err != nil {
		return nil, err
	}

	var result ListServiceAccountsResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetServiceAccount retrieves a service account by ID.
func (c *Client) GetServiceAccount(ctx context.Context, projectID, serviceAccountID string) (*ServiceAccount, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/organization/projects/"+projectID+"/service_accounts/"+serviceAccountID, nil)
	if err != nil {
		return nil, err
	}

	var result ServiceAccount
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateServiceAccount creates a new service account and returns it with its API key.
func (c *Client) CreateServiceAccount(ctx context.Context, projectID string, req *CreateServiceAccountRequest) (*ServiceAccountWithKey, error) {
	resp, err := c.doRequest(ctx, http.MethodPost, "/organization/projects/"+projectID+"/service_accounts", req)
	if err != nil {
		return nil, err
	}

	var result ServiceAccountWithKey
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteServiceAccount deletes a service account.
func (c *Client) DeleteServiceAccount(ctx context.Context, projectID, serviceAccountID string) error {
	_, err := c.doRequest(ctx, http.MethodDelete, "/organization/projects/"+projectID+"/service_accounts/"+serviceAccountID, nil)
	return err
}
