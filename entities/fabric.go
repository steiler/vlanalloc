package entities

import (
	"fmt"

	. "github.com/steiler/vlanalloc/enum"
	. "github.com/steiler/vlanalloc/utils"
)

type Fabric struct {
	global  *Global
	name    string
	routers map[string]*Router
	vlans   map[int]*Vlan
}

func NewFabric(name string, global *Global) *Fabric {
	f := &Fabric{
		name:    name,
		routers: map[string]*Router{},
		global:  global,
		vlans:   map[int]*Vlan{},
	}
	global.AddFabric(f)
	return f
}

func (f *Fabric) Name() string {
	return f.name
}

func (f *Fabric) NewRouter(name string) *Router {
	r := NewRouter(name, f)
	f.routers[r.name] = r
	return r
}

func (f *Fabric) Identifier() string {
	return fmt.Sprintf("%s%s%s", f.global.Identifier(), IdentifierSep, f.name)
}

func (f *Fabric) String(indent int) string {
	white := GetWhitespaces(indent)
	result := fmt.Sprintf("%sFabric: %s\n", white, f.name)
	for _, v := range f.vlans {
		result = result + fmt.Sprintf("%s+ VLAN: %s\n", white, v.Identifier())
	}
	for _, r := range f.routers {
		result = result + r.String(indent+PerLevelIndent)
	}
	return result
}

// GetVlanIndex returns the reference to the overall VLAN Index
func (f *Fabric) GetVlanIndex() *VlanIndex {
	return f.global.GetVlanIndex()
}

// GetAssignedVlanIDsUp retrieves all the assigned Vlan IDs from this level and further up
func (f *Fabric) GetAssignedVlanIDsUp(vlanMap map[int]struct{}) {
	// add the fabric assigned vlans
	for k := range f.vlans {
		vlanMap[k] = struct{}{}
	}
	f.global.GetAssignedVlanIDsUp(vlanMap)
}

// GetAssignedVlanIDsDown retrieves all the assigned Vlan IDs from this level and further down
func (f *Fabric) GetAssignedVlanIDsDown(vlanMap map[int]struct{}) {
	// add the fabrics assigned vlans
	for k, _ := range f.vlans {
		vlanMap[k] = struct{}{}
	}
	// add the interface assigned vlans
	for _, v := range f.routers {
		v.GetAssignedVlanIDsDown(vlanMap)
	}
}

func (f *Fabric) AssignVlan(router, interf string, scope Scope, vlanId int) (string, error) {
	switch scope {
	case Scope_Fabric:
		vidMapDown := map[int]struct{}{}
		f.GetAssignedVlanIDsDown(vidMapDown)
		if _, exists := vidMapDown[vlanId]; exists {
			return "", fmt.Errorf("VlanId %d already exists in or under %s context", vlanId, Scope_Fabric)
		}
		vidMapUp := map[int]struct{}{}
		f.GetAssignedVlanIDsUp(vidMapUp)
		if _, exists := vidMapUp[vlanId]; exists {
			return "", fmt.Errorf("VlanId %d already exists in or above %s context", vlanId, Scope_Fabric)
		}
		v := NewVlan(vlanId, f)
		f.vlans[v.id] = v
		f.GetVlanIndex().AddVlan(v)
		return v.Identifier(), nil
	default:
		var r *Router
		var exists bool
		if r, exists = f.routers[router]; !exists {
			r = f.NewRouter(router)
		}
		return r.AssignVlan(interf, scope, vlanId)
	}
}

func (f *Fabric) Scope() Scope {
	return Scope_Fabric
}

func (f *Fabric) GenerateVlanIdentifier(v *Vlan) string {
	return fmt.Sprintf("%s%s%d", f.Identifier(), IdentifierSep, v.id)
}
