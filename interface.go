package llmapi

import "context"

// ReasoningEffort requests how much the model reasons before answering. Providers
// map it to their native control: anthropic → extended-thinking budget tokens,
// openai (vLLM) → chat_template_kwargs reasoning_effort / enable_thinking, novelai
// → think on/off. The zero value is ReasoningOff — callers opt IN to reasoning,
// matching the assumption that reasoning is off by default.
type ReasoningEffort int

const (
	ReasoningOff ReasoningEffort = iota
	ReasoningLow
	ReasoningMedium
	ReasoningHigh
	ReasoningMax
)

// String returns the lowercase effort level ("off", "low", "medium", "high",
// "max") — the wire value used by providers that take a reasoning_effort string.
func (e ReasoningEffort) String() string {
	switch e {
	case ReasoningLow:
		return "low"
	case ReasoningMedium:
		return "medium"
	case ReasoningHigh:
		return "high"
	case ReasoningMax:
		return "max"
	default:
		return "off"
	}
}

// Sampling contains per-call sampling parameters.
// Zero values mean "use conversation defaults", except ReasoningEffort whose zero
// value is ReasoningOff (reasoning disabled).
type Sampling struct {
	TopK            int             // 0 = use default, 1 = deterministic
	Temperature     float64         // 0 = use default
	TopP            float64         // 0 = use default
	ReasoningEffort ReasoningEffort // zero value = ReasoningOff (reasoning disabled)
}

// Conversation is the primary interface for LLM interactions.
// Both anthropic.Conversation and novelai.Conversation implement this.
type Conversation interface {
	// Send sends a user message and returns the assistant's reply.
	// If text is empty, continues from the last assistant message (for max_tokens continuation).
	// Sampling parameters override conversation defaults for this call only.
	//
	// Returns:
	//   - reply: The assistant's response text
	//   - stopReason: Normalized stop reason ("end_turn", "max_tokens", "stop_sequence")
	//   - inputTokens: Tokens used for this request's input
	//   - outputTokens: Tokens generated in this response
	//   - cacheCreationTokens: Tokens written to cache (Anthropic only, 0 for others)
	//   - cacheReadTokens: Tokens read from cache (Anthropic only, 0 for others)
	//   - err: Any error that occurred
	Send(text string, sampling Sampling) (reply, stopReason string, inputTokens, outputTokens, cacheCreationTokens, cacheReadTokens int, err error)

	// SendStreaming sends a message with real-time token streaming via SSE.
	// The callback is invoked for each token received.
	// Sampling parameters override conversation defaults for this call only.
	SendStreaming(text string, sampling Sampling, callback StreamCallback) (reply, stopReason string, inputTokens, outputTokens, cacheCreationTokens, cacheReadTokens int, err error)

	// SendUntilDone repeatedly calls Send until stopReason != "max_tokens".
	// Returns the complete accumulated output.
	SendUntilDone(text string, sampling Sampling) (reply, stopReason string, inputTokens, outputTokens, cacheCreationTokens, cacheReadTokens int, err error)

	// SendStreamingUntilDone combines streaming with auto-continuation.
	SendStreamingUntilDone(text string, sampling Sampling, callback StreamCallback) (reply, stopReason string, inputTokens, outputTokens, cacheCreationTokens, cacheReadTokens int, err error)

	// AddMessage manually adds a message to the conversation history.
	AddMessage(role Role, content string)

	// GetMessages returns the current conversation history.
	GetMessages() []Message

	// GetUsage returns cumulative token usage for this conversation.
	GetUsage() Usage

	// GetSystem returns the system prompt.
	GetSystem() string

	// Clear resets the conversation history but keeps the system prompt and settings.
	Clear()

	// SetContext sets the context for cancellation and timeouts.
	// The context applies to all subsequent API calls until changed.
	// Pass nil to clear the context (will use context.Background()).
	SetContext(ctx context.Context)

	// SetModel changes the model for subsequent API calls.
	SetModel(model string)

	// SetEndpoint overrides the API endpoint URL for this conversation.
	// Pass empty string to revert to the provider's default endpoint.
	// This is useful for testing, proxies, or alternative API-compatible services.
	SetEndpoint(endpoint string)

	// SendRich sends a message with rich content blocks and returns a full response.
	// This enables multimodal input (images, documents) and captures all response
	// types (text, thinking, tool use).
	//
	// If content is nil or empty, continues from the last message (for max_tokens continuation).
	// Sampling parameters override conversation defaults for this call only.
	SendRich(content []ContentBlock, sampling Sampling) (*RichResponse, error)

	// SendRichStreaming sends rich content with streaming.
	// The callback receives text fragments as they arrive.
	// Returns the complete RichResponse after streaming completes.
	SendRichStreaming(content []ContentBlock, sampling Sampling, callback StreamCallback) (*RichResponse, error)

	// AddRichMessage adds a message with multiple content blocks to the history.
	// Use this for adding tool results, images, or other structured content.
	AddRichMessage(role Role, content []ContentBlock)

	// GetRichMessages returns the conversation history with full content blocks.
	// This preserves images, tool use, thinking, etc that GetMessages(Fix) loses.
	GetRichMessages() []RichMessage

	// SetTools configures the available tools for this conversation.
	// Pass nil or empty slice to disable tools.
	SetTools(tools []ToolDefinition)

	// GetTools returns the currently configured tools.
	GetTools() []ToolDefinition

	// EnableSystemCaching enables caching for the system prompt.
	// Returns an error if the provider does not support caching.
	EnableSystemCaching() error

	// EnableConversationCaching enables automatic cache breakpoints on
	// conversation turns. Before each API call, the conversation prefix is
	// marked for caching so subsequent turns can be served from cache.
	// Returns an error if the provider does not support caching.
	EnableConversationCaching() error

	// DisableConversationCaching disables automatic conversation turn caching.
	// Returns an error if the provider does not support caching.
	DisableConversationCaching() error
}

// CapabilityProvider is optionally implemented by Conversation implementations
// to advertise their capabilities.
type CapabilityProvider interface {
	// GetCapabilities returns the provider's capabilities.
	GetCapabilities() Capabilities
}

// ConversationFactory creates new conversations.
// Each provider implements this.
type ConversationFactory interface {
	NewConversation(system string) Conversation
}

// Provider identifies an LLM provider.
type Provider string

const (
	ProviderAnthropic Provider = "anthropic"
	ProviderNovelAI   Provider = "novelai"
	ProviderOpenAI    Provider = "openai"
)
