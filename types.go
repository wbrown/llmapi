// Package llmapi provides a unified interface for LLM providers.
// Both Anthropic and NovelAI implementations can implement this interface,
// allowing code to swap providers with minimal changes.
package llmapi

import "encoding/json"

// ==========================================================================
// Content Block Types
// ==========================================================================

// ContentType Identifies the type of content block.
type ContentType string

const (
	ContentTypeText       ContentType = "text"
	ContentTypeImage      ContentType = "image"
	ContentTypeToolUse    ContentType = "tool_use"
	ContentTypeToolResult ContentType = "tool_result"
	ContentTypeThinking   ContentType = "thinking"
	ContentTypeDocument   ContentType = "document"
)

// Role identifies the sender of a message.
type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

// MediaType identifies the MIME type of content.
type MediaType string

const (
	// Image types
	MediaTypePNG  MediaType = "image/png"
	MediaTypeJPEG MediaType = "image/jpeg"
	MediaTypeGIF  MediaType = "image/gif"
	MediaTypeWebP MediaType = "image/webp"

	// Document types
	MediaTypePDF MediaType = "application/pdf"
)

// ContentBlock represents a single block of content within a message.
// This is the provider-agnostic representation that both Anthropic
// and NovelAI implementations can work with.
type ContentBlock struct {
	Type       ContentType        `json:"type"`
	Text       string             `json:"text,omitempty"`
	Image      *ImageContent      `json:"image,omitempty"`
	ToolUse    *ToolUseContent    `json:"tool_use,omitempty"`
	ToolResult *ToolResultContent `json:"tool_result,omitempty"`
	Thinking   *ThinkingContent   `json:"thinking,omitempty"`
	Document   *DocumentContent   `json:"document,omitempty"`
}

// ==========================================================================
// Image Content
// ==========================================================================

// ImageContent represents image data within a content block.
type ImageContent struct {
	Source ImageSource `json:"source"`
}

// ImageSource contains the actual image data or reference.
type ImageSource struct {
	// Type is the source type: "base64" or "url".
	Type string `json:"type"`
	// MediaType is the MIME type of the image data (MediaTypePNG, MediaTypeJPEG, etc.)
	MediaType MediaType `json:"media_type"`
	// Data is the base64-encoded image data. (When Type is "base64".)
	Data string `json:"data,omitempty"`
	// URL is the URL of the image. (When Type is "url".)
	URL string `json:"url,omitempty"`
}

// ==========================================================================
// Tool Use Content
// ==========================================================================

// ToolUseContent represents a tool call made by the assistant.
type ToolUseContent struct {
	// ID is a unique identifier for this tool use, used to match results.
	ID string `json:"id"`
	// Name is the name of the tool being called.
	Name string `json:"name"`
	// Input is the tool arguments as JSON.
	Input json.RawMessage `json:"input"`
}

// ToolResultContent represents the result of a tool call.
type ToolResultContent struct {
	// ToolUseID matches the ID from the corresponding ToolUseContent.
	ToolUseID string `json:"tool_use_id"`
	// Content is the result of the tool execution.
	Content string `json:"content"`
	// IsError indicates if the tool execution failed.
	IsError bool `json:"is_error,omitempty"`
}

// ToolDefinition describes a tool that can be called by the assistant.
type ToolDefinition struct {
	// Name is the name of the tool.
	Name string `json:"name"`
	// Description is a short description of the tool.
	Description string `json:"description"`
	// InputSchema is a JSON Schema describing the tool's input parameters.
	InputSchema json.RawMessage `json:"input_schema"`
}

// ==========================================================================
// Thinking Content
// ==========================================================================

// ThinkingContent represents internal reasoning/chain-of-thought.
type ThinkingContent struct {
	// Thinking is the reasoning text.
	Thinking string `json:"thinking"`
	// Signature is used for verification (Anthropic-specific, optional).
	Signature string `json:"signature,omitempty"`
}

// ==========================================================================
// Document Content
// ==========================================================================

// DocumentContent represents embedded documents (PDFs, etc.).
type DocumentContent struct {
	// Source specifies how the document is provided
	Source DocumentSource `json:"source"`
	// Title is an optional title for the document.
	Title string `json:"title,omitempty"`
}

// DocumentSource contains document data.
type DocumentSource struct {
	// Type is the source type: "base64"
	Type string `json:"type"`
	// MediaType is the MIME type of the document data (MediaTypePDF, etc.)
	MediaType MediaType `json:"media_type"`
	// Data is the base64-encoded document data.
	Data string `json:"data"`
}

