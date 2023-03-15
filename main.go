package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/steiler/vlanalloc/entities"
	. "github.com/steiler/vlanalloc/enum"
)

func main() {

	G := entities.NewGlobal()

	/////////////////////////////////////

	id, err := G.AssignVlan("", "", "", Scope_Global, 5)
	if err != nil {
		log.Errorf("assignment failed with %v", err)
	}
	id, err = G.AssignVlan("", "", "", Scope_Global, 5)
	if err != nil {
		log.Errorf("assignment failed with %v", err)
	}

	id, err = G.AssignVlan("", "", "", Scope_Global, 7)
	if err != nil {
		log.Errorf("assignment failed with %v", err)
	}

	id, err = G.AssignVlan("MyFabric", "", "", Scope_Fabric, 7)
	if err != nil {
		log.Errorf("assignment failed with %v", err)
	}

	id, err = G.AssignVlan("MyFabric", "", "", Scope_Fabric, 78)
	if err != nil {
		log.Errorf("assignment failed with %v", err)
	}

	id, err = G.AssignVlan("Fab2", "R5.1", "eth0", Scope_Interface, 7)
	if err != nil {
		log.Errorf("assignment failed with %v", err)
	}

	id, err = G.AssignVlan("Fab2", "R5.1", "eth0", Scope_Interface, 27)
	if err != nil {
		log.Errorf("assignment failed with %v", err)
	}

	id, err = G.AssignVlan("Fab2", "R5.1", "eth0", Scope_Router, 8)
	if err != nil {
		log.Errorf("assignment failed with %v", err)
	}

	////////////////////////////////////////////

	_ = id

	// print a treeview of the collected data
	fmt.Println(G.String(0))

	// print a list of vlans
	vidx := G.GetVlanIndex()
	fmt.Println(vidx.String())

}
