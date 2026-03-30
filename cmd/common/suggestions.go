package common

// APISuggestions returns standard API error suggestions.
func APISuggestions() []string {
	return []string{
		"Проверьте переменные окружения: IAM_KEY_ID, IAM_SECRET_KEY, IAM_ENDPOINT",
		"Убедитесь что вы авторизованы: ai-agents-cli auth login",
		"Проверьте доступность API: curl -I $IAM_ENDPOINT",
		"📚 Документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution",
	}
}
