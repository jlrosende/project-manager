package list

func RenderNames(title string, names []string) string {
	it := []Item{}
	for _, n := range names {
		it = append(it, Item{Name: n})
	}
	return NewList(title, it).View()
}
