package connect

import (
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
)

type DaemonAPI struct {
	Manager *Manager
}

func (a *DaemonAPI) SetVariable(key, value string) {
	a.Manager.SetVariable(key, value)
}
func (a *DaemonAPI) GetVariable(key string) string {
	return a.Manager.GetVariable(key)
}

func (a *DaemonAPI) ListConnects() []string {
	list := []string{}
	a.Manager.Connects.Range(func(key, value interface{}) bool {
		conn := value.(*Connect)
		if conn != nil {
			list = append(list, conn.ID)
		}
		return true
	})
	return list
}
func (a *DaemonAPI) getConnect(id string) *Connect {
	value, ok := a.Manager.Connects.Load(id)
	if !ok {
		return nil
	}
	return value.(*Connect)
}
func (a *DaemonAPI) GetConnectSession(id, key string) string {
	conn := a.getConnect(id)
	if conn == nil {
		return ""
	}
	value, ok := conn.Session.Load(key)
	if !ok {
		return ""
	}
	return value.(string)
}
func (a *DaemonAPI) SetConnectSession(id, key, value string) {
	conn := a.getConnect(id)
	if conn != nil {
		conn.Session.Store(key, value)
	}
}
func (a *DaemonAPI) GetConnectTags(id string) []string {
	conn := a.getConnect(id)
	if conn == nil {
		return nil
	}
	return conn.GetTags()
}

func (a *DaemonAPI) SetConnectTag(id string, key string, value bool) {
	conn := a.getConnect(id)
	if conn != nil {
		if value {
			conn.Tags.Store(key, true)
		} else {
			conn.Tags.Delete(key)
		}
	}
}

func (a *DaemonAPI) CloseConnect(id string) bool {
	conn := a.getConnect(id)
	if conn == nil {
		return false
	}
	ok, err := conn.Close()
	if err != nil {
		a.Manager.OnError(err)
		return false
	}
	return ok
}
func (a *DaemonAPI) GetConnectInfo(id string) *Info {
	conn := a.getConnect(id)
	if conn == nil {
		return nil
	} else {
		return conn.Info()
	}
}

type ConnectAPI struct {
	Connect *Connect
}

func (a *ConnectAPI) Close() bool {
	ok, err := a.Connect.Close()
	if err != nil {
		a.Connect.OnError(err)
		return false
	}
	return ok
}

func (a *ConnectAPI) SetSession(key, value string) {
	a.Connect.Session.Store(key, value)
}
func (a *ConnectAPI) GetSession(key string) string {
	v, ok := a.Connect.Session.Load(key)
	if !ok {
		return ""
	}
	return v.(string)
}

func (a *ConnectAPI) GetTags() []string {
	return a.Connect.GetTags()
}

func (a *ConnectAPI) SetTag(key string, value bool) {
	if value {
		a.Connect.Tags.Store(key, true)
	} else {
		a.Connect.Tags.Delete(key)
	}
}
func (a *ConnectAPI) HasTag(key string) bool {
	value, ok := a.Connect.Tags.Load(key)
	if !ok {
		return false
	}
	v, ok := value.(bool)
	if !ok {
		return false
	}
	return v
}

func (a *ConnectAPI) ID() string {
	return a.Connect.ID
}
func (a *ConnectAPI) Main() string {
	return a.Connect.MainServer
}
func (a *ConnectAPI) Port() string {
	return a.Connect.Port
}
func (a *ConnectAPI) GetIP() string {
	return a.Connect.GetIP()
}
func (a *ConnectAPI) SetIP(v string) {
	a.Connect.SetIP(v)
}
func (a *ConnectAPI) GetUser() string {
	return a.Connect.GetUser()
}
func (a *ConnectAPI) SetUser(v string) {
	a.Connect.SetUser(v)
}

func (a *ConnectAPI) WriteDisplayBufferRaw(data []byte) bool {
	b := NewRawBlock(data)
	a.Connect.WriteDisplayBuffer(b.Bytes())
	return true
}

func (a *ConnectAPI) WriteDisplayBufferEscaped(data []byte) bool {
	b := NewTextBlock(data)
	a.Connect.WriteDisplayBuffer(b.Bytes())
	return true
}
func (a *ConnectAPI) FlushDisplayBuffer() {
	a.Connect.flushDisplayBuffer()
}

func (a *ConnectAPI) WriteSendBufferRaw(data []byte) bool {
	b := NewRawBlock(data)
	a.Connect.WriteSendBuffer(b.Bytes())
	return true
}

func (a *ConnectAPI) WriteSendBufferEscaped(data []byte) bool {
	b := NewTextBlock(data)
	a.Connect.WriteSendBuffer(b.Bytes())
	return true
}
func (a *ConnectAPI) FlushSendBuffer() {
	a.Connect.flushSendBuffer()
}

type OutputAPI struct {
	Connect *Connect
}

func (a *OutputAPI) WriteRaw(data []byte) bool {
	b := NewRawBlock(data)
	err := b.WriteTo(a.Connect.Input)
	if err != nil {
		a.Connect.OnError(err)
		return false
	}
	return true
}

func (a *OutputAPI) WriteEscaped(data []byte) bool {
	b := NewTextBlock(data)
	err := b.WriteTo(a.Connect.Input)
	if err != nil {
		a.Connect.OnError(err)
		return false
	}
	return true
}

type BinaryAPI struct {
	Manager *Manager
}

func (a *BinaryAPI) Encode(data string) []byte {
	return []byte(data)
}

func (a *BinaryAPI) EncodeGBK(data string) []byte {
	encoder := simplifiedchinese.GBK.NewEncoder()
	d, err := encoder.Bytes([]byte(data))
	if err != nil {
		a.Manager.OnError(err)
		return nil
	}
	return d
}
func (a *BinaryAPI) EncodeBIG5(data string) []byte {
	encoder := traditionalchinese.Big5.NewEncoder()
	d, err := encoder.Bytes([]byte(data))
	if err != nil {
		a.Manager.OnError(err)
		return nil
	}
	return d
}
func (a *BinaryAPI) Decode(data []byte) *string {
	result := string(data)
	return &result
}

func (a *BinaryAPI) DecodeGBK(data []byte) *string {
	decoder := simplifiedchinese.GBK.NewDecoder()
	d, err := decoder.Bytes(data)
	if err != nil {
		a.Manager.OnError(err)
		return nil
	}
	result := string(d)
	return &result
}

func (a *BinaryAPI) DecodeBIG5(data []byte) *string {
	decoder := traditionalchinese.Big5.NewDecoder()
	d, err := decoder.Bytes(data)
	if err != nil {
		a.Manager.OnError(err)
		return nil
	}
	result := string(d)
	return &result
}

func (a *BinaryAPI) Cmd(cmd byte) []byte {
	return NewCommandNoOptBlock(cmd).Bytes()
}

func (a *BinaryAPI) CmdWithOpt(cmd byte, opt byte) []byte {
	return NewCommandBlock(cmd, opt).Bytes()
}

func (a *BinaryAPI) SubNegotiation(opt byte, data []byte) []byte {
	return NewSubNegotiationBlock(opt, data).Bytes()
}
