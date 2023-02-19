package connect

type Daemon interface {
	OnDaemonStart()
	OnDaemonClose()
	OnTick()
	OnDaemonCommand(string)
	OnDaemonInitConnect(*Connect)
}

type NopDaemon struct{}

func (NopDaemon) OnDaemonStart()         {}
func (NopDaemon) OnDaemonClose()         {}
func (NopDaemon) OnTick()                {}
func (NopDaemon) OnDaemonCommand(string) {}
func (NopDaemon) OnDaemonInitConnect(c *Connect) {
	c.VM = NopVM{}
}

var _ Daemon = &NopDaemon{}

type DaemonFactory interface {
	Create(*Manager) (Daemon, error)
}

var Factories = map[string]DaemonFactory{}
