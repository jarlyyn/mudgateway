package jsvm

import (
	"mudgateway/modules/connect"

	"github.com/dop251/goja"
)

type ConnectAPI struct {
	api *connect.ConnectAPI
}

func (a *ConnectAPI) Close(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	return r.ToValue(a.api.Close())
}

func (a *ConnectAPI) SetSession(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	key := call.Argument(0).String()
	value := call.Argument(0).String()
	a.api.SetSession(key, value)
	return nil
}
func (a *ConnectAPI) GetSession(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	key := call.Argument(0).String()
	return r.ToValue(a.api.GetSession(key))
}

func (a *ConnectAPI) GetTags(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	return r.ToValue(a.api.GetTags())
}

func (a *ConnectAPI) SetTag(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	key := call.Argument(0).String()
	value := call.Argument(0).ToBoolean()
	a.api.SetTag(key, value)
	return nil
}
func (a *ConnectAPI) HasTag(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	key := call.Argument(0).String()
	return r.ToValue(a.api.HasTag(key))
}
func (a *ConnectAPI) ID(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	return r.ToValue(a.api.ID())
}
func (a *ConnectAPI) Main(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	return r.ToValue(a.api.Main())
}
func (a *ConnectAPI) Port(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	return r.ToValue(a.api.Port())
}
func (a *ConnectAPI) GetIP(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	return r.ToValue(a.api.GetIP())
}
func (a *ConnectAPI) SetIP(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	value := call.Argument(0).String()
	a.api.SetIP(value)
	return nil
}
func (a *ConnectAPI) GetUser(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	return r.ToValue(a.api.GetUser())
}
func (a *ConnectAPI) SetUser(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	value := call.Argument(0).String()
	a.api.SetUser(value)
	return nil
}

func (a *ConnectAPI) WriteDisplayBufferRaw(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	data := call.Argument(0).Export()
	if data == nil {
		return r.ToValue(false)
	}
	bs, ok := data.(goja.ArrayBuffer)
	if !ok {
		return r.ToValue(false)
	}
	return r.ToValue(a.api.WriteDisplayBufferRaw(bs.Bytes()))
}

func (a *ConnectAPI) WriteDisplayBufferEscaped(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	data := call.Argument(0).Export()
	if data == nil {
		return r.ToValue(false)
	}
	bs, ok := data.(goja.ArrayBuffer)
	if !ok {
		return r.ToValue(false)
	}
	return r.ToValue(a.api.WriteDisplayBufferEscaped(bs.Bytes()))
}
func (a *ConnectAPI) FlushDisplayBuffer(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	a.api.FlushDisplayBuffer()
	return nil
}

func (a *ConnectAPI) WriteSendBufferRaw(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	data := call.Argument(0).Export()
	if data == nil {
		return r.ToValue(false)
	}
	bs, ok := data.(goja.ArrayBuffer)
	if !ok {
		return r.ToValue(false)
	}
	return r.ToValue(a.api.WriteSendBufferRaw(bs.Bytes()))
}

func (a *ConnectAPI) WriteSendBufferEscaped(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	data := call.Argument(0).Export()
	if data == nil {
		return r.ToValue(false)
	}
	bs, ok := data.(goja.ArrayBuffer)
	if !ok {
		return r.ToValue(false)
	}
	return r.ToValue(a.api.WriteSendBufferEscaped(bs.Bytes()))
}
func (a *ConnectAPI) FlushSendBuffer(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	a.api.FlushSendBuffer()
	return nil
}

type DaemonAPI struct {
	api *connect.DaemonAPI
}

func (a *DaemonAPI) SetVariable(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	key := call.Argument(0).String()
	value := call.Argument(1).String()
	a.api.SetVariable(key, value)
	return nil
}
func (a *DaemonAPI) GetVariable(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	key := call.Argument(0).String()
	value := a.api.GetVariable(key)
	return r.ToValue(value)
}

func (a *DaemonAPI) ListConnects(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	value := a.api.ListConnects()
	return r.ToValue(value)
}

func (a *DaemonAPI) GetConnectSession(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	id := call.Argument(0).String()
	key := call.Argument(1).String()
	value := a.api.GetConnectSession(id, key)
	return r.ToValue(value)
}
func (a *DaemonAPI) SetConnectSession(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	id := call.Argument(0).String()
	key := call.Argument(1).String()
	value := call.Argument(3).String()
	a.api.SetConnectSession(id, key, value)
	return nil
}
func (a *DaemonAPI) GetConnectTags(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	id := call.Argument(0).String()
	key := call.Argument(1).String()
	value := a.api.GetConnectSession(id, key)
	return r.ToValue(value)
}

func (a *DaemonAPI) SetConnectTag(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	id := call.Argument(0).String()
	key := call.Argument(1).String()
	value := call.Argument(2).ToBoolean()
	a.api.SetConnectTag(id, key, value)
	return nil
}

