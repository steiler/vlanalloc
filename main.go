package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/steiler/vlanalloc/entities"
	"github.com/steiler/vlanalloc/enum"
	// . "github.com/steiler/vlanalloc/enum"
)

func main() {

	G := entities.NewGlobal()

	/////////////////////////////////////

	// id, err := G.AssignVlan("", "", "", Scope_Global, 5)
	// if err != nil {
	// 	log.Errorf("assignment failed with %v", err)
	// }
	// id, err = G.AssignVlan("", "", "", Scope_Global, 5)
	// if err != nil {
	// 	log.Errorf("assignment failed with %v", err)
	// }

	// id, err = G.AssignVlan("", "", "", Scope_Global, 7)
	// if err != nil {
	// 	log.Errorf("assignment failed with %v", err)
	// }

	// id, err = G.AssignVlan("MyFabric", "", "", Scope_Fabric, 7)
	// if err != nil {
	// 	log.Errorf("assignment failed with %v", err)
	// }

	// id, err = G.AssignVlan("MyFabric", "", "", Scope_Fabric, 78)
	// if err != nil {
	// 	log.Errorf("assignment failed with %v", err)
	// }

	// id, err = G.AssignVlan("Fab2", "R5.1", "eth0", Scope_Interface, 7)
	// if err != nil {
	// 	log.Errorf("assignment failed with %v", err)
	// }

	// id, err = G.AssignVlan("Fab2", "R5.1", "eth0", Scope_Interface, 27)
	// if err != nil {
	// 	log.Errorf("assignment failed with %v", err)
	// }

	// id, err = G.AssignVlan("Fab2", "R5.1", "", Scope_Router, 8)
	// if err != nil {
	// 	log.Errorf("assignment failed with %v", err)
	// }

	////////////////////////////////////////////

	// _ = id

	id, err := G.AssignVlanBD("MyBD01", 5)
	if err != nil {
		log.Errorf("assigning vid to bd failed with %v", err)
	}

	id, err = G.AssignVlanBD("MyBD01", 5)
	if err != nil {
		log.Infof("EXPECTED: %v", err)
	} else {
		log.Errorf("should raise error")
	}

	_, err = G.AssignVlan("", "", "", enum.Scope_Global, 7)
	if err != nil {
		log.Error(err)
	}

	id, err = G.AssignVlanBD("MyBD01", 8)
	if err != nil {
		log.Errorf("assigning vid to bd failed with %v", err)
	}

	err = G.AssignRouterToBD("myRouter", "MyBD01")
	if err != nil {
		log.Errorf("assigning router to bd failed with %v", err)
	}

	err = G.AssignInterfaceToBD("myOtherRouter", "eth0", "MyBD01")
	if err != nil {
		log.Errorf("assigning interface to bd failed with %v", err)
	}

	err = G.AssignInterfaceToBD("myNextOtherRouter", "eth0", "MyBD02")
	if err != nil {
		log.Errorf("assigning interface to bd failed with %v", err)
	}

	id, err = G.AssignVlanBD("MyBD02", 8)
	if err != nil {
		log.Errorf("assigning vid to bd failed with %v", err)
	}

	err = G.AssignRouterToBD("myRouter", "MyBD02")
	if err != nil {
		log.Infof("EXPECTED: %v", err)
	} else {
		log.Errorf("should raise error")
	}

	// DISCUSS: Should it be allowed to have a router assigned to a BD and one to many of its interfaces can be assigned to another BD
	// right now that is possible, see:
	err = G.AssignRouterToBD("myOtherRouter", "MyBD03")
	if err != nil {
		log.Errorf("assigning router to bd failed with %v", err)
	}

	for x := entities.VlanStart; x <= entities.VlanEnd; x++ {
		_, err := G.GetFreeBDVlan("MyBD02")
		if err != nil {
			log.Errorf("creating the %dth vlan caused %v", x, err)
		}
	}

	_ = id

	// print a treeview of the collected data
	fmt.Println("\n--------------------\nTopology:\n--------------------")
	fmt.Println(G.String(0))

	fmt.Println("\n--------------------\nVlanIndex:\n--------------------")
	// print a list of vlans
	vidx := G.GetVlanIndex()
	fmt.Println(vidx.String(0))

	fmt.Println("\n--------------------\nBridgeDomainIndex:\n--------------------")
	fmt.Println(G.GetBridgeDomainIndex().String(0))

}
