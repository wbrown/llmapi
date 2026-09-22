package llmapi

import (
	"encoding/json"
	"testing"
)

// TestNewTextBlock tests the NewTextBlock helper constructor.
func TestNewTextBlock(t *testing.T) {
	block := NewTextBlock("Hello world")

	if block.Type != ContentTypeText {
		t.Errorf("Expected type %s, got %s", ContentTypeText, block.Type)
	}
	if block.Text != "Hello world" {
		t.Errorf("Expected text 'Hello world', got '%s'", block.Text)
	}
	// Other fields should be nil/empty
	if block.Image != nil {
		t.Error("Expected Image to be nil")
	}
	if block.ToolUse != nil {
		t.Error("Expected ToolUse to be nil")
	}
}

// TestNewImageBlock tests the NewImageBlock helper constructor.
func TestNewImageBlock(t *testing.T) {
	block := NewImageBlock(MediaTypePNG, "base64data")

	if block.Type != ContentTypeImage {
		t.Errorf("Expected type %s, got %s", ContentTypeImage, block.Type)
	}
	if block.Image == nil {
		t.Fatal("Expected Image to not be nil")
	}
	if block.Image.Source.Type != "base64" {
		t.Errorf("Expected source type 'base64', got '%s'", block.Image.Source.Type)
	}
	if block.Image.Source.MediaType != MediaTypePNG {
		t.Errorf("Expected media type '%s', got '%s'", MediaTypePNG, block.Image.Source.MediaType)
	}
	if block.Image.Source.Data != "base64data" {
		t.Errorf("Expected data 'base64data', got '%s'", block.Image.Source.Data)
	}
}

// TestNewImageBlockFromURL tests the NewImageBlockFromURL helper constructor.
func TestNewImageBlockFromURL(t *testing.T) {
	block := NewImageBlockFromURL(MediaTypeJPEG, "https://example.com/image.jpg")

	if block.Type != ContentTypeImage {
		t.Errorf("Expected type %s, got %s", ContentTypeImage, block.Type)
	}
	if block.Image == nil {
		t.Fatal("Expected Image to not be nil")
	}
	if block.Image.Source.Type != "url" {
		t.Errorf("Expected source type 'url', got '%s'", block.Image.Source.Type)
	}
	if block.Image.Source.MediaType != MediaTypeJPEG {
		t.Errorf("Expected media type '%s', got '%s'", MediaTypeJPEG, block.Image.Source.MediaType)
	}
	if block.Image.Source.URL != "https://example.com/image.jpg" {
		t.Errorf("Expected URL 'https://example.com/image.jpg', got '%s'", block.Image.Source.URL)
	}
}

// TestNewToolResultBlock tests the NewToolResultBlock helper constructor.
func TestNewToolResultBlock(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		block := NewToolResultBlock("tool_123", "Result content", false)

		if block.Type != ContentTypeToolResult {
			t.Errorf("Expected type %s, got %s", ContentTypeToolResult, block.Type)
		}
		if block.ToolResult == nil {
			t.Fatal("Expected ToolResult to not be nil")
		}
		if block.ToolResult.ToolUseID != "tool_123" {
			t.Errorf("Expected ToolUseID 'tool_123', got '%s'", block.ToolResult.ToolUseID)
		}
		if block.ToolResult.Content != "Result content" {
			t.Errorf("Expected content 'Result content', got '%s'", block.ToolResult.Content)
		}
		if block.ToolResult.IsError {
			t.Error("Expected IsError to be false")
		}
	})

	t.Run("Error", func(t *testing.T) {
		block := NewToolResultBlock("tool_456", "Error message", true)

		if !block.ToolResult.IsError {
			t.Error("Expected IsError to be true")
		}
	})
}

// TestNewThinkingBlock tests the NewThinkingBlock helper constructor.
func TestNewThinkingBlock(t *testing.T) {
	block := NewThinkingBlock("Let me think about this...")

	if block.Type != ContentTypeThinking {
		t.Errorf("Expected type %s, got %s", ContentTypeThinking, block.Type)
	}
	if block.Thinking == nil {
		t.Fatal("Expected Thinking to not be nil")
	}
	if block.Thinking.Thinking != "Let me think about this..." {
		t.Errorf("Expected thinking text, got '%s'", block.Thinking.Thinking)
	}
}

