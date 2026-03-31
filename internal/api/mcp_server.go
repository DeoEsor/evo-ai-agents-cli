package api

import (
	"context"
	"fmt"
)

// MCPServer представляет MCP сервер
type MCPServer struct {
	ID                 string                 `json:"id"`
	Name               string                 `json:"name"`
	Description        string                 `json:"description"`
	Status             string                 `json:"status"`
	StatusReason       StatusReason           `json:"statusReason,omitempty"`
	InstanceType       InstanceType           `json:"instanceType,omitempty"`
	ImageSource        map[string]interface{} `json:"imageSource,omitempty"`
	ExposedPorts       []int                  `json:"exposedPorts,omitempty"`
	EnvironmentOptions map[string]interface{} `json:"environmentOptions,omitempty"`
	Scaling            map[string]interface{} `json:"scaling,omitempty"`
	IntegrationOptions map[string]interface{} `json:"integrationOptions,omitempty"`
	Options            map[string]interface{} `json:"options,omitempty"`
	Tools              []Tool                 `json:"tools,omitempty"`
	PublicURL          string                 `json:"publicUrl,omitempty"`
	CreatedAt          CustomTime             `json:"createdAt"`
	UpdatedAt          CustomTime             `json:"updatedAt"`
	CreatedBy          string                 `json:"createdBy,omitempty"`
	UpdatedBy          string                 `json:"updatedBy,omitempty"`
}

// Tool представляет инструмент MCP сервера
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"input_schema"`
}

// MCPServerCreateRequest представляет запрос на создание MCP сервера
type MCPServerCreateRequest struct {
	Name               string                 `json:"name"`
	Description        string                 `json:"description,omitempty"`
	InstanceTypeID     string                 `json:"instance_type_id,omitempty"`
	ImageSource        map[string]interface{} `json:"image_source,omitempty"`
	ExposedPorts       []int                  `json:"exposed_ports,omitempty"`
	EnvironmentOptions map[string]interface{} `json:"environment_options,omitempty"`
	Scaling            map[string]interface{} `json:"scaling,omitempty"`
	IntegrationOptions map[string]interface{} `json:"integration_options,omitempty"`
	Options            map[string]interface{} `json:"options,omitempty"`
}

// MCPServerUpdateRequest представляет запрос на обновление MCP сервера
type MCPServerUpdateRequest struct {
	Name               string                 `json:"name,omitempty"`
	Description        string                 `json:"description,omitempty"`
	InstanceTypeID     string                 `json:"instance_type_id,omitempty"`
	ImageSource        map[string]interface{} `json:"image_source,omitempty"`
	ExposedPorts       []int                  `json:"exposed_ports,omitempty"`
	EnvironmentOptions map[string]interface{} `json:"environment_options,omitempty"`
	Scaling            map[string]interface{} `json:"scaling,omitempty"`
	IntegrationOptions map[string]interface{} `json:"integration_options,omitempty"`
	Options            map[string]interface{} `json:"options,omitempty"`
}

// MCPServerListResponse представляет ответ со списком MCP серверов
type MCPServerListResponse struct {
	Data  []MCPServer `json:"data"`
	Total int         `json:"total"`
}

// MCPServerHistoryResponse представляет ответ с историей MCP сервера
type MCPServerHistoryResponse struct {
	Data []HistoryEntry `json:"data"`
}

// MCPServerService предоставляет методы для работы с MCP серверами
type MCPServerService struct {
	client *Client
}

// NewMCPServerService создает новый сервис для работы с MCP серверами
func NewMCPServerService(client *Client) *MCPServerService {
	return &MCPServerService{client: client}
}

// List возвращает список MCP серверов
func (s *MCPServerService) List(ctx context.Context, limit, offset int) (*MCPServerListResponse, error) {
	query := map[string]string{
		"limit":  fmt.Sprintf("%d", limit),
		"offset": fmt.Sprintf("%d", offset),
	}

	var result MCPServerListResponse
	err := s.client.Get(ctx, fmt.Sprintf("/api/v1/%s/mcpServers", s.client.projectID), query, &result)
	return &result, err
}

// Get возвращает информацию о конкретном MCP сервере
func (s *MCPServerService) Get(ctx context.Context, serverID string) (*MCPServer, error) {
	var result MCPServer
	err := s.client.Get(ctx, fmt.Sprintf("/api/v1/%s/mcpServers/%s", s.client.projectID, serverID), nil, &result)
	return &result, err
}

// Create создает новый MCP сервер
func (s *MCPServerService) Create(ctx context.Context, req *MCPServerCreateRequest) (*MCPServer, error) {
	var result MCPServer
	err := s.client.Post(ctx, fmt.Sprintf("/api/v1/%s/mcpServers", s.client.projectID), req, &result)
	return &result, err
}

// Update обновляет существующий MCP сервер
func (s *MCPServerService) Update(ctx context.Context, serverID string, req *MCPServerUpdateRequest) (*MCPServer, error) {
	var result MCPServer
	err := s.client.Patch(ctx, fmt.Sprintf("/api/v1/%s/mcpServers/%s", s.client.projectID, serverID), req, &result)
	return &result, err
}

// Delete удаляет MCP сервер
func (s *MCPServerService) Delete(ctx context.Context, serverID string) error {
	return s.client.Delete(ctx, fmt.Sprintf("/api/v1/%s/mcpServers/%s", s.client.projectID, serverID), nil)
}

// Resume возобновляет работу MCP сервера
func (s *MCPServerService) Resume(ctx context.Context, serverID string) error {
	return s.client.Patch(ctx, fmt.Sprintf("/api/v1/%s/mcpServers/resume/%s", s.client.projectID, serverID), struct{}{}, nil)
}

// Suspend приостанавливает работу MCP сервера
func (s *MCPServerService) Suspend(ctx context.Context, serverID string) error {
	return s.client.Patch(ctx, fmt.Sprintf("/api/v1/%s/mcpServers/suspend/%s", s.client.projectID, serverID), struct{}{}, nil)
}

// GetHistory возвращает историю операций MCP сервера
func (s *MCPServerService) GetHistory(ctx context.Context, serverID string) (*MCPServerHistoryResponse, error) {
	var result MCPServerHistoryResponse
	err := s.client.Get(ctx, fmt.Sprintf("/api/v1/%s/mcpServers/%s/history", s.client.projectID, serverID), nil, &result)
	return &result, err
}

// GetTools возвращает список инструментов MCP сервера
func (s *MCPServerService) GetTools(ctx context.Context, serverID string) ([]Tool, error) {
	var result struct {
		Tools []Tool `json:"tools"`
	}
	err := s.client.Get(ctx, fmt.Sprintf("/api/v1/%s/mcpServers/%s/tools", s.client.projectID, serverID), nil, &result)
	return result.Tools, err
}

// ExecuteTool выполняет инструмент MCP сервера
func (s *MCPServerService) ExecuteTool(ctx context.Context, serverID, toolName string, params map[string]interface{}) (interface{}, error) {
	req := map[string]interface{}{
		"tool_name":  toolName,
		"parameters": params,
	}

	var result interface{}
	err := s.client.Post(ctx, fmt.Sprintf("/api/v1/%s/mcpServers/%s/execute", s.client.projectID, serverID), req, &result)
	return result, err
}
