package entities

import (
	"fmt"
	//. "github.com/steiler/vlanalloc/utils"
)

type BridgeDomainIndex struct {
	bds map[string]*BridgeDomain
}

func NewBridgeDomainIndex() *BridgeDomainIndex {
	return &BridgeDomainIndex{
		bds: map[string]*BridgeDomain{},
	}
}

func (bdi *BridgeDomainIndex) String(indent int) string {
	var result string
	for _, bd := range bdi.bds {
		result = result + fmt.Sprintf("%s", bd.String(0))
	}

	return result
}

func (bdi *BridgeDomainIndex) AddBridgeDomain(bd *BridgeDomain) error {
	if _, exists := bdi.bds[bd.Identifier()]; exists {
		return fmt.Errorf("BridgeDomain with identifier %q already exists", bd.Identifier())
	}
	bdi.bds[bd.Identifier()] = bd
	return nil
}

func (bdi *BridgeDomainIndex) GetBridgeDomain(name string) *BridgeDomain {
	return bdi.bds[name]
}