// TestRichMessageToMessage tests the RichMessage.ToMessage() conversion.
func TestRichMessageToMessage(t *testing.T) {
	t.Run("TextOnly", func(t *testing.T) {
		rm := RichMessage{
			Role: RoleAssistant,
			Content: []ContentBlock{
				NewTextBlock("Hello "),
				NewTextBlock("world!"),
			},
		}

		msg := rm.ToMessage()
		if msg.Role != RoleAssistant {
			t.Errorf("Expected role 'assistant', got '%s'", msg.Role)
		}
		if msg.Content != "Hello world!" {
			t.Errorf("Expected content 'Hello world!', got '%s'", msg.Content)
		}
	})

	t.Run("WithThinking", func(t *testing.T) {
		rm := RichMessage{
			Role: RoleAssistant,
			Content: []ContentBlock{
				NewThinkingBlock("My reasoning"),
				NewTextBlock("My response"),
			},
		}

		msg := rm.ToMessage()
		expected := "<thinking>\nMy reasoning\n</thinking>\nMy response"
		if msg.Content != expected {
			t.Errorf("Expected content '%s', got '%s'", expected, msg.Content)
		}
	})

	t.Run("MixedContent", func(t *testing.T) {
		rm := RichMessage{
			Role: RoleUser,
			Content: []ContentBlock{
				NewTextBlock("Look at this image:"),
				NewImageBlock(MediaTypePNG, "base64data"), // Should be ignored
				NewTextBlock(" What do you see?"),
			},
		}

		msg := rm.ToMessage()
		if msg.Content != "Look at this image: What do you see?" {
			t.Errorf("Expected only text content, got '%s'", msg.Content)
		}
	})
}

// TestRichResponseText tests the RichResponse.Text() method.
func TestRichResponseText(t *testing.T) {
	rr := RichResponse{
		Content: []ContentBlock{
			NewThinkingBlock("thinking"),
			NewTextBlock("Hello "),
			NewTextBlock("world"),
		},
	}

	text := rr.Text()
	if text != "Hello world" {
		t.Errorf("Expected 'Hello world', got '%s'", text)
	}
}

// TestRichResponseThinkingText tests the RichResponse.ThinkingText() method.
func TestRichResponseThinkingText(t *testing.T) {
	rr := RichResponse{
		Content: []ContentBlock{
			NewThinkingBlock("First thought"),
			NewTextBlock("response"),
			NewThinkingBlock("Second thought"),
		},
	}

	thinking := rr.ThinkingText()
	if thinking != "First thoughtSecond thought" {
		t.Errorf("Expected 'First thoughtSecond thought', got '%s'", thinking)
	}
}

// TestRichResponseToolUses tests the RichResponse.ToolUses() method.
func TestRichResponseToolUses(t *testing.T) {
	rr := RichResponse{
		Content: []ContentBlock{
			NewTextBlock("I'll use a tool"),
			{
				Type: ContentTypeToolUse,
				ToolUse: &ToolUseContent{
					ID:    "tool_1",
					Name:  "get_weather",
					Input: []byte(`{"location": "NYC"}`),
				},
			},
			{
				Type: ContentTypeToolUse,
				ToolUse: &ToolUseContent{
					ID:    "tool_2",
					Name:  "get_time",
					Input: []byte(`{}`),
				},
			},
		},
	}

	uses := rr.ToolUses()
	if len(uses) != 2 {
		t.Fatalf("Expected 2 tool uses, got %d", len(uses))
	}
	if uses[0].Name != "get_weather" {
		t.Errorf("Expected first tool 'get_weather', got '%s'", uses[0].Name)
	}
	if uses[1].Name != "get_time" {
		t.Errorf("Expected second tool 'get_time', got '%s'", uses[1].Name)
	}
}

