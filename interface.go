package llmapi

// Conversation is the primary interface for LLM interactions.
// Both anthropic.Conversation and novelai.Conversation implement this.
type Conversation interface {
	// Send sends a user message and returns the assistant's reply.
	// If text is empty, continues from the last assistant message (for max_tokens continuation).
	//
	// Returns:
	//   - reply: The assistant's response text
	//   - stopReason: Normalized stop reason ("end_turn", "max_tokens", "stop_sequence")
	//   - inputTokens: Tokens used for this request's input
	//   - outputTokens: Tokens generated in this response
	//   - err: Any error that occurred
	Send(text string) (reply, stopReason string, inputTokens, outputTokens int, err error)

	// SendStreaming sends a message with real-time token streaming via SSE.
	// The callback is invoked for each token received.
	SendStreaming(text string, callback StreamCallback) (reply, stopReason string, inputTokens, outputTokens int, err error)

	// SendUntilDone repeatedly calls Send until stopReason != "max_tokens".
	// Returns the complete accumulated output.
	SendUntilDone(text string) (reply, stopReason string, inputTokens, outputTokens int, err error)

	// SendStreamingUntilDone combines streaming with auto-continuation.
	SendStreamingUntilDone(text string, callback StreamCallback) (reply, stopReason string, inputTokens, outputTokens int, err error)

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
