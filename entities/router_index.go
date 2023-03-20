package entities

import (
	"fmt"

	. "github.com/steiler/vlanalloc/utils"
)

type RouterIndex struct {
	routers map[string]*Router
}

func NewRouterIndex() *RouterIndex {
	return &RouterIndex{
		routers: map[string]*Router{},
	}
}

func (ri *RouterIndex) String(indent int) string {
	white := GetWhitespaces(indent)
	result := ""
	for _, r := range ri.routers {
		result = result + fmt.Sprintf("%s", white+r.String(indent))
	}

	return result
}

func (ri *RouterIndex) AddRouter(r *Router) error {
	if _, exists := ri.routers[r.Identifier()]; exists {
		return fmt.Errorf("Router with identifier %q already exists", r.Identifier())
	}
	// TODO: Should maybe be r.identifier()
	ri.routers[r.name] = r
	return nil
}

func (ri *RouterIndex) GetRouter(name string) *Router {
	return ri.routers[name]
}
