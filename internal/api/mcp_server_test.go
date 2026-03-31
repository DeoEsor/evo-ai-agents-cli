package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMCPServerService_List(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем URL
		expectedPath := "/api/v1/test-project/mcpServers"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		// Проверяем query параметры
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("Expected limit=10, got %s", r.URL.Query().Get("limit"))
		}
		if r.URL.Query().Get("offset") != "0" {
			t.Errorf("Expected offset=0, got %s", r.URL.Query().Get("offset"))
		}

		// Отправляем ответ
		response := MCPServerListResponse{
			Data: []MCPServer{
				{
					ID:          "1",
					Name:        "test-server",
					Description: "Test MCP server",
					Status:      "active",
					CreatedAt:   CustomTime{Time: time.Now()},
					UpdatedAt:   CustomTime{Time: time.Now()},
					Options:     map[string]interface{}{"host": "localhost"},
				},
			},
			Total: 1,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Создаем клиент и сервис
	mockAuth := &MockIAMService{token: "test-token"}
	client := NewClient(server.URL, "test-project", mockAuth)
	service := NewMCPServerService(client)

	// Тестируем List
	result, err := service.List(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.Total != 1 {
		t.Errorf("Expected total 1, got %d", result.Total)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 server, got %d", len(result.Data))
	}

	if result.Data[0].Name != "test-server" {
		t.Errorf("Expected name 'test-server', got '%s'", result.Data[0].Name)
	}
}

func TestMCPServerService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/v1/test-project/mcpServers/test-id"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		server := MCPServer{
			ID:          "test-id",
			Name:        "test-server",
			Description: "Test MCP server",
			Status:      "active",
			CreatedAt:   CustomTime{Time: time.Now()},
			UpdatedAt:   CustomTime{Time: time.Now()},
			Options:     map[string]interface{}{"host": "localhost"},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(server)
	}))
	defer server.Close()

	mockAuth := &MockIAMService{token: "test-token"}
	client := NewClient(server.URL, "test-project", mockAuth)
	service := NewMCPServerService(client)

	result, err := service.Get(context.Background(), "test-id")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", result.ID)
	}

	if result.Name != "test-server" {
		t.Errorf("Expected name 'test-server', got '%s'", result.Name)
	}
}

func TestMCPServerService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/v1/test-project/mcpServers"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Проверяем тело запроса
		var req MCPServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if req.Name != "new-server" {
			t.Errorf("Expected name 'new-server', got '%s'", req.Name)
		}

		// Отправляем ответ
		createdServer := MCPServer{
			ID:          "new-id",
			Name:        req.Name,
			Description: req.Description,
			Status:      "active",
			CreatedAt:   CustomTime{Time: time.Now()},
			UpdatedAt:   CustomTime{Time: time.Now()},
			Options:     req.Options,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdServer)
	}))
	defer server.Close()

	mockAuth := &MockIAMService{token: "test-token"}
	client := NewClient(server.URL, "test-project", mockAuth)
	service := NewMCPServerService(client)

	req := &MCPServerCreateRequest{
		Name:        "new-server",
		Description: "New MCP server",
		Options:     map[string]interface{}{"host": "localhost"},
	}

	result, err := service.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.Name != "new-server" {
		t.Errorf("Expected name 'new-server', got '%s'", result.Name)
	}

	if result.Status != "active" {
		t.Errorf("Expected status 'active', got '%s'", result.Status)
	}
}

