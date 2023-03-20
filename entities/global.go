package entities

import (
	"fmt"

	. "github.com/steiler/vlanalloc/enum"
	. "github.com/steiler/vlanalloc/utils"
)

type Global struct {
	fabrics     map[string]*Fabric
	vlans       map[int]*Vlan
	vlanIndex   *VlanIndex
	bdIndex     *BridgeDomainIndex
	routerIndex *RouterIndex
}

func NewGlobal() *Global {

	g := &Global{
		fabrics:     map[string]*Fabric{},
		vlans:       map[int]*Vlan{},
		vlanIndex:   NewVlanIndex(),
		bdIndex:     NewBridgeDomainIndex(),
		routerIndex: NewRouterIndex(),
	}

	g.NewFabric("Default")
	return g
}

func (g *Global) NewFabric(name string) *Fabric {
	fab := NewFabric(name, g)
	g.fabrics[fab.name] = fab
	return fab
}

func (g *Global) GetVlanIndex() *VlanIndex {
	return g.vlanIndex
}

func (g *Global) GetBridgeDomainIndex() *BridgeDomainIndex {
	return g.bdIndex
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
	for _, f := range g.fabrics {
		result = result + f.String(indent+PerLevelIndent)
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

func (g *Global) AssignVlanBD(bdname string, vlanId int) (string, error) {
	var bd *BridgeDomain
	if bd = g.bdIndex.GetBridgeDomain(bdname); bd == nil {
		bd = NewBridgeDomain(bdname, g.vlanIndex)
		err := g.bdIndex.AddBridgeDomain(bd)
		if err != nil {
			return "", err
		}
	}
	vlan, err := bd.AssignVlan(vlanId)
	if err != nil {
		return "", err
	}
	return vlan.Identifier(), nil
}

func (g *Global) AssignInterfaceToBD(router, ifname, bridgedomain string) error {
	// Get router from index
	r := g.routerIndex.GetRouter(router)
	// if router does not exist create it
	if r == nil {
		r = g.fabrics["Default"].NewRouter(router)
		// add new router to index
		err := g.routerIndex.AddRouter(r)
		if err != nil {
			return err
		}
	}
	// try get interface from router
	interf := r.GetInterface(ifname)
	if interf == nil {
		// if interface does not exist, create it
		interf = r.AddInterface(ifname)
	}

	// get the bridgedomain
	bd := g.bdIndex.GetBridgeDomain(bridgedomain)
	if bd == nil {
		// if it does not exist, create it
		bd = NewBridgeDomain(bridgedomain, g.vlanIndex)
		// add bridgedomain to index
		err := g.bdIndex.AddBridgeDomain(bd)
		if err != nil {
			return err
		}
	}
	// finally add the interface to the bridgedomain
	return bd.AddInterface(interf)
}

func (g *Global) AssignRouterToBD(router, bridgedomain string) error {
	// Get router from index
	r := g.routerIndex.GetRouter(router)
	// if router does not exist create it
	if r == nil {
		// TODO: For now we just take the "Default" Fabric
		r = g.fabrics["Default"].NewRouter(router)
		// add new router to index
		err := g.routerIndex.AddRouter(r)
		if err != nil {
			return err
		}
	}

	// get the bridgedomain
	bd := g.bdIndex.GetBridgeDomain(bridgedomain)
	if bd == nil {
		// if it does not exist, create it
		bd = NewBridgeDomain(bridgedomain, g.vlanIndex)
		// add bridgedomain to index
		err := g.bdIndex.AddBridgeDomain(bd)
		if err != nil {
			return err
		}
	}
	// finally add the router to the bridgedomain
	return bd.AddRouter(r)
}

func (g *Global) GetFreeBDVlan(bdname string) (*Vlan, error) {
	// get the bridgedomain
	bd := g.bdIndex.GetBridgeDomain(bdname)
	if bd == nil {
		// if it does not exist, create it
		bd = NewBridgeDomain(bdname, g.vlanIndex)
		// add bridgedomain to index
		err := g.bdIndex.AddBridgeDomain(bd)
		if err != nil {
			return nil, err
		}
	}
	return bd.GetFreeVlan()
}
