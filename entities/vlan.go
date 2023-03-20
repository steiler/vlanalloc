package entities

import (
	"fmt"

	. "github.com/steiler/vlanalloc/utils"
)

type Vlan struct {
	id          int
	targetLevel TargetLevel
}

func NewVlan(id int, t TargetLevel) *Vlan {
	return &Vlan{
		id:          id,
		targetLevel: t,
	}
}

func (v *Vlan) ID() int {
	return v.id
}

func (v *Vlan) Identifier() string {
	return fmt.Sprintf("%s%s%d", v.targetLevel.Identifier(), IdentifierSep, v.id)
}

func (v *Vlan) String(indent int) string {
	white := GetWhitespaces(indent)
	result := fmt.Sprintf("%sVlanId: %d, VlanIdentifier: %s\n", white, v.id, v.Identifier())
	return result
}

func (v *Vlan) StringOneLine(indent int) string {
	white := GetWhitespaces(indent)
	return fmt.Sprintf("%sVlan: %s", white, v.Identifier())
}
