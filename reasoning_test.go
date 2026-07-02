package llmapi

import "testing"

// TestParseReasoningEffort pins ParseReasoningEffort as the inverse of
// ReasoningEffort.String across every level, plus case-insensitivity and the
// error on an unrecognized level.
func TestParseReasoningEffort(t *testing.T) {
	for _, e := range []ReasoningEffort{ReasoningOff, ReasoningLow, ReasoningMedium, ReasoningHigh, ReasoningXHigh, ReasoningMax} {
		got, err := ParseReasoningEffort(e.String())
		if err != nil {
			t.Errorf("ParseReasoningEffort(%q): unexpected error %v", e.String(), err)
			continue
		}
		if got != e {
			t.Errorf("ParseReasoningEffort(%q) = %v, want %v", e.String(), got, e)
		}
	}

	// Case-insensitive: a CLI/user typing "HIGH" resolves to ReasoningHigh.
	if got, err := ParseReasoningEffort("HIGH"); err != nil || got != ReasoningHigh {
		t.Errorf(`ParseReasoningEffort("HIGH") = %v, %v; want %v, nil`, got, err, ReasoningHigh)
	}

	// Unrecognized input is an error (the flag layer rejects it loudly).
	if _, err := ParseReasoningEffort("bogus"); err == nil {
		t.Error(`ParseReasoningEffort("bogus"): want error, got nil`)
	}
}

// TestReasoningEffort_XHigh pins the xhigh tier: it slots between high and max
// (Anthropic's own effort vocabulary, recommended for the most demanding
// coding/agentic work), stringifies to "xhigh", and parses back. The ordering
// matters — effort levels are comparable, and ParseReasoningEffort's
// ReasoningOff..ReasoningMax sweep only covers every level if xhigh sits
// inside that range.
func TestReasoningEffort_XHigh(t *testing.T) {
	if got := ReasoningXHigh.String(); got != "xhigh" {
		t.Errorf(`ReasoningXHigh.String() = %q, want "xhigh"`, got)
	}
	if got, err := ParseReasoningEffort("xhigh"); err != nil || got != ReasoningXHigh {
		t.Errorf(`ParseReasoningEffort("xhigh") = %v, %v; want %v, nil`, got, err, ReasoningXHigh)
	}
	if !(ReasoningHigh < ReasoningXHigh && ReasoningXHigh < ReasoningMax) {
		t.Errorf("ordering: want ReasoningHigh < ReasoningXHigh < ReasoningMax, got %d, %d, %d",
			ReasoningHigh, ReasoningXHigh, ReasoningMax)
	}
}
