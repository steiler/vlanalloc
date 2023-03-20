package entities

import (
	"fmt"

	. "github.com/steiler/vlanalloc/utils"
)

type Entry struct {
	parent    *Entry
	childs    map[string]*Entry
	name      string
	labels    map[string]string
	vlans     map[int]*Vlan
	vlanIndex *VlanIndex
}

func NewEntry(name string, labels map[string]string, parent *Entry) *Entry {
	e := &Entry{
		name:   name,
		labels: labels,
		parent: parent,
		childs: map[string]*Entry{},
		vlans:  map[int]*Vlan{},
	}
	return e
}

func (e *Entry) IsRoot() bool {
	return e.parent == nil
}

func (e *Entry) Name() string {
	return e.name
}

func (e *Entry) NewChildEntry(name string, labels map[string]string) *Entry {
	r := NewEntry(name, labels, e)
	e.childs[r.name] = r
	return r
}

func (e *Entry) String(indent int) string {
	white := GetWhitespaces(indent)
	// convert labels to printable string
	labelString := MapStringString2String(e.labels, LabelKVSep, LabelEntrySep)

	result := fmt.Sprintf("%sEntry: %s [%s]\n", white, e.name, labelString)
	for _, v := range e.vlans {
		result = result + fmt.Sprintf("%s+ VLAN: %s\n", white, v.Identifier())
	}
	for _, i := range e.childs {
		result = result + i.String(indent+PerLevelIndent)
	}
	return result
}

func (e *Entry) Identifier() string {
	return fmt.Sprintf("%d%s%s", e.GetLevel(), IdentifierSep, e.name)
}

func (e *Entry) GetLevel() int {
	if e.IsRoot() {
		return 0
	}
	return e.parent.GetLevel() + 1
}

// GetAssignedVlanIDsUp retrieves all the assigned Vlan IDs from this level and further up
func (e *Entry) GetAssignedVlanIDsUp(vlanMap map[int]struct{}) {
	// add the router assigned vlans
	for k := range e.vlans {
		vlanMap[k] = struct{}{}
	}
	e.parent.GetAssignedVlanIDsUp(vlanMap)
}

// GetAssignedVlanIDsDown retrieves all the assigned Vlan IDs from this level and further down
func (e *Entry) GetAssignedVlanIDsDown(vlanMap map[int]struct{}) {
	// add the router assigned vlans
	for k, _ := range e.vlans {
		vlanMap[k] = struct{}{}
	}
	// add the interface assigned vlans
	for _, v := range e.childs {
		v.GetAssignedVlanIDsDown(vlanMap)
	}
}
