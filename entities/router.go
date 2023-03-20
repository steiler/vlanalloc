package entities

import (
	"fmt"

	. "github.com/steiler/vlanalloc/enum"
	. "github.com/steiler/vlanalloc/utils"
)

type Router struct {
	name         string             // the identifying name of this router
	fabric       *Fabric            // the fabric this router belongs to
	interfces    map[string]*Interf // interfaces indexed by their name
	vlans        map[int]*Vlan      // vlans indexed by the Vlan ID
	bridgeDomain *BridgeDomain      // bridgedomains indexed by their name
}

func NewRouter(name string, fabric *Fabric) *Router {
	return &Router{
		name:      name,
		interfces: map[string]*Interf{},
		vlans:     map[int]*Vlan{},
		fabric:    fabric,
	}
}

func (r *Router) _addBridgeDomain(bd *BridgeDomain) error {
	r.bridgeDomain = bd
	return nil
}

func (r *Router) AddInterface(ifname string) *Interf {
	// TODO: Add error in case of existence
	i := NewInterf(ifname, r)
	r.interfces[i.Name()] = i
	return i
}

func (r *Router) GetInterface(name string) *Interf {
	return r.interfces[name]
}

func (r *Router) Identifier() string {
	return fmt.Sprintf("%s%s%s", r.fabric.Identifier(), IdentifierSep, r.name)
}

func (r *Router) String(indent int) string {
	white := GetWhitespaces(indent)
	result := fmt.Sprintf("%sRouter: %s\n", white, r.name)
	if r.bridgeDomain != nil {
		result = result + GetWhitespaces(indent+PerLevelIndent) + "*" + r.bridgeDomain.StringOneLine(0) + "\n"
	}
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
	// add the bridgedomain vlans
	if r.bridgeDomain != nil {
		r.bridgeDomain.GetAssignedVlanIDs(vlanMap)
	}
}

// GetVlanIndex returns the reference to the overall VLAN Index
func (r *Router) GetVlanIndex() *VlanIndex {
	return r.fabric.GetVlanIndex()
}

func (r *Router) GenerateVlanIdentifier(v *Vlan) string {
	return fmt.Sprintf("%s%s%d", r.Identifier(), IdentifierSep, v.id)
}

// GetAssignedVlans use x to provide the result struct
func (r *Router) GetAssignedVlans(x map[int]struct{}) {
	r.GetAssignedVlanIDsDown(x)
	r.GetAssignedVlanIDsUp(x)
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

// StringOneLine returns a string representation of the entity without a trailing newline
func (r *Router) StringOneLine(indent int) string {
	white := GetWhitespaces(indent)
	return fmt.Sprintf("%sRouter: %s", white, r.name)
}
