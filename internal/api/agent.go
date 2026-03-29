package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// CustomTime представляет кастомный тип времени для парсинга
type CustomTime struct {
	time.Time
}

// UnmarshalJSON кастомный парсинг JSON для времени
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	// Убираем кавычки
	s := string(b[1 : len(b)-1])

	// Пробуем разные форматы времени
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05.000000Z",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			ct.Time = t
			return nil
		}
	}

	// Если ничего не сработало, возвращаем ошибку
	return fmt.Errorf("unable to parse time: %s", s)
}

// StatusReason представляет причину статуса
type StatusReason struct {
	ReasonType string                 `json:"reasonType,omitempty"`
	Key        string                 `json:"key,omitempty"`
	Message    string                 `json:"message,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// InstanceType представляет тип инстанса
type InstanceType struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	SKUCode      string `json:"skuCode,omitempty"`
	ResourceCode string `json:"resourceCode,omitempty"`
	IsActive     bool   `json:"isActive,omitempty"`
	MCPU         int    `json:"mCpu,omitempty"`
	MibRAM       int    `json:"mibRam,omitempty"`
	CreatedAt    string `json:"createdAt,omitempty"`
	UpdatedAt    string `json:"updatedAt,omitempty"`
	CreatedBy    string `json:"createdBy,omitempty"`
	UpdatedBy    string `json:"updatedBy,omitempty"`
}

// Agent представляет агента
type Agent struct {
	ID                    string                 `json:"id"`
	ProjectID             string                 `json:"projectId,omitempty"`
	Name                  string                 `json:"name"`
	Description           string                 `json:"description"`
	Status                string                 `json:"status"`
	StatusReason          StatusReason           `json:"statusReason,omitempty"`
	InstanceType          InstanceType           `json:"instanceType,omitempty"`
	MCPServers            []MCPServerReference   `json:"mcpServers,omitempty"`
	ImageSource           map[string]interface{} `json:"imageSource,omitempty"`
	AgentType             string                 `json:"agentType,omitempty"`
	Options               map[string]interface{} `json:"options,omitempty"`
	IntegrationOptions    map[string]interface{} `json:"integrationOptions,omitempty"`
	UsedInAgentSystems    []AgentSystemPreview   `json:"usedInAgentSystems,omitempty"`
	PublicURL             string                 `json:"publicUrl,omitempty"`
	ArizePhoenixPublicURL string                 `json:"arizePhoenixPublicUrl,omitempty"`
	CreatedAt             CustomTime             `json:"createdAt"`
	UpdatedAt             CustomTime             `json:"updatedAt"`
	CreatedBy             string                 `json:"createdBy,omitempty"`
	UpdatedBy             string                 `json:"updatedBy,omitempty"`
	// Обратная совместимость
	MCPs []string `json:"mcp_servers,omitempty"`
}

// MCPServerReference представляет ссылку на MCP сервер
type MCPServerReference struct {
	ID     string                 `json:"mcpServerId"`
	Name   string                 `json:"name"`
	Status string                 `json:"status,omitempty"`
	Source map[string]interface{} `json:"source,omitempty"`
	Tools  []MCPTool              `json:"tools,omitempty"`
}

// MCPTool представляет инструмент MCP сервера
type MCPTool struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Args        []MCPToolArg `json:"args,omitempty"`
}

// MCPToolArg представляет аргумент инструмента MCP
type MCPToolArg struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// AgentSystemPreview представляет превью системы агентов
type AgentSystemPreview struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AgentListResponse представляет ответ со списком агентов
type AgentListResponse struct {
	Data  []Agent `json:"data"`
	Total int     `json:"total"`
}

// AgentHistoryResponse представляет ответ с историей агента
type AgentHistoryResponse struct {
	Data []HistoryEntry `json:"data"`
}

// MarketplaceAgent представляет агента из маркетплейса
type MarketplaceAgent struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"`
	Type        string                 `json:"type"`
	Categories  []string               `json:"categories"`
	Tags        []string               `json:"tags"`
	Options     map[string]interface{} `json:"options"`
}

