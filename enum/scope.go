package enum

type Scope string

const (
	Scope_Global    Scope = "global"
	Scope_Fabric    Scope = "fabric"
	Scope_Router    Scope = "router"
	Scope_Interface Scope = "interface"

	Scope_BridgeDomain Scope = "bridge_domain"
)
