package unit

import (
	"testing"

	list "github.com/jlrosende/project-manager/pkg/ui/list"
)

func TestRenderNames(t *testing.T) {
	out := list.RenderNames("", []string{"dev", "pre"})
	if len(out) == 0 {
		t.Fatal("expected non-empty render output")
	}
}