// MarketplaceAgentListResponse представляет ответ со списком агентов из маркетплейса
type MarketplaceAgentListResponse struct {
	Data       []MarketplaceAgent `json:"data"`
	Total      int                `json:"total"`
	Categories []string           `json:"categories"`
	Tags       []string           `json:"tags"`
}

// AgentService предоставляет методы для работы с агентами
type AgentService struct {
	client *Client
}

// NewAgentService создает новый сервис для работы с агентами
func NewAgentService(client *Client) *AgentService {
	return &AgentService{client: client}
}

// List возвращает список агентов
func (s *AgentService) List(ctx context.Context, limit, offset int) (*AgentListResponse, error) {
	query := map[string]string{
		"limit":  fmt.Sprintf("%d", limit),
		"offset": fmt.Sprintf("%d", offset),
	}

	var result AgentListResponse
	err := s.client.Get(ctx, fmt.Sprintf("/api/v1/%s/agents", s.client.projectID), query, &result)
	return &result, err
}

// AgentGetResponse представляет ответ с информацией об агенте
type AgentGetResponse struct {
	Agent Agent `json:"agent"`
}

// Get возвращает информацию о конкретном агенте
func (s *AgentService) Get(ctx context.Context, agentID string) (*Agent, error) {
	// Сначала парсим в map, чтобы понять структуру
	var rawResponse map[string]interface{}
	err := s.client.Get(ctx, fmt.Sprintf("/api/v1/%s/agents/%s", s.client.projectID, agentID), nil, &rawResponse)
	if err != nil {
		return nil, err
	}

	// Пробуем разные варианты структуры
	if agentData, ok := rawResponse["agent"]; ok {
		// Если есть поле "agent"
		agentBytes, err := json.Marshal(agentData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal agent data: %w", err)
		}

		var agent Agent
		if err := json.Unmarshal(agentBytes, &agent); err != nil {
			return nil, fmt.Errorf("failed to unmarshal agent: %w", err)
		}

		return &agent, nil
	}

	// Если нет поля "agent", возможно данные прямо в корне
	agentBytes, err := json.Marshal(rawResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var agent Agent
	if err := json.Unmarshal(agentBytes, &agent); err != nil {
		return nil, fmt.Errorf("failed to unmarshal agent from root: %w", err)
	}

	return &agent, nil
}

// AgentCreateRequest представляет запрос на создание агента
type AgentCreateRequest struct {
	Name               string                 `json:"name"`
	Description        string                 `json:"description,omitempty"`
	InstanceTypeID     string                 `json:"instance_type_id,omitempty"`
	ExportedPorts      []int                  `json:"exported_ports,omitempty"`
	ImageSource        map[string]interface{} `json:"image_source,omitempty"`
	Options            map[string]interface{} `json:"options,omitempty"`
	MCPServers         []string               `json:"mcp_servers,omitempty"`
	IntegrationOptions map[string]interface{} `json:"integration_options,omitempty"`
}

// Create создает нового агента
func (s *AgentService) Create(ctx context.Context, req *AgentCreateRequest) (*Agent, error) {
	var result Agent
	err := s.client.Post(ctx, fmt.Sprintf("/api/v1/%s/agents", s.client.projectID), req, &result)
	return &result, err
}

// AgentUpdateRequest представляет запрос на обновление агента
type AgentUpdateRequest struct {
	Name               string                 `json:"name,omitempty"`
	Description        string                 `json:"description,omitempty"`
	InstanceTypeID     string                 `json:"instance_type_id,omitempty"`
	ExportedPorts      []int                  `json:"exported_ports,omitempty"`
	ImageSource        map[string]interface{} `json:"image_source,omitempty"`
	Options            map[string]interface{} `json:"options,omitempty"`
	MCPServerID        string                 `json:"mcp_server_id,omitempty"`
	IntegrationOptions map[string]interface{} `json:"integration_options,omitempty"`
}

// Update обновляет существующего агента
func (s *AgentService) Update(ctx context.Context, agentID string, req *AgentUpdateRequest) (*Agent, error) {
	var result Agent
	err := s.client.Put(ctx, fmt.Sprintf("/api/v1/%s/agents/%s", s.client.projectID, agentID), req, &result)
	return &result, err
}

// Delete удаляет агента
func (s *AgentService) Delete(ctx context.Context, agentID string) error {
	return s.client.Delete(ctx, fmt.Sprintf("/api/v1/%s/agents/%s", s.client.projectID, agentID), nil)
}

// Resume возобновляет работу агента
func (s *AgentService) Resume(ctx context.Context, agentID string) error {
	return s.client.Post(ctx, fmt.Sprintf("/api/v1/%s/agents/resume/%s", s.client.projectID, agentID), nil, nil)
}

// Suspend приостанавливает работу агента
func (s *AgentService) Suspend(ctx context.Context, agentID string) error {
	return s.client.Post(ctx, fmt.Sprintf("/api/v1/%s/agents/suspend/%s", s.client.projectID, agentID), nil, nil)
}

// GetHistory возвращает историю операций агента
func (s *AgentService) GetHistory(ctx context.Context, agentID string) (*AgentHistoryResponse, error) {
	var result AgentHistoryResponse
	err := s.client.Get(ctx, fmt.Sprintf("/api/v1/%s/agents/%s/history", s.client.projectID, agentID), nil, &result)
	return &result, err
}

// SearchMarketplace ищет агентов в маркетплейсе
func (s *AgentService) SearchMarketplace(ctx context.Context, req *MarketplaceSearchRequest) (*MarketplaceAgentListResponse, error) {
	query := make(map[string]string)
	if req.Limit > 0 {
		query["limit"] = fmt.Sprintf("%d", req.Limit)
	}
	if req.Offset > 0 {
		query["offset"] = fmt.Sprintf("%d", req.Offset)
	}
	if req.Name != "" {
		query["name"] = req.Name
	}
	if len(req.Tags) > 0 {
		query["tags"] = joinStrings(req.Tags, ",")
	}
	if len(req.Categories) > 0 {
		query["categories"] = joinStrings(req.Categories, ",")
	}
	if len(req.Statuses) > 0 {
		query["statuses"] = joinStrings(req.Statuses, ",")
	}
	if len(req.Types) > 0 {
		query["types"] = joinStrings(req.Types, ",")
	}

	var result MarketplaceAgentListResponse
	err := s.client.Get(ctx, "/api/v1/marketplace/agents", query, &result)
	return &result, err
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

// GetMarketplaceAgent возвращает информацию об агенте из маркетплейса
func (s *AgentService) GetMarketplaceAgent(ctx context.Context, agentID string) (*MarketplaceAgent, error) {
	var result struct {
		PredefinedAgent MarketplaceAgent `json:"predefined_agent"`
	}
	err := s.client.Get(ctx, fmt.Sprintf("/api/v1/marketplace/agents/%s", agentID), nil, &result)
	return &result.PredefinedAgent, err
}

// MarketplaceSearchRequest представляет запрос поиска в маркетплейсе
type MarketplaceSearchRequest struct {
	Limit      int      `json:"limit,omitempty"`
	Offset     int      `json:"offset,omitempty"`
	Name       string   `json:"name,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	Categories []string `json:"categories,omitempty"`
	Statuses   []string `json:"statuses,omitempty"`
	Types      []string `json:"types,omitempty"`
}

// User представляет пользователя
type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Username  string     `json:"username"`
	CreatedAt CustomTime `json:"createdAt"`
	UpdatedAt CustomTime `json:"updatedAt"`
	IsActive  bool       `json:"isActive"`
	Roles     []string   `json:"roles,omitempty"`
}