// TestRichResponseHasToolUse tests the RichResponse.HasToolUse() method.
func TestRichResponseHasToolUse(t *testing.T) {
	t.Run("WithToolUse", func(t *testing.T) {
		rr := RichResponse{
			Content: []ContentBlock{
				{
					Type:    ContentTypeToolUse,
					ToolUse: &ToolUseContent{ID: "1", Name: "test"},
				},
			},
		}
		if !rr.HasToolUse() {
			t.Error("Expected HasToolUse() to return true")
		}
	})

	t.Run("WithoutToolUse", func(t *testing.T) {
		rr := RichResponse{
			Content: []ContentBlock{
				NewTextBlock("Just text"),
			},
		}
		if rr.HasToolUse() {
			t.Error("Expected HasToolUse() to return false")
		}
	})

	t.Run("Empty", func(t *testing.T) {
		rr := RichResponse{}
		if rr.HasToolUse() {
			t.Error("Expected HasToolUse() to return false for empty response")
		}
	})
}

// TestRichResponse_CarriesRequestAccount pins the account a RichResponse
// carries beside its content: the completion budget the request carried on
// the wire, the server's own finish reason, and the output tokens attributed
// per channel. The split's Known flag
// separates a server that attributed nothing from one that attributed zero
// reasoning tokens, and both readings survive a JSON round trip.
func TestRichResponse_CarriesRequestAccount(t *testing.T) {
	roundTrip := func(t *testing.T, sent RichResponse) RichResponse {
		encoded, err := json.Marshal(sent)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		var got RichResponse
		if err := json.Unmarshal(encoded, &got); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		return got
	}

	t.Run("Attributed", func(t *testing.T) {
		sent := RichResponse{
			Content:          []ContentBlock{NewTextBlock("answer")},
			StopReason:       "end_turn",
			FinishReason:     "stop",
			InputTokens:      29,
			OutputTokens:     553,
			CompletionBudget: 24576,
			OutputSplit:      OutputTokenSplit{Reasoning: 522, Content: 31, Known: true},
		}
		got := roundTrip(t, sent)
		if got.CompletionBudget != 24576 {
			t.Errorf("CompletionBudget = %d, want 24576", got.CompletionBudget)
		}
		if got.FinishReason != "stop" {
			t.Errorf("FinishReason = %q, want %q", got.FinishReason, "stop")
		}
		if got.StopReason != "end_turn" {
			t.Errorf("StopReason = %q, want %q", got.StopReason, "end_turn")
		}
		if got.OutputSplit != sent.OutputSplit {
			t.Errorf("OutputSplit = %+v, want %+v", got.OutputSplit, sent.OutputSplit)
		}
	})

	t.Run("Unattributed", func(t *testing.T) {
		got := roundTrip(t, RichResponse{OutputTokens: 553})
		if got.OutputSplit.Known {
			t.Errorf("OutputSplit = %+v, want an unknown split for a server that attributed no channel", got.OutputSplit)
		}
	})

	t.Run("KnownZeroReasoning", func(t *testing.T) {
		sent := RichResponse{
			OutputTokens: 31,
			OutputSplit:  OutputTokenSplit{Reasoning: 0, Content: 31, Known: true},
		}
		got := roundTrip(t, sent)
		if !got.OutputSplit.Known {
			t.Errorf("OutputSplit = %+v, want a known split: zero reasoning tokens the server attributed is a real zero", got.OutputSplit)
		}
		if got.OutputSplit != sent.OutputSplit {
			t.Errorf("OutputSplit = %+v, want %+v", got.OutputSplit, sent.OutputSplit)
		}
	})
}

// TestContentTypes tests that content type constants are correct.
func TestContentTypes(t *testing.T) {
	tests := []struct {
		ct       ContentType
		expected string
	}{
		{ContentTypeText, "text"},
		{ContentTypeImage, "image"},
		{ContentTypeToolUse, "tool_use"},
		{ContentTypeToolResult, "tool_result"},
		{ContentTypeThinking, "thinking"},
		{ContentTypeDocument, "document"},
	}

	for _, tt := range tests {
		if string(tt.ct) != tt.expected {
			t.Errorf("Expected ContentType %s, got %s", tt.expected, tt.ct)
		}
	}
}
