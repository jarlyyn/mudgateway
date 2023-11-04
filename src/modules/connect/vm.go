package connect

type VM interface {
	//连接成功事件，返回true终止连接
	OnConnectStart() bool
	//连接关闭事件
	OnConnectClosed()
	OnConnectLocalServerConnected()
	OnConnectLocalServerCloseded()
	OnConnectUserCommand(*Block) bool
	OnConnectUserInput([]byte) bool
	OnConnectMainServerCommand(*Block) bool
	OnConnectMainServerSubNegotiation(*Block) bool
	OnConnectLocalServerCommand(*Block) bool
	OnConnectLocalServerSubNegotiation(*Block) bool
	OnConnectUserSubNegotiation(*Block) bool
}

type NopVM struct{}

func (NopVM) OnConnectStart() bool             { return false }
func (NopVM) OnConnectClosed()                 {}
func (NopVM) OnConnectUserCommand(*Block) bool { return false }
func (NopVM) OnConnectUserInput([]byte) bool   { return false }
func (NopVM) OnConnectLocalServerConnected()
func (NopVM) OnConnectLocalServerCloseded()
func (NopVM) OnConnectMainServerCommand(*Block) bool         { return false }
func (NopVM) OnConnectMainServerSubNegotiation(*Block) bool  { return false }
func (NopVM) OnConnectLocalServerCommand(*Block) bool        { return false }
func (NopVM) OnConnectLocalServerSubNegotiation(*Block) bool { return false }
func (NopVM) OnConnectUserSubNegotiation(*Block) bool        { return false }

var _ VM = &NopVM{}
