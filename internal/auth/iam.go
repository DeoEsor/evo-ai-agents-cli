package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

// IAMTokenResponse представляет ответ от IAM API
type IAMTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// IAMTokenRequest представляет запрос к IAM API
type IAMTokenRequest struct {
	KeyID  string `json:"keyId"`
	Secret string `json:"secret"`
}

// IAMAuthServiceInterface определяет интерфейс для IAM аутентификации
type IAMAuthServiceInterface interface {
	GetToken(ctx context.Context) (string, error)
	IsAuthenticated() bool
	ClearToken()
}

// IAMAuthService управляет аутентификацией через IAM
type IAMAuthService struct {
	keyID    string
	secret   string
	endpoint string
	client   *http.Client

	// Кэш токена
	token     string
	expiresAt time.Time
	mutex     sync.RWMutex
}

// NewIAMAuthService создает новый сервис IAM аутентификации
func NewIAMAuthService(keyID, secret, endpoint string) *IAMAuthService {
	return &IAMAuthService{
		keyID:    keyID,
		secret:   secret,
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetToken возвращает действующий токен доступа
func (s *IAMAuthService) GetToken(ctx context.Context) (string, error) {
	s.mutex.RLock()
	if s.token != "" && time.Now().Before(s.expiresAt) {
		token := s.token
		s.mutex.RUnlock()
		return token, nil
	}
	s.mutex.RUnlock()

	// Токен истек или отсутствует, получаем новый
	return s.refreshToken(ctx)
}

// refreshToken получает новый токен от IAM API
func (s *IAMAuthService) refreshToken(ctx context.Context) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Проверяем еще раз, возможно токен уже обновился
	if s.token != "" && time.Now().Before(s.expiresAt) {
		log.Debug("Токен еще действителен, используем кэшированный")
		return s.token, nil
	}

	log.Debug("Получение нового токена от IAM API", "endpoint", s.endpoint, "key_id", s.keyID)

	// Подготавливаем запрос
	reqBody := IAMTokenRequest{
		KeyID:  s.keyID,
		Secret: s.secret,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Создаем HTTP запрос
	url := s.endpoint + "/api/v1/auth/token"
	log.Debug("Создание IAM запроса", "url", url, "body_size", len(jsonData))

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error("Ошибка создания IAM запроса", "error", err)
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	log.Debug("Отправка IAM запроса", "headers", req.Header)

	// Выполняем запрос
	resp, err := s.client.Do(req)
	if err != nil {
		log.Error("Ошибка выполнения IAM запроса", "error", err, "url", url)
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	log.Debug("IAM ответ получен", "status", resp.StatusCode, "headers", resp.Header)

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Ошибка чтения IAM ответа", "error", err)
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	log.Debug("IAM ответ прочитан", "body_length", len(body), "body", string(body))

	if resp.StatusCode != http.StatusOK {
		log.Error("IAM API вернул ошибку", "status", resp.StatusCode, "body", string(body))
		return "", fmt.Errorf("IAM API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Парсим ответ
	var tokenResp IAMTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		log.Error("Ошибка парсинга IAM ответа", "error", err, "body", string(body))
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	s.token = tokenResp.AccessToken
	margin := tokenResp.ExpiresIn - 300
	if margin < 60 {
		margin = 60
	}
	s.expiresAt = time.Now().Add(time.Duration(margin) * time.Second)

	log.Debug("Токен успешно получен", "expires_at", s.expiresAt, "expires_in", tokenResp.ExpiresIn)

	return s.token, nil
}

// IsAuthenticated проверяет, есть ли действующий токен
func (s *IAMAuthService) IsAuthenticated() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.token != "" && time.Now().Before(s.expiresAt)
}

// ClearToken очищает сохраненный токен
func (s *IAMAuthService) ClearToken() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.token = ""
	s.expiresAt = time.Time{}
}
