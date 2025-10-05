package delete

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

func TestDefaultConfirmationFallbackPromptsViaStdio(t *testing.T) {
	in := strings.NewReader("yes\n")

	var out bytes.Buffer

	confirm := defaultConfirmation(in, &out, nil)

	ok, err := confirm(context.Background(), ConfirmationRequest{
		Target: domain.ProjectIdentifier{Name: "sample-app"},
		Scope:  domain.DeleteScopeMetadata,
	})
	if err != nil {
		t.Fatalf("confirmation returned error: %v", err)
	}

	if !ok {
		t.Fatalf("expected confirmation to succeed when user accepts")
	}

	prompt := out.String()
	if !strings.Contains(prompt, "Delete sample-app") {
		t.Fatalf("expected prompt to include target label, got %q", prompt)
	}
}
