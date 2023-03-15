package entities

import "github.com/steiler/vlanalloc/enum"

type TargetLevel interface {
	Identifier() string
	Scope() enum.Scope
	GenerateVlanIdentifier(*Vlan) string
}
