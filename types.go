// Package llmapi provides a unified interface for LLM providers.
// Both Anthropic and NovelAI implementations can implement this interface,
// allowing code to swap providers with minimal changes.
package llmapi

// Message represents a single message in a conversation.
type Message struct {
	Role    string `json:"role"`    // "user", "assistant", "system"
	Content string `json:"content"` // The message text
}

// Usage tracks token consumption for a conversation.
type Usage struct {
	InputTokens  int
	OutputTokens int
}

// StreamCallback is called for each token during streaming.
// text contains the new token(s), done indicates if streaming is complete.
type StreamCallback func(text string, done bool)

// Settings configures generation parameters.
// Provider implementations map these to their native formats.
type Settings struct {
	Model         string
	MaxTokens     int
	Temperature   float64
	TopP          float64
	TopK          int
	StopSequences []string

	// Provider-specific extensions
	Extra map[string]any
}

// DefaultSettings provides reasonable defaults.
var DefaultSettings = Settings{
	MaxTokens:   2048,
	Temperature: 1.0,
}
