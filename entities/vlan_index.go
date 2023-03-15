package entities

import "fmt"

type VlanIndex struct {
	vlans map[string]*Vlan
}

func NewVlanIndex() *VlanIndex {
	return &VlanIndex{
		vlans: map[string]*Vlan{},
	}
}

func (vi *VlanIndex) String() string {
	result := ""
	for _, v := range vi.vlans {
		result = result + fmt.Sprintf("%s", v.String(0))
	}

	return result
}

func (vi *VlanIndex) AddVlan(v *Vlan) error {
	if _, exists := vi.vlans[v.Identifier()]; exists {
		return fmt.Errorf("vlan with identifier %q already exists", v.Identifier())
	}
	vi.vlans[v.Identifier()] = v
	return nil
}
