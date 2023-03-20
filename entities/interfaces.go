package entities

import (
	"fmt"

	. "github.com/steiler/vlanalloc/enum"
	. "github.com/steiler/vlanalloc/utils"
)

type Interf struct {
	name         string
	router       *Router
	vlans        map[int]*Vlan // vlans indexed by their vlan id
	bridgeDomain *BridgeDomain // bridgedomains indexed by their name
}

func NewInterf(name string, router *Router) *Interf {
	return &Interf{
		name:   name,
		router: router,
		vlans:  map[int]*Vlan{},
	}
}

func (i *Interf) Name() string {
	return i.name
}

func (i *Interf) Identifier() string {
	return fmt.Sprintf("%s%s%s", i.router.Identifier(), IdentifierSep, i.name)
}

func (i *Interf) String(indent int) string {
	white := GetWhitespaces(indent)
	result := fmt.Sprintf("%sInterface: %s\n", white, i.name)
	if i.bridgeDomain != nil {
		result = result + GetWhitespaces(indent+PerLevelIndent) + "*" + i.bridgeDomain.StringOneLine(0) + "\n"
	}
	for _, v := range i.vlans {
		result = result + fmt.Sprintf("%s+ VLAN: %s\n", white, v.Identifier())
	}

	return result
}

// GetAssignedVlanIDsUp retrieves all the assigned Vlan IDs from this level and further up
func (i *Interf) GetAssignedVlanIDsUp(vlanMap map[int]struct{}) {
	// add the interface assigned vlans
	for k := range i.vlans {
		vlanMap[k] = struct{}{}
	}
	i.router.GetAssignedVlanIDsUp(vlanMap)
}

// GetAssignedVlanIDs retrieves all the assigned Vlan IDs from this level and further down
func (i *Interf) GetAssignedVlanIDsDown(vlanMap map[int]struct{}) {
	for k, _ := range i.vlans {
		vlanMap[k] = struct{}{}
	}
	if i.bridgeDomain != nil {
		i.bridgeDomain.GetAssignedVlanIDs(vlanMap)
	}
}

// GetVlanIndex returns the reference to the overall VLAN Index
func (i *Interf) GetVlanIndex() *VlanIndex {
	return i.router.GetVlanIndex()
}

func (i *Interf) AssignVlan(scope Scope, vlanId int) (string, error) {

	// TODO: CHECK the UP stuff (GetAssignedVlanIDsUp) before the switch and before going down

	switch scope {
	case Scope_Interface:
		vidMap := map[int]struct{}{}
		i.GetAssignedVlanIDsDown(vidMap)
		if _, exists := vidMap[vlanId]; exists {
			return "", fmt.Errorf("VlanId %d already exists in or under %s context", vlanId, Scope_Interface)
		}
		vidMapUp := map[int]struct{}{}
		i.GetAssignedVlanIDsUp(vidMapUp)
		if _, exists := vidMapUp[vlanId]; exists {
			return "", fmt.Errorf("VlanId %d already exists in or above %s context", vlanId, Scope_Interface)
		}

		v := NewVlan(vlanId, i)
		i.vlans[v.id] = v
		i.GetVlanIndex().AddVlan(v)
		return v.Identifier(), nil
	default:
		return "", fmt.Errorf("Error, issues with scope, reached interface, scope does not match.")
	}
}

func (i *Interf) GenerateVlanIdentifier(v *Vlan) string {
	return fmt.Sprintf("%s%s%d", i.Identifier(), IdentifierSep, v.id)
}

func (i *Interf) Scope() Scope {
	return Scope_Interface
}

// StringOneLine returns a string representation of the entity without a trailing newline
func (i *Interf) StringOneLine(indent int) string {
	white := GetWhitespaces(indent)
	return fmt.Sprintf("%sInterface: %s:%s", white, i.router.name, i.name)
}

// GetAssignedVlans use x to provide the result struct
func (i *Interf) GetAssignedVlans(x map[int]struct{}) {
	i.GetAssignedVlanIDsDown(x)
	i.GetAssignedVlanIDsUp(x)
}

func (i *Interf) _setBridgeDomain(bd *BridgeDomain) error {
	i.bridgeDomain = bd
	return nil
}
