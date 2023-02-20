package connect

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/herb-go/uniqueid"
	"github.com/herb-go/util"
)

//Connect 连接对象
type Connect struct {
	//连接级锁
	Lock                        sync.Mutex
	VM                          VM
	InputRaw                    net.Conn
	Input                       *bufio.ReadWriter
	MainConnect                 *bufio.ReadWriter
	MainRaw                     net.Conn
	MainServer                  string
	Port                        string
	InstanceConnect             *bufio.ReadWriter
	InstanceRaw                 net.Conn
	Current                     atomic.Value
	ID                          string
	User                        atomic.Value
	IP                          atomic.Value
	CreatedAt                   time.Time
	Session                     sync.Map
	Tags                        sync.Map
	Manager                     *Manager
	displayBuffer               *bytes.Buffer
	displayBufferFlushable      bool
	displayBufferTicker         *time.Ticker
	displayBufferTickerDuraiton time.Duration
	sendBuffer                  *bytes.Buffer
	sendBufferFlushable         bool
	userinputBuffer             *bytes.Buffer
}

func (c *Connect) GetTags() []string {
	result := []string{}
	c.Tags.Range(func(key, value interface{}) bool {
		v, ok := value.(bool)
		if ok && v {
			result = append(result, key.(string))
		}
		return true
	})
	return result
}
func (c *Connect) GetSession() map[string]string {
	result := map[string]string{}
	c.Session.Range(func(key, value interface{}) bool {
		result[value.(string)] = key.(string)
		return true
	})
	return result
}
func (c *Connect) GetIP() string {
	v := c.IP.Load()
	if v == nil {
		return ""
	}
	return v.(string)
}
func (c *Connect) SetIP(v string) {
	c.IP.Store(v)
}
func (c *Connect) GetUser() string {
	v := c.User.Load()
	if v == nil {
		return ""
	}
	return v.(string)
}
func (c *Connect) SetUser(v string) {
	c.User.Store(v)
}
func (c *Connect) Info() *Info {
	i := &Info{
		ID:        c.ID,
		User:      c.GetUser(),
		IP:        c.GetIP(),
		Main:      c.MainServer,
		Port:      c.Port,
		Timestamp: strconv.FormatInt(c.CreatedAt.Unix(), 10),
		Tags:      c.GetTags(),
		Session:   c.GetSession(),
	}
	return i
}
func (c *Connect) Start() (bool, error) {
	if c.VM.OnConnectStart() {
		return false, nil
	}
	timeout := c.Manager.Config.TimeoutInSeconds
	if timeout <= 0 {
		timeout = DefaultTimeoutInSeconds
	}
	main, err := net.DialTimeout("tcp", c.MainServer, time.Second*(time.Duration(timeout)))
	if err != nil {
		return false, nil
	}
	c.MainConnect = bufio.NewReadWriter(bufio.NewReader(main), bufio.NewWriter(main))
	c.Current.Store(c.MainConnect)
	t := c.Manager.ScriptConfig.DisplayBufferFlushIntervalInMilliseconds
	if t == 0 {
		t = DefualtDisplayBufferFlushIntervalInMilliseconds
	}
	if t > 0 {
		c.displayBufferTickerDuraiton = time.Duration(t) * time.Millisecond
		c.displayBufferTicker = time.NewTicker(c.displayBufferTickerDuraiton)
		go c.Tick()
	}
	go c.ListenMainServer()
	go c.ListenUserInput()
	return true, nil
}
func (c *Connect) Tick() {
	for {
		select {
		case _, isClose := <-c.displayBufferTicker.C:
			if isClose {
				return
			} else {
				go c.OnTick()
			}

		}
	}
}
func (c *Connect) Close() (bool, error) {
	if c.VM.OnConnectClose() {
		return false, nil
	}
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return true, nil
}
func (c *Connect) ServerClosed() {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	if c.InputRaw != nil {
		c.InputRaw.Close()
	}
	c.InputRaw = nil
	c.Input = nil
	if c.MainRaw != nil {
		c.MainRaw.Close()
	}
	c.MainRaw = nil
	c.MainConnect = nil
	c.Manager.ConnectClosed(c.ID)
	if c.displayBufferTicker != nil {
		c.displayBufferTicker.Stop()
		c.displayBufferTicker = nil
	}
}
func (c *Connect) OnError(err error) {
	util.LogError(err)
}
func (c *Connect) SendBlockToServer(b *Block) error {
	current := c.Current.Load()
	if current != nil {
		conn := current.(*bufio.ReadWriter)
		err := b.WriteTo(conn)
		if err != nil {
			return err
		}
		return conn.Flush()
	}
	return nil
}
func (c *Connect) SendBlockToUser(b *Block) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return c.sendBlockToUser(b)
}
func (c *Connect) sendBlockToUser(b *Block) error {
	conn := c.Input
	if conn != nil {
		err := b.WriteTo(conn)
		if err != nil {
			return err
		}
		return conn.Flush()
	}
	return nil
}
func (c *Connect) flushUserinput() {
	c.displayBufferFlushable = true
	c.flushSendBuffer()
	if c.userinputBuffer.Len() == 0 {
		return
	}
	bs := c.userinputBuffer.Bytes()
	result := c.VM.OnConnectUserInput(bs)
	if !result {
		err := c.SendBlockToServer(NewRawBlock(bs))
		if err != nil {
			c.OnError(err)
		}
	}
	c.userinputBuffer.Reset()
}
func (c *Connect) ListenUserInput() {
	var last byte
	for {
		c.Lock.Lock()
		in := c.Input
		c.Lock.Unlock()
		if in == nil {
			return
		}
		b, err := ReadBlock(in)
		c.Lock.Lock()
		c.displayBufferFlushable = false
		if err != nil {
			if errors.Is(err, net.ErrClosed) || err == io.EOF {
				c.Lock.Unlock()
				c.ServerClosed()
				return
			}
			c.OnError(err)
			c.Lock.Unlock()
			continue
		}
		switch b.Type {
		case BlockTypeCommand, BlockTypeCommandNoOpt:
			c.flushUserinput()
			last = 0
			result := c.VM.OnConnectUserCommand(b)
			if !result {
				err := c.SendBlockToServer(b)
				if err != nil {
					c.OnError(err)
					c.Lock.Unlock()
					continue
				}
			}
		case BlockTypeSubNegotiation:
			c.flushUserinput()
			last = 0
			if !c.Manager.ScriptConfig.ReservedCommands[b.Opt] && !c.Manager.ScriptConfig.ControlCommand[b.Opt] {
				result := c.VM.OnConnectUserSubNegotiation(b)
				if !result {
					err := c.SendBlockToServer(b)
					if err != nil {
						c.OnError(err)
						c.Lock.Unlock()
						continue
					}
				}
			}

		default:
			if len(b.Data) == 1 {
				switch b.Data[0] {
				case CR:
					c.userinputBuffer.WriteByte(LF)
					c.flushUserinput()
				case LF:
					if last != CR {
						c.userinputBuffer.WriteByte(LF)
						c.flushUserinput()
					}
				default:
					c.userinputBuffer.WriteByte(b.Data[0])
				}
				last = b.Data[0]

			}
		}
		c.Lock.Unlock()
	}
}
func (c *Connect) FlushDisplayBuffer() {
	if c.displayBuffer.Len() > 0 {
		c.flushDisplayBuffer()
	}
}
func (c *Connect) flushDisplayBuffer() {
	if c.displayBuffer.Len() == 0 {
		return
	}
	if c.Input != nil {
		_, err := c.Input.Write(c.displayBuffer.Bytes())
		c.displayBuffer.Reset()
		if err != nil {
			c.OnError(err)
		}
		err = c.Input.Flush()
		if err != nil {
			c.OnError(err)
		}

	}
}
func (c *Connect) WriteDisplayBuffer(data []byte) {
	_, err := c.displayBuffer.Write(data)
	if err != nil {
		c.OnError(err)
	}
	if c.displayBufferFlushable {
		c.flushDisplayBuffer()
	}
}

