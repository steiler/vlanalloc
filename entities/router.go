package entities

import (
	"fmt"

	. "github.com/steiler/vlanalloc/enum"
	. "github.com/steiler/vlanalloc/utils"
)

type Router struct {
	name      string
	interfces map[string]*Interf
	fabric    *Fabric
	vlans     map[int]*Vlan
}

func NewRouter(name string, fabric *Fabric) *Router {
	return &Router{
		name:      name,
		interfces: map[string]*Interf{},
		vlans:     map[int]*Vlan{},
		fabric:    fabric,
	}
}

func (r *Router) AddInterface(ifname string) *Interf {
	i := NewInterf(ifname, r)
	r.interfces[i.Name()] = i
	return i
}

func (r *Router) Identifier() string {
	return fmt.Sprintf("%s%s%s", r.fabric.Identifier(), IdentifierSep, r.name)
}

func (r *Router) String(indent int) string {
	white := GetWhitespaces(indent)
	result := fmt.Sprintf("%sRouter: %s\n", white, r.name)
	for _, v := range r.vlans {
		result = result + fmt.Sprintf("%s+ VLAN: %s\n", white, v.Identifier())
	}
	for _, i := range r.interfces {
		result = result + i.String(indent+PerLevelIndent)
	}
	return result
}

// GetAssignedVlanIDsUp retrieves all the assigned Vlan IDs from this level and further up
func (r *Router) GetAssignedVlanIDsUp(vlanMap map[int]struct{}) {
	// add the router assigned vlans
	for k := range r.vlans {
		vlanMap[k] = struct{}{}
	}
	r.fabric.GetAssignedVlanIDsUp(vlanMap)
}

// GetAssignedVlanIDsDown retrieves all the assigned Vlan IDs from this level and further down
func (r *Router) GetAssignedVlanIDsDown(vlanMap map[int]struct{}) {
	// add the router assigned vlans
	for k, _ := range r.vlans {
		vlanMap[k] = struct{}{}
	}
	// add the interface assigned vlans
	for _, v := range r.interfces {
		v.GetAssignedVlanIDsDown(vlanMap)
	}
}

// GetVlanIndex returns the reference to the overall VLAN Index
func (r *Router) GetVlanIndex() *VlanIndex {
	return r.fabric.GetVlanIndex()
}

func (r *Router) GenerateVlanIdentifier(v *Vlan) string {
	return fmt.Sprintf("%s%s%d", r.Identifier(), IdentifierSep, v.id)
}

func (r *Router) AssignVlan(interf string, scope Scope, vlanId int) (string, error) {
	switch scope {
	case Scope_Router:
		vidMapDown := map[int]struct{}{}
		r.GetAssignedVlanIDsDown(vidMapDown)
		if _, exists := vidMapDown[vlanId]; exists {
			return "", fmt.Errorf("VlanId %d already exists in or under %s context", vlanId, Scope_Router)
		}
		vidMapUp := map[int]struct{}{}
		r.GetAssignedVlanIDsUp(vidMapUp)
		if _, exists := vidMapUp[vlanId]; exists {
			return "", fmt.Errorf("VlanId %d already exists in or above %s context", vlanId, Scope_Router)
		}
		v := NewVlan(vlanId, r)
		r.vlans[v.id] = v
		r.GetVlanIndex().AddVlan(v)
		return v.Identifier(), nil
	default:
		var i *Interf
		var exists bool
		if i, exists = r.interfces[interf]; !exists {
			i = r.AddInterface(interf)
		}
		return i.AssignVlan(scope, vlanId)
	}
}

func (r *Router) Scope() Scope {
	return Scope_Router
}
