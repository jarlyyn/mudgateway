package connect

type VM interface {
	OnConnectStart() bool
	OnConnectClose() bool
	OnConnectUserCommand(*Block) bool
	OnConnectUserInput([]byte) bool
	OnConnectServerCommand(*Block) bool
	OnConnectServerSubNegotiation(*Block) bool
	OnConnectUserSubNegotiation(*Block) bool
}

type NopVM struct{}

func (NopVM) OnConnectStart() bool                      { return false }
func (NopVM) OnConnectClose() bool                      { return false }
func (NopVM) OnConnectUserCommand(*Block) bool          { return false }
func (NopVM) OnConnectUserInput([]byte) bool            { return false }
func (NopVM) OnConnectServerCommand(*Block) bool        { return false }
func (NopVM) OnConnectServerSubNegotiation(*Block) bool { return false }
func (NopVM) OnConnectUserSubNegotiation(*Block) bool   { return false }

var _ VM = &NopVM{}
