package llmapi

// Sampling contains per-call sampling parameters.
// Zero values mean "use conversation defaults".
type Sampling struct {
	TopK        int     // 0 = use default, 1 = deterministic
	Temperature float64 // 0 = use default
	TopP        float64 // 0 = use default
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
	//   - err: Any error that occurred
	Send(text string, sampling Sampling) (reply, stopReason string, inputTokens, outputTokens int, err error)

	// SendStreaming sends a message with real-time token streaming via SSE.
	// The callback is invoked for each token received.
	// Sampling parameters override conversation defaults for this call only.
	SendStreaming(text string, sampling Sampling, callback StreamCallback) (reply, stopReason string, inputTokens, outputTokens int, err error)

	// SendUntilDone repeatedly calls Send until stopReason != "max_tokens".
	// Returns the complete accumulated output.
	SendUntilDone(text string, sampling Sampling) (reply, stopReason string, inputTokens, outputTokens int, err error)

	// SendStreamingUntilDone combines streaming with auto-continuation.
	SendStreamingUntilDone(text string, sampling Sampling, callback StreamCallback) (reply, stopReason string, inputTokens, outputTokens int, err error)

	// AddMessage manually adds a message to the conversation history.
	AddMessage(role, content string)

	// GetMessages returns the current conversation history.
	GetMessages() []Message

	// GetUsage returns cumulative token usage for this conversation.
	GetUsage() Usage

	// GetSystem returns the system prompt.
	GetSystem() string

	// Clear resets the conversation history but keeps the system prompt and settings.
	Clear()

	// SetModel changes the model for subsequent API calls.
	SetModel(model string)

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
	AddRichMessage(role string, content []ContentBlock)

	// GetRichMessages returns the conversation history with full content blocks.
	// This preserves images, tool use, thinking, etc that GetMessages(Fix) loses.
	GetRichMessages() []RichMessage

	// SetTools configures the available tools for this conversation.
	// Pass nil or empty slice to disable tools.
	SetTools(tools []ToolDefinition)

	// GetTools returns the currently configured tools.
	GetTools() []ToolDefinition
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
)
