package router

type Route int

const (
	RouteProjects Route = iota
	RouteProjectForm
	RouteEnvVars
	RouteUpdateFlow
)

type Params map[string]string

type Router struct {
	current Route
	stack   []Route
	params  Params
}

func New(start Route) *Router {
	return &Router{current: start, stack: []Route{}, params: Params{}}
}

func (r *Router) Current() Route { return r.current }

func (r *Router) Params() Params { return r.params }

func (r *Router) NavigateTo(route Route, p Params) {
	r.stack = append(r.stack, r.current)
	r.current = route
	r.params = p
}

func (r *Router) Back() {
	if len(r.stack) == 0 {
		return
	}

	idx := len(r.stack) - 1
	r.current = r.stack[idx]
	r.stack = r.stack[:idx]
	r.params = nil
}
