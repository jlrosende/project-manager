package router

import "testing"

func TestRouterCompile(_ *testing.T) {
	r := New(RouteProjects)
	_ = r.Current()
	r.NavigateTo(RouteProjectForm, nil)
	r.Back()
}