func TestMCPServerService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/v1/test-project/mcpServers/test-id"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		if r.Method != "PATCH" {
			t.Errorf("Expected method PATCH, got %s", r.Method)
		}

		// Проверяем тело запроса
		var req MCPServerUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if req.Name != "updated-server" {
			t.Errorf("Expected name 'updated-server', got '%s'", req.Name)
		}

		// Отправляем ответ
		updatedServer := MCPServer{
			ID:          "test-id",
			Name:        req.Name,
			Description: req.Description,
			Status:      "active",
			CreatedAt:   CustomTime{Time: time.Now()},
			UpdatedAt:   CustomTime{Time: time.Now()},
			Options:     req.Options,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedServer)
	}))
	defer server.Close()

	mockAuth := &MockIAMService{token: "test-token"}
	client := NewClient(server.URL, "test-project", mockAuth)
	service := NewMCPServerService(client)

	req := &MCPServerUpdateRequest{
		Name:        "updated-server",
		Description: "Updated MCP server",
		Options:     map[string]interface{}{"host": "updated-host"},
	}

	result, err := service.Update(context.Background(), "test-id", req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.Name != "updated-server" {
		t.Errorf("Expected name 'updated-server', got '%s'", result.Name)
	}
}

func TestMCPServerService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/v1/test-project/mcpServers/test-id"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		if r.Method != "DELETE" {
			t.Errorf("Expected method DELETE, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	mockAuth := &MockIAMService{token: "test-token"}
	client := NewClient(server.URL, "test-project", mockAuth)
	service := NewMCPServerService(client)

	err := service.Delete(context.Background(), "test-id")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestMCPServerService_Resume(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/v1/test-project/mcpServers/resume/test-id"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		if r.Method != "PATCH" {
			t.Errorf("Expected method PATCH, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	mockAuth := &MockIAMService{token: "test-token"}
	client := NewClient(server.URL, "test-project", mockAuth)
	service := NewMCPServerService(client)

	err := service.Resume(context.Background(), "test-id")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestMCPServerService_Suspend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/v1/test-project/mcpServers/suspend/test-id"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		if r.Method != "PATCH" {
			t.Errorf("Expected method PATCH, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	mockAuth := &MockIAMService{token: "test-token"}
	client := NewClient(server.URL, "test-project", mockAuth)
	service := NewMCPServerService(client)

	err := service.Suspend(context.Background(), "test-id")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestMCPServerService_GetHistory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/v1/test-project/mcpServers/test-id/history"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		response := MCPServerHistoryResponse{
			Data: []HistoryEntry{
				{
					ID:        "1",
					Action:    "create",
					Status:    "success",
					Message:   "Server created successfully",
					CreatedAt: time.Now(),
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	mockAuth := &MockIAMService{token: "test-token"}
	client := NewClient(server.URL, "test-project", mockAuth)
	service := NewMCPServerService(client)

	result, err := service.GetHistory(context.Background(), "test-id")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 history entry, got %d", len(result.Data))
	}

	if result.Data[0].Action != "create" {
		t.Errorf("Expected action 'create', got '%s'", result.Data[0].Action)
	}
}

func TestMCPServerService_GetTools(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/v1/test-project/mcpServers/test-id/tools"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		response := struct {
			Tools []Tool `json:"tools"`
		}{
			Tools: []Tool{
				{
					Name:        "test-tool",
					Description: "Test tool",
					InputSchema: map[string]interface{}{"type": "object"},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	mockAuth := &MockIAMService{token: "test-token"}
	client := NewClient(server.URL, "test-project", mockAuth)
	service := NewMCPServerService(client)

	result, err := service.GetTools(context.Background(), "test-id")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(result))
	}

	if result[0].Name != "test-tool" {
		t.Errorf("Expected tool name 'test-tool', got '%s'", result[0].Name)
	}
}

func TestMCPServerService_ExecuteTool(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/v1/test-project/mcpServers/test-id/execute"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}

		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Проверяем тело запроса
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if req["tool_name"] != "test-tool" {
			t.Errorf("Expected tool_name 'test-tool', got '%s'", req["tool_name"])
		}

		// Отправляем ответ
		response := map[string]interface{}{
			"result": "success",
			"data":   "test result",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	mockAuth := &MockIAMService{token: "test-token"}
	client := NewClient(server.URL, "test-project", mockAuth)
	service := NewMCPServerService(client)

	params := map[string]interface{}{
		"param1": "value1",
		"param2": "value2",
	}

	result, err := service.ExecuteTool(context.Background(), "test-id", "test-tool", params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Errorf("Expected result, got nil")
	}
}
