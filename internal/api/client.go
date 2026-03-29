package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/auth"
)

// AuthenticationError представляет ошибку аутентификации
type AuthenticationError struct {
	StatusCode int
	Message    string
	Details    string
}

// Error возвращает строковое представление ошибки
func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("%s (статус %d): %s", e.Message, e.StatusCode, e.Details)
}

// Client представляет HTTP клиент для работы с AI Agents API
type Client struct {
	baseURL    string
	httpClient *http.Client
	projectID  string
	auth       auth.IAMAuthServiceInterface
}

// NewClient создает новый экземпляр API клиента с IAM аутентификацией
func NewClient(baseURL, projectID string, authService auth.IAMAuthServiceInterface) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		projectID: projectID,
		auth:      authService,
	}
}

// RequestOptions содержит опции для HTTP запроса
type RequestOptions struct {
	Method  string
	Path    string
	Body    interface{}
	Headers map[string]string
	Query   map[string]string
}

// doRequest выполняет HTTP запрос к API
func (c *Client) doRequest(ctx context.Context, opts RequestOptions) (*http.Response, error) {
	url := c.baseURL + opts.Path

	// Добавляем query параметры
	if len(opts.Query) > 0 {
		url += "?"
		first := true
		for key, value := range opts.Query {
			if !first {
				url += "&"
			}
			url += fmt.Sprintf("%s=%s", key, value)
			first = false
		}
	}

	var body io.Reader
	if opts.Body != nil {
		jsonData, err := json.Marshal(opts.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, opts.Method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")

	// Получаем токен для аутентификации
	if c.auth != nil {
		log.Debug("Getting auth token for API request")
		token, err := c.auth.GetToken(ctx)
		if err != nil {
			log.Error("Failed to get auth token", "error", err)
			return nil, fmt.Errorf("failed to get auth token: %w", err)
		}
		log.Debug("Auth token obtained successfully", "token_length", len(token))
		req.Header.Set("Authorization", "Bearer "+token)
	}

	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	log.Debug("Making API request", "method", opts.Method, "url", url, "headers", req.Header)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error("Failed to execute API request", "error", err, "url", url)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	log.Debug("API request completed", "status", resp.StatusCode, "url", url)
	return resp, nil
}

// parseResponse парсит ответ API в указанную структуру
func (c *Client) parseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read response body", "error", err, "status", resp.StatusCode)
		return fmt.Errorf("failed to read response body: %w", err)
	}

	log.Debug("API response received", "status", resp.StatusCode, "body_length", len(body))

	if resp.StatusCode >= 400 {
		log.Error("API error response", "status", resp.StatusCode, "body", string(body))

		// Проверяем на ошибки аутентификации
		if resp.StatusCode == 401 || resp.StatusCode == 403 {
			return &AuthenticationError{
				StatusCode: resp.StatusCode,
				Message:    "Ошибка аутентификации",
				Details:    "Проверьте настройки аутентификации и попробуйте снова",
			}
		}

		var errorResp struct {
			Error struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}

		if err := json.Unmarshal(body, &errorResp); err != nil {
			log.Error("Failed to parse error response", "error", err, "body", string(body))
			return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}

		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, errorResp.Error.Message)
	}

	if target != nil {
		// Отладочный вывод для агентов
		if len(body) > 0 && len(body) < 1000 {
			log.Debug("Raw response body", "body", string(body))
		}

		if err := json.Unmarshal(body, target); err != nil {
			log.Error("Failed to parse response JSON", "error", err, "body", string(body))
			return fmt.Errorf("failed to parse response: %w", err)
		}
		log.Debug("Response parsed successfully")
	}

	return nil
}

// Get выполняет GET запрос
func (c *Client) Get(ctx context.Context, path string, query map[string]string, result interface{}) error {
	resp, err := c.doRequest(ctx, RequestOptions{
		Method: "GET",
		Path:   path,
		Query:  query,
	})
	if err != nil {
		log.Error("Failed to execute GET request", "error", err, "url", fmt.Sprintf("%s%s", c.baseURL, path))
		return err
	}

	return c.parseResponse(resp, result)
}

// Post выполняет POST запрос
func (c *Client) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
	resp, err := c.doRequest(ctx, RequestOptions{
		Method: "POST",
		Path:   path,
		Body:   body,
	})
	if err != nil {
		log.Error("Failed to execute POST request", "error", err, "url", fmt.Sprintf("%s%s", c.baseURL, path))
		return err
	}

	return c.parseResponse(resp, result)
}

// Put выполняет PUT запрос
func (c *Client) Put(ctx context.Context, path string, body interface{}, result interface{}) error {
	resp, err := c.doRequest(ctx, RequestOptions{
		Method: "PUT",
		Path:   path,
		Body:   body,
	})
	if err != nil {
		log.Error("Failed to execute PUT request", "error", err, "url", fmt.Sprintf("%s%s", c.baseURL, path))
		return err
	}

	return c.parseResponse(resp, result)
}

// Patch выполняет PATCH запрос
func (c *Client) Patch(ctx context.Context, path string, body interface{}, result interface{}) error {
	resp, err := c.doRequest(ctx, RequestOptions{
		Method: "PATCH",
		Path:   path,
		Body:   body,
	})
	if err != nil {
		log.Error("Failed to execute PATCH request", "error", err, "url", fmt.Sprintf("%s%s", c.baseURL, path))
		return err
	}

	return c.parseResponse(resp, result)
}

// Delete выполняет DELETE запрос
func (c *Client) Delete(ctx context.Context, path string, result interface{}) error {
	resp, err := c.doRequest(ctx, RequestOptions{
		Method: "DELETE",
		Path:   path,
	})
	if err != nil {
		log.Error("Failed to execute DELETE request", "error", err, "url", fmt.Sprintf("%s%s", c.baseURL, path))
		return err
	}

	return c.parseResponse(resp, result)
}