func (a *DaemonAPI) CloseConnect(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	id := call.Argument(0).String()
	value := a.api.CloseConnect(id)
	return r.ToValue(value)
}
func (a *DaemonAPI) GetConnectInfo(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	id := call.Argument(0).String()
	info := a.api.GetConnectInfo(id)
	if info == nil {
		return nil
	}
	result := r.NewObject()
	result.Set("ID", info.ID)
	result.Set("Tags", info.Tags)
	result.Set("Session", info.Session)
	result.Set("Main", info.Main)
	result.Set("Port", info.Port)
	result.Set("User", info.User)
	result.Set("IP", info.IP)
	result.Set("Timestamp", info.Timestamp)
	return result
}

type OutputAPI struct {
	api *connect.OutputAPI
}

func (a *OutputAPI) WriteRaw(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	data := call.Argument(0).Export()
	if data == nil {
		return r.ToValue(false)
	}
	bs, ok := data.(goja.ArrayBuffer)
	if !ok {
		return r.ToValue(false)
	}
	return r.ToValue(a.api.WriteRaw(bs.Bytes()))
}

func (a *OutputAPI) WriteEscaped(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	data := call.Argument(0).Export()
	if data == nil {
		return r.ToValue(false)
	}
	bs, ok := data.(goja.ArrayBuffer)
	if !ok {
		return r.ToValue(false)
	}
	return r.ToValue(a.api.WriteEscaped(bs.Bytes()))
}

type SendAPI struct {
	api *connect.SendAPI
}

func (a *SendAPI) WriteRaw(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	data := call.Argument(0).Export()
	if data == nil {
		return r.ToValue(false)
	}
	bs, ok := data.(goja.ArrayBuffer)
	if !ok {
		return r.ToValue(false)
	}
	return r.ToValue(a.api.WriteRaw(bs.Bytes()))
}

func (a *SendAPI) WriteEscaped(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	data := call.Argument(0).Export()
	if data == nil {
		return r.ToValue(false)
	}
	bs, ok := data.(goja.ArrayBuffer)
	if !ok {
		return r.ToValue(false)
	}
	return r.ToValue(a.api.WriteEscaped(bs.Bytes()))
}

type BinaryAPI struct {
	api *connect.BinaryAPI
}

func (a *BinaryAPI) Encode(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	data := call.Argument(0).String()
	result := a.api.Encode(data)
	if result == nil {
		return nil
	}
	return r.ToValue(r.NewArrayBuffer(result))
}

func (a *BinaryAPI) EncodeGBK(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	data := call.Argument(0).String()
	result := a.api.EncodeGBK(data)
	if result == nil {
		return nil
	}
	return r.ToValue(r.NewArrayBuffer(result))
}
func (a *BinaryAPI) EncodeBIG5(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	data := call.Argument(0).String()
	result := a.api.EncodeBIG5(data)
	if result == nil {
		return nil
	}
	return r.ToValue(r.NewArrayBuffer(result))
}
func (a *BinaryAPI) Decode(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	data := call.Argument(0).Export()
	if data == nil {
		return nil
	}
	ab, ok := data.(goja.ArrayBuffer)
	if !ok {
		return nil
	}
	result := a.api.Decode(ab.Bytes())
	if result == nil {
		return nil
	}
	return r.ToValue(*result)

}

func (a *BinaryAPI) DecodeGBK(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	data := call.Argument(0).Export()
	if data == nil {
		return nil
	}
	ab, ok := data.(goja.ArrayBuffer)
	if !ok {
		return nil
	}
	result := a.api.DecodeGBK(ab.Bytes())
	if result == nil {
		return nil
	}
	return r.ToValue(*result)
}

func (a *BinaryAPI) DecodeBIG5(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	data := call.Argument(0).Export()
	if data == nil {
		return nil
	}
	ab, ok := data.(goja.ArrayBuffer)
	if !ok {
		return nil
	}
	result := a.api.DecodeBIG5(ab.Bytes())
	if result == nil {
		return nil
	}
	return r.ToValue(*result)

}

func (a *BinaryAPI) Cmd(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	data := call.Argument(0).ToInteger()
	result := a.api.Cmd(byte(data))
	if result == nil {
		return nil
	}
	return r.ToValue(r.NewArrayBuffer(result))
}

func (a *BinaryAPI) CmdWithOpt(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	data := call.Argument(0).ToInteger()
	opt := call.Argument(1).ToInteger()
	result := a.api.CmdWithOpt(byte(data), byte(opt))
	if result == nil {
		return nil
	}
	return r.ToValue(r.NewArrayBuffer(result))
}

func (a *BinaryAPI) SubNegotiation(call goja.FunctionCall, r *goja.Runtime) goja.Value {
	data := call.Argument(0).ToInteger()
	ab := call.Argument(1).Export()
	if ab == nil {
		return nil
	}
	bs := ab.(goja.ArrayBuffer)
	result := a.api.SubNegotiation(byte(data), bs.Bytes())
	if result == nil {
		return nil
	}
	return r.ToValue(r.NewArrayBuffer(result))
}
