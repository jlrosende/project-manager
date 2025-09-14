package unit

import (
	"github.com/jlrosende/project-manager/pkg/ui/styles"
	"testing"
)

func TestDefaultStyleWidth(t *testing.T) {
	s := styles.DefaultStyle
	if s.GetWidth() == 0 {
		t.Fatal("expected non-zero width")
	}
}
