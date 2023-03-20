package entities

import (
	"fmt"
	"math/rand"

	. "github.com/steiler/vlanalloc/enum"
	. "github.com/steiler/vlanalloc/utils"
)

type BridgeDomain struct {
	name       string
	labels     map[string]string
	vlans      map[int]*Vlan
	routers    map[string]*Router
	interfaces []*Interf
	vlanindex  *VlanIndex
}

func NewBridgeDomain(name string, vlanIndex *VlanIndex) *BridgeDomain {
	return &BridgeDomain{
		name:       name,
		labels:     map[string]string{},
		vlans:      map[int]*Vlan{},
		routers:    map[string]*Router{},
		interfaces: []*Interf{},
		vlanindex:  vlanIndex,
	}
}

func (b *BridgeDomain) AssignVlan(vlanId int) (*Vlan, error) {

	if _, exists := b.vlans[vlanId]; exists {
		return nil, fmt.Errorf("Vlan %d already assigned to bridgedomain %s", vlanId, b.name)
	}

	takenVlans := map[int]struct{}{}
	for _, r := range b.routers {
		r.GetAssignedVlans(takenVlans)
		if _, taken := takenVlans[vlanId]; taken {
			return nil, fmt.Errorf("False assignment of VlanID %d", vlanId)
		}
	}
	for _, i := range b.interfaces {
		i.GetAssignedVlans(takenVlans)
		if _, taken := takenVlans[vlanId]; taken {
			return nil, fmt.Errorf("False assignment of VlanID %d", vlanId)
		}
	}
	v := NewVlan(vlanId, b)
	b.vlans[v.id] = v
	b.vlanindex.AddVlan(v)
	return v, nil
}

func (b *BridgeDomain) GenerateVlanIdentifier(v *Vlan) string {
	return fmt.Sprintf("%s%s%d", b.Identifier(), IdentifierSep, v.id)
}

func (b *BridgeDomain) Identifier() string {
	return b.name
}

func (b *BridgeDomain) Scope() Scope {
	return Scope_BridgeDomain
}

// StringOneLine returns a string representation of the entity without a trailing newline
func (b *BridgeDomain) StringOneLine(indent int) string {
	white := GetWhitespaces(indent)
	labelsString := MapStringString2String(b.labels, LabelKVSep, LabelEntrySep)
	return fmt.Sprintf("%sBridgeDomain: %s [%s]", white, b.name, labelsString)
}

func (b *BridgeDomain) String(indent int) string {
	white := GetWhitespaces(indent)
	labelsString := MapStringString2String(b.labels, LabelKVSep, LabelEntrySep)
	result := fmt.Sprintf("%sBridgeDomain: %s [%s]\n", white, b.name, labelsString)
	for _, v := range b.vlans {
		result = result + fmt.Sprintf("%s%s\n", white, v.StringOneLine(indent+PerLevelIndent))
	}
	for _, r := range b.routers {
		result = result + fmt.Sprintf("%s%s\n", white, r.StringOneLine(indent+PerLevelIndent))
	}
	for _, i := range b.interfaces {
		result = result + fmt.Sprintf("%s%s\n", white, i.StringOneLine(indent+PerLevelIndent))
	}
	return result
}

func (b *BridgeDomain) StringBDCentered(indent int) string {
	white := GetWhitespaces(indent)
	labelsString := MapStringString2String(b.labels, LabelKVSep, LabelEntrySep)
	result := fmt.Sprintf("%sBridgeDomain: %s [%s]\n", white, b.name, labelsString)
	for _, v := range b.vlans {
		result = result + fmt.Sprintf("%s+ VLAN: %s\n", white, v.Identifier())
	}
	return result
}

func (b *BridgeDomain) AddInterface(interf *Interf) error {
	interfVlans := map[int]struct{}{}
	interf.GetAssignedVlans(interfVlans)

	for vid, _ := range b.vlans {
		if _, exists := interfVlans[vid]; exists {
			// maybe also check that if it is the same VID if it is also the same
			// vlan Object, then this operation should not cause an error
			return fmt.Errorf("VlanId %d is already taken on interface %s", vid, interf.Identifier())
		}
	}
	b.interfaces = append(b.interfaces, interf)
	return interf._setBridgeDomain(b)
}

func (b *BridgeDomain) AddRouter(router *Router) error {
	routerVlans := map[int]struct{}{}
	router.GetAssignedVlans(routerVlans)

	for vid, _ := range b.vlans {
		if _, exists := routerVlans[vid]; exists {
			// maybe also check that if it is the same VID if it is also the same
			// vlan Object, then this operation should not cause an error
			return fmt.Errorf("VlanId %d is already taken on router %s", vid, router.Identifier())
		}
	}

	b.routers[router.name] = router
	return router._addBridgeDomain(b)
}

func (b *BridgeDomain) GetAssignedVlanIDs(result map[int]struct{}) {
	for k := range b.vlans {
		result[k] = struct{}{}
	}
}

// GetAssignedVlanIDsRelated returns the list of vlan ids assigned to the bridge domain
// as well as all the related object to determine a list of unavailable vlanids
func (b *BridgeDomain) GetAssignedVlanIDsRelated(result map[int]struct{}) {
	b.GetAssignedVlanIDs(result)
	for _, i := range b.interfaces {
		i.GetAssignedVlans(result)
	}
	for _, r := range b.routers {
		r.GetAssignedVlans(result)
	}
}

func (b *BridgeDomain) GetFreeVlan() (*Vlan, error) {
	vlanMap := map[int]struct{}{}
	b.GetAssignedVlanIDsRelated(vlanMap)

	if len(vlanMap) >= VlanEnd-VlanStart {
		return nil, fmt.Errorf("no more assignable vlans available %d assigned in range %d-%d []", len(vlanMap), VlanStart, VlanEnd)
	}

	var rangeIndex int

	randOffset := rand.Intn(VlanEnd-1) + 1
	// start at a random offset between 1 and vlanEnd
	// the offset shifts the end to randoOffset + VlanEnd, while the index that is to be accessed is vlanId (randOffset + VlanStart) % (modulo) VlanEnd
	var vlanId int
	for rangeIndex = VlanStart; rangeIndex <= VlanEnd; rangeIndex++ {
		if _, exists := vlanMap[rangeIndex+randOffset]; exists {
			continue
		}
		vlanId = rangeIndex + randOffset
		break
	}

	v, err := b.AssignVlan(vlanId)
	return v, err
}