// ==========================================================================
// Rich Message and Response Types
// ==========================================================================

// RichMessage represents a message with multiple content blocks.
// This extends the simple Message type for advanced use cases.
type RichMessage struct {
	Role    Role           `json:"role"`
	Content []ContentBlock `json:"content"`
}

// ToMessage converts a RichMessage to a simple Message by extracting text.
// This is useful for providers that don't support rich content.
func (rm RichMessage) ToMessage() Message {
	var text string
	for _, block := range rm.Content {
		switch block.Type {
		case ContentTypeText:
			text += block.Text
		case ContentTypeThinking:
			if block.Thinking != nil {
				text += "<thinking>\n" + block.Thinking.Thinking + "\n</thinking>\n"
			}
		}
	}

	return Message{
		Role:    rm.Role,
		Content: text,
	}
}

// RichResponse contains the full response from a SendRich operation,
// including all content blocks, not just text.
type RichResponse struct {
	// Content contains all response content blocks.
	Content []ContentBlock `json:"content"`
	// StopReason indicates why the generation stopped.
	StopReason string `json:"stop_reason"`
	// InputTokens is the number of input tokens consumed.
	InputTokens int `json:"input_tokens"`
	// OutputTokens is the number of output tokens generated.
	OutputTokens int `json:"output_tokens"`
}

// Text returns the concatenated text from the response.
func (rr RichResponse) Text() string {
	var text string
	for _, block := range rr.Content {
		if block.Type == ContentTypeText {
			text += block.Text
		}
	}
	return text
}

// ThinkingText returns the concatenated thinking text from the response.
func (rr RichResponse) ThinkingText() string {
	var text string
	for _, block := range rr.Content {
		if block.Type == ContentTypeThinking && block.Thinking != nil {
			text += block.Thinking.Thinking
		}
	}
	return text
}

// ToolUses returns all tool uses from the response.
func (rr RichResponse) ToolUses() []ToolUseContent {
	var uses []ToolUseContent
	for _, block := range rr.Content {
		if block.Type == ContentTypeToolUse && block.ToolUse != nil {
			uses = append(uses, *block.ToolUse)
		}
	}
	return uses
}

// HasToolUse returns true if the response contains tool use requests.
func (rr RichResponse) HasToolUse() bool {
	return len(rr.ToolUses()) > 0
}

// ==========================================================================
// Helper Constructors
// ==========================================================================

// NewTextBlock creates a text content block.
func NewTextBlock(text string) ContentBlock {
	return ContentBlock{
		Type: ContentTypeText,
		Text: text,
	}
}

// NewImageBlock creates an image content block.
func NewImageBlock(mediaType MediaType, base64Data string) ContentBlock {
	return ContentBlock{
		Type: ContentTypeImage,
		Image: &ImageContent{
			Source: ImageSource{
				Type:      "base64",
				MediaType: mediaType,
				Data:      base64Data,
			},
		},
	}
}

// NewImageBlockFromURL creates an image content block from a URL.
func NewImageBlockFromURL(mediaType MediaType, url string) ContentBlock {
	return ContentBlock{
		Type: ContentTypeImage,
		Image: &ImageContent{
			Source: ImageSource{
				Type:      "url",
				MediaType: mediaType,
				URL:       url,
			},
		},
	}
}

// NewToolResultBlock creates a tool result content block.
func NewToolResultBlock(toolUseID string, content string, isError bool) ContentBlock {
	return ContentBlock{
		Type: ContentTypeToolResult,
		ToolResult: &ToolResultContent{
			ToolUseID: toolUseID,
			Content:   content,
			IsError:   isError,
		},
	}
}

// NewThinkingBlock creates a thinking content block.
func NewThinkingBlock(thinking string) ContentBlock {
	return ContentBlock{
		Type: ContentTypeThinking,
		Thinking: &ThinkingContent{
			Thinking: thinking,
		},
	}
}

// ==========================================================================
// Compatibility Detection
// ==========================================================================

// Capabilities describes what features a conversation implementation
// supports. This allows code to gracefully handle provider differences.
type Capabilities struct {
	SupportsImages      bool
	SupportsDocuments   bool
	SupportsToolUse     bool
	SupportsThinking    bool
	SupportsStreaming   bool
	MaxImageSize        int64    // bytes, 0 = no limit
	SupportedImageTypes []string // eg. ["image/png", "image/jpeg"]
}

// Message represents a single message in a conversation.
type Message struct {
	Role    Role   `json:"role"`    // RoleUser, RoleAssistant, RoleSystem
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
