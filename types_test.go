package llmapi

import "testing"

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
	block := NewImageBlock("image/png", "base64data")

	if block.Type != ContentTypeImage {
		t.Errorf("Expected type %s, got %s", ContentTypeImage, block.Type)
	}
	if block.Image == nil {
		t.Fatal("Expected Image to not be nil")
	}
	if block.Image.Source.Type != "base64" {
		t.Errorf("Expected source type 'base64', got '%s'", block.Image.Source.Type)
	}
	if block.Image.Source.MediaType != "image/png" {
		t.Errorf("Expected media type 'image/png', got '%s'", block.Image.Source.MediaType)
	}
	if block.Image.Source.Data != "base64data" {
		t.Errorf("Expected data 'base64data', got '%s'", block.Image.Source.Data)
	}
}

// TestNewImageBlockFromURL tests the NewImageBlockFromURL helper constructor.
func TestNewImageBlockFromURL(t *testing.T) {
	block := NewImageBlockFromURL("image/jpeg", "https://example.com/image.jpg")

	if block.Type != ContentTypeImage {
		t.Errorf("Expected type %s, got %s", ContentTypeImage, block.Type)
	}
	if block.Image == nil {
		t.Fatal("Expected Image to not be nil")
	}
	if block.Image.Source.Type != "url" {
		t.Errorf("Expected source type 'url', got '%s'", block.Image.Source.Type)
	}
	if block.Image.Source.MediaType != "image/jpeg" {
		t.Errorf("Expected media type 'image/jpeg', got '%s'", block.Image.Source.MediaType)
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
			Role: "assistant",
			Content: []ContentBlock{
				NewTextBlock("Hello "),
				NewTextBlock("world!"),
			},
		}

		msg := rm.ToMessage()
		if msg.Role != "assistant" {
			t.Errorf("Expected role 'assistant', got '%s'", msg.Role)
		}
		if msg.Content != "Hello world!" {
			t.Errorf("Expected content 'Hello world!', got '%s'", msg.Content)
		}
	})

	t.Run("WithThinking", func(t *testing.T) {
		rm := RichMessage{
			Role: "assistant",
			Content: []ContentBlock{
				NewThinkingBlock("My reasoning"),
				NewTextBlock("My response"),
			},
		}

		msg := rm.ToMessage()
		expected := "<thinking>\nMy reasoning\n</thinkingCan>\nMy response"
		if msg.Content != expected {
			t.Errorf("Expected content '%s', got '%s'", expected, msg.Content)
		}
	})

	t.Run("MixedContent", func(t *testing.T) {
		rm := RichMessage{
			Role: "user",
			Content: []ContentBlock{
				NewTextBlock("Look at this image:"),
				NewImageBlock("image/png", "base64data"), // Should be ignored
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