func (c *Connect) FlushSendBuffer() {
	if c.sendBuffer.Len() > 0 {
		c.flushSendBuffer()
	}
}
func (c *Connect) flushSendBuffer() {
	if c.sendBuffer.Len() == 0 {
		return
	}
	v := c.Current.Load()
	if v == nil {
		return
	}
	conn := v.(*bufio.ReadWriter)
	if conn != nil {
		_, err := conn.Write(c.sendBuffer.Bytes())
		c.sendBuffer.Reset()
		if err != nil {
			c.OnError(err)
		}
		err = conn.Flush()
		if err != nil {
			c.OnError(err)
		}
	}
}
func (c *Connect) WriteSendBuffer(data []byte) {
	_, err := c.sendBuffer.Write(data)
	if err != nil {
		c.OnError(err)
	}
	if c.sendBufferFlushable {
		c.flushSendBuffer()
	}
}
func (c *Connect) ListenMainServer() {
	for {
		c.Lock.Lock()
		conn := c.MainConnect
		c.Lock.Unlock()
		if conn == nil {
			return
		}
		b, err := ReadBlock(conn)
		c.Lock.Lock()
		c.displayBufferFlushable = false
		if err != nil {
			if errors.Is(err, net.ErrClosed) || err == io.EOF {
				c.Lock.Unlock()
				c.ServerClosed()
				return
			}
			c.OnError(err)
			c.Lock.Unlock()
			continue
		}
		if c.displayBufferTicker != nil {
			c.displayBufferTicker.Reset(c.displayBufferTickerDuraiton)
		}
		switch b.Type {
		case BlockTypeCommand, BlockTypeCommandNoOpt:
			result := c.VM.OnConnectServerCommand(b)
			if !result {
				err := c.sendBlockToUser(b)
				if err != nil {
					c.OnError(err)
				}
				if b.Command == IACCmdGA {
					c.displayBufferFlushable = true
					c.flushDisplayBuffer()
				}
			}
		case BlockTypeSubNegotiation:
			if !c.Manager.ScriptConfig.ReservedCommands[b.Opt] && !c.Manager.ScriptConfig.ControlCommand[b.Opt] {
				result := c.VM.OnConnectServerSubNegotiation(b)
				if !result {
					err := c.sendBlockToUser(b)
					if err != nil {
						c.OnError(err)
					}
				}
			}
		default:
			err := c.sendBlockToUser(b)
			if err != nil {
				c.OnError(err)
			}
		}
		c.Lock.Unlock()
	}
}
func (c *Connect) OnTick() {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.displayBufferFlushable = true
	c.flushDisplayBuffer()
}
func New(in net.Conn) *Connect {
	inconn := bufio.NewReadWriter(bufio.NewReader(in), bufio.NewWriter(in))
	h, _, err := net.SplitHostPort(in.RemoteAddr().String())
	if err != nil {
		panic(err)
	}
	c := &Connect{
		Input:                  inconn,
		InputRaw:               in,
		ID:                     uniqueid.MustGenerateID(),
		CreatedAt:              time.Now(),
		displayBuffer:          bytes.NewBuffer(nil),
		sendBuffer:             bytes.NewBuffer(nil),
		userinputBuffer:        bytes.NewBuffer(nil),
		sendBufferFlushable:    true,
		displayBufferFlushable: true,
	}
	c.SetIP(h)
	return c
}

type Info struct {
	ID        string
	Tags      []string
	Session   map[string]string
	Main      string
	Port      string
	User      string
	IP        string
	Timestamp string
}
