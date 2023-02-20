package connect

import (
	"fmt"
	"io/ioutil"
	"mudgateway/modules/app/gatewayconfig"
	"net"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/herb-go/util"
)

type Manager struct {
	Connects     sync.Map
	Config       *gatewayconfig.Config
	ScriptConfig *ScriptConfig
	Daemon       Daemon
	listeners    sync.Map
	Variables    sync.Map
}

func (m *Manager) GetVariable(key string) string {
	v, ok := m.Variables.Load(key)
	if !ok {
		return ""
	}
	return v.(string)

}
func (m *Manager) SetVariable(key, value string) {
	m.Variables.Store(key, value)
}
func (m *Manager) Start(c *gatewayconfig.Config) {
	m.Config = c
	script := c.Script

	if script != "" {
		cf := &ScriptFile{}
		data, err := ioutil.ReadFile(util.AppData(c.Script, "script.toml"))
		if err != nil {
			panic(err)
		}
		err = toml.Unmarshal(data, cf)
		if err != nil {
			panic(err)
		}
		m.ScriptConfig = MustConvert(cf)
		switch m.ScriptConfig.Engine {
		case "":
			DefaultManager.Daemon = NopDaemon{}
		default:
			f := Factories[m.ScriptConfig.Engine]
			if f == nil {
				panic(fmt.Errorf("未知的 脚本引擎:[" + m.ScriptConfig.Engine + "]"))
			}
			d, err := f.Create(m)
			if err != nil {
				panic(err)
			}
			m.Daemon = d
		}
	} else {
		m.ScriptConfig = NewScriptConfig()
	}

	for _, v := range c.Servers {
		server, err := net.Listen("tcp", v.Port)
		if err != nil {
			panic(err)
		}
		m.listeners.Store(v.Port, server)
		port := v.Port
		main := v.Target
		tags := v.Tags
		go func() {
			m.NewIncome(server, port, main, tags)
		}()
	}
	return
}
func (m *Manager) ConnectClosed(id string) {
	m.Connects.Delete(id)
}
func (m *Manager) OnError(err error) {
	util.LogError(err)
}
func (m *Manager) NewIncome(listener net.Listener, port, main string, tags []string) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		if conn != nil {
			m.OnConnect(conn, port, main, tags)
		}
	}
}
func (m *Manager) Close() error {
	return nil
}
func (m *Manager) OnConnect(rawconn net.Conn, port, main string, tags []string) {
	conn := New(rawconn)
	conn.Manager = m
	conn.MainServer = main
	conn.Port = port
	for k := range tags {
		conn.Tags.Store(tags[k], true)
	}
	m.Daemon.OnDaemonInitConnect(conn)
	ok, err := conn.Start()
	if err != nil {
		go m.OnError(err)
		return
	}
	if !ok {
		return
	}
	m.Connects.Store(conn.ID, conn)
}

func NewManager() *Manager {
	return &Manager{}
}

var DefaultManager *Manager
