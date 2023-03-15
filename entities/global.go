package entities

import (
	"fmt"

	. "github.com/steiler/vlanalloc/enum"
	. "github.com/steiler/vlanalloc/utils"
)

type Global struct {
	fabrics   map[string]*Fabric
	vlans     map[int]*Vlan
	vlanIndex *VlanIndex
}

func NewGlobal() *Global {
	return &Global{
		fabrics:   map[string]*Fabric{},
		vlans:     map[int]*Vlan{},
		vlanIndex: NewVlanIndex(),
	}
}

func (g *Global) NewFabric(name string) *Fabric {
	fab := NewFabric(name, g)
	g.fabrics[fab.name] = fab
	return fab
}

func (g *Global) GetVlanIndex() *VlanIndex {
	return g.vlanIndex
}

func (g *Global) AddFabric(fab *Fabric) {
	g.fabrics[fab.name] = fab
}

func (g *Global) String(indent int) string {
	white := GetWhitespaces(indent)

	result := white + "Global:\n"
	for _, v := range g.vlans {
		result = result + fmt.Sprintf("%s+ VLAN: %s\n", white, v.Identifier())
	}
	for _, v := range g.fabrics {
		result = result + v.String(indent+PerLevelIndent)
	}
	return result
}

// GetAssignedVlanIDsUp retrieves all the assigned Vlan IDs from this level and further up
func (g *Global) GetAssignedVlanIDsUp(vlanMap map[int]struct{}) {
	// add the globally assigned vlans
	for k := range g.vlans {
		vlanMap[k] = struct{}{}
	}
}

// GetAssignedVlanIDsDown retrieves all the assigned Vlan IDs from this level and further down
func (g *Global) GetAssignedVlanIDsDown(vlanMap map[int]struct{}) {
	// add the globally assigned vlans
	for k := range g.vlans {
		vlanMap[k] = struct{}{}
	}
	// add the interface assigned vlans
	for _, f := range g.fabrics {
		f.GetAssignedVlanIDsDown(vlanMap)
	}
}

func (g *Global) Scope() Scope {
	return Scope_Global
}

func (g *Global) Identifier() string {
	return string(g.Scope())
}

func (g *Global) GenerateVlanIdentifier(v *Vlan) string {
	return fmt.Sprintf("%s%s%d", g.Scope(), IdentifierSep, v.id)
}

func (g *Global) AssignVlan(fabric, router, interf string, scope Scope, vlanId int) (string, error) {
	switch scope {
	case Scope_Global:
		vidMapDown := map[int]struct{}{}
		g.GetAssignedVlanIDsDown(vidMapDown)
		if _, exists := vidMapDown[vlanId]; exists {
			return "", fmt.Errorf("VlanId %d already exists in or under %s context", vlanId, Scope_Global)
		}
		vidMapUp := map[int]struct{}{}
		g.GetAssignedVlanIDsUp(vidMapUp)
		if _, exists := vidMapUp[vlanId]; exists {
			return "", fmt.Errorf("VlanId %d already exists in or above %s context", vlanId, Scope_Global)
		}
		v := NewVlan(vlanId, g)
		g.vlans[v.id] = v
		g.vlanIndex.AddVlan(v)
		return v.Identifier(), nil
	default:
		var fab *Fabric
		var exists bool
		if fab, exists = g.fabrics[fabric]; !exists {
			fab = g.NewFabric(fabric)
		}
		return fab.AssignVlan(router, interf, scope, vlanId)
	}
}
