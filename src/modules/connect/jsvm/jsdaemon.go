package jsvm

import (
	"fmt"
	"io/ioutil"
	"mudgateway/modules/connect"
	"sync"

	"github.com/dop251/goja"
	"github.com/herb-go/util"
)

type JSDaemon struct {
	Lock          sync.Locker
	runtime       *goja.Runtime
	Manager       *connect.Manager
	ConnectScript string
}

func (vm *JSDaemon) Call(source string, args ...interface{}) goja.Value {
	s, err := vm.runtime.RunString(source)
	if err != nil {
		vm.Manager.OnError(err)
		return nil
	}
	fn, ok := goja.AssertFunction(s)
	if !ok {
		vm.Manager.OnError(fmt.Errorf("js function %s not found", source))
		return nil
	}
	jargs := []goja.Value{}
	for _, v := range args {
		jargs = append(jargs, vm.runtime.ToValue(v))
	}
	var result goja.Value
	var scripterr error
	err = util.Catch(func() {
		result, scripterr = fn(goja.Undefined(), jargs...)
	})
	if scripterr != nil {
		vm.Manager.OnError(scripterr)
		return nil
	}
	if err != nil {
		vm.Manager.OnError(err)
		return nil
	}
	return result

}
func (d *JSDaemon) OnDaemonStart() {
	if d.Manager.ScriptConfig.OnDaemonStart == "" {
		return
	}
	d.Lock.Lock()
	defer d.Lock.Unlock()
	d.Call(d.Manager.ScriptConfig.OnDaemonStart)
}
func (d *JSDaemon) OnDaemonClose() {
	if d.Manager.ScriptConfig.OnDaemonClose == "" {
		return
	}
	d.Lock.Lock()
	defer d.Lock.Unlock()
	d.Call(d.Manager.ScriptConfig.OnDaemonClose)
}

func (d *JSDaemon) OnDaemonCommand(cmd string) {
	if d.Manager.ScriptConfig.OnDaemonCommand == "" {
		return
	}
	d.Lock.Lock()
	defer d.Lock.Unlock()
	d.Call(d.Manager.ScriptConfig.OnDaemonCommand, cmd)
}
func (d *JSDaemon) OnDaemonInitConnect(c *connect.Connect) {
	runtime := goja.New()
	vm := &JSVM{
		Connect: c,
		runtime: runtime,
		output: runtime.ToValue(&OutputAPI{
			&connect.OutputAPI{
				Connect: c,
			},
		}),
		send: runtime.ToValue(&SendAPI{
			&connect.SendAPI{
				Connect: c,
			},
		}),
	}
	vm.runtime.Set("Binary", &BinaryAPI{
		api: &connect.BinaryAPI{
			Manager: d.Manager,
		},
	})
	vm.runtime.Set("Connect", &ConnectAPI{
		api: &connect.ConnectAPI{
			Connect: c,
		},
	})
	if d.ConnectScript != "" {
		_, err := vm.runtime.RunScript(d.Manager.ScriptConfig.ConnectScript, d.ConnectScript)
		if err != nil {
			d.Manager.OnError(err)
			return
		}
	}
	c.VM = vm
}

var _ connect.Daemon = &JSDaemon{}

type Factory struct {
}

func (f Factory) Create(m *connect.Manager) (connect.Daemon, error) {
	d := &JSDaemon{
		Manager: m,
		runtime: goja.New(),
	}
	d.runtime.Set("Daemon", &DaemonAPI{
		api: &connect.DaemonAPI{
			Manager: m,
		},
	})
	d.runtime.Set("Binary", &BinaryAPI{
		api: &connect.BinaryAPI{
			Manager: m,
		},
	})
	script := m.ScriptConfig.DaemonScript
	if script != "" {
		data, err := ioutil.ReadFile(util.AppData(m.Config.Script, script))
		if err != nil {
			return nil, err
		}
		_, err = d.runtime.RunScript(script, string(data))
		if err != nil {
			return nil, err
		}
	}
	connscript := m.ScriptConfig.ConnectScript
	if connscript != "" {
		connscriptdata, err := ioutil.ReadFile(util.AppData(m.Config.Script, connscript))
		if err != nil {
			return nil, err
		}
		d.ConnectScript = string(connscriptdata)
	}
	return d, nil
}

func init() {
	connect.Factories["js"] = Factory{}
}
