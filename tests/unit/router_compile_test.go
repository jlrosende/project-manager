//go:build unit
// +build unit

package unit_test

import (
	"testing"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui/router"
)

func TestRouterCompile(_ *testing.T) {
	r := router.New(router.RouteProjects)
	_ = r.Current()
	r.NavigateTo(router.RouteProjectForm, nil)
	r.Back()
}
