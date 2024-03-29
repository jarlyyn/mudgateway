package jsvm

import (
	"fmt"
	"mudgateway/modules/connect"
	"sync"

	"github.com/dop251/goja"
	"github.com/herb-go/util"
)

type JSVM struct {
	Lock    sync.Mutex
	runtime *goja.Runtime
	output  goja.Value
	send    goja.Value
	connect.NopVM
	Connect *connect.Connect
}

func (vm *JSVM) Call(source string, args ...interface{}) goja.Value {
	s, err := vm.runtime.RunString(source)
	if err != nil {
		vm.Connect.OnError(err)
		return nil
	}
	fn, ok := goja.AssertFunction(s)
	if !ok {
		vm.Connect.OnError(fmt.Errorf("js function %s not found", source))
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
		vm.Connect.OnError(scripterr)
		return nil
	}
	if err != nil {
		vm.Connect.OnError(err)
		return nil
	}
	return result
}
func (v *JSVM) OnConnectStart() bool {
	if v.Connect.Manager.ScriptConfig.OnConnectStart == "" {
		return false
	}
	v.Lock.Lock()
	defer v.Lock.Unlock()
	result := v.Call(v.Connect.Manager.ScriptConfig.OnConnectStart)
	if result == nil {
		return false
	}
	return result.ToBoolean()
}
func (v *JSVM) OnConnectClose() bool {
	if v.Connect.Manager.ScriptConfig.OnConnectClose == "" {
		return false
	}
	v.Lock.Lock()
	defer v.Lock.Unlock()

	result := v.Call(v.Connect.Manager.ScriptConfig.OnConnectClose)
	if result == nil {
		return false
	}
	return result.ToBoolean()
}
func (v *JSVM) OnConnectUserCommand(b *connect.Block) bool {
	if v.Connect.Manager.ScriptConfig.OnConnectUserCommand[int(b.Command)] == "" {
		return false
	}
	v.Lock.Lock()
	defer v.Lock.Unlock()
	result := v.Call(v.Connect.Manager.ScriptConfig.OnConnectUserCommand[int(b.Command)], v.send, b.Command, b.Opt, v.runtime.NewArrayBuffer(b.Data))
	if result == nil {
		return false
	}
	return result.ToBoolean()
}
func (v *JSVM) OnConnectUserInput(data []byte) bool {
	if v.Connect.Manager.ScriptConfig.OnConnectUserInput == "" {
		return false
	}
	v.Lock.Lock()
	defer v.Lock.Unlock()
	result := v.Call(v.Connect.Manager.ScriptConfig.OnConnectUserInput, v.send, v.runtime.NewArrayBuffer(data))
	if result == nil {
		return false
	}
	return result.ToBoolean()
}
func (v *JSVM) OnConnectServerCommand(b *connect.Block) bool {
	if v.Connect.Manager.ScriptConfig.OnConnectServerCommand[int(b.Command)] == "" {
		return false
	}
	v.Lock.Lock()
	defer v.Lock.Unlock()
	result := v.Call(v.Connect.Manager.ScriptConfig.OnConnectServerCommand[int(b.Command)], v.output, b.Command, b.Opt, v.runtime.NewArrayBuffer(b.Data), 0)
	if result == nil {
		return false
	}
	return result.ToBoolean()
}
func (v *JSVM) OnConnectServerSubNegotiation(b *connect.Block) bool {
	if v.Connect.Manager.ScriptConfig.OnConnectServerSubNegotiation[int(b.Opt)] == "" {
		return false
	}
	v.Lock.Lock()
	defer v.Lock.Unlock()
	result := v.Call(v.Connect.Manager.ScriptConfig.OnConnectServerSubNegotiation[int(b.Opt)], v.output, b.Opt, v.runtime.NewArrayBuffer(b.Data), 0)
	if result == nil {
		return false
	}
	return result.ToBoolean()
}
func (v *JSVM) OnConnectUserSubNegotiation(b *connect.Block) bool {
	if v.Connect.Manager.ScriptConfig.OnConnectUserSubNegotiation[int(b.Opt)] == "" {
		return false
	}
	v.Lock.Lock()
	defer v.Lock.Unlock()
	result := v.Call(v.Connect.Manager.ScriptConfig.OnConnectUserSubNegotiation[int(b.Opt)], v.send, b.Opt, v.runtime.NewArrayBuffer(b.Data))
	if result == nil {
		return false
	}
	return result.ToBoolean()
}

var _ connect.VM = &JSVM{}
