package connect

import (
	"strconv"
)

const DefaultTimeoutInSeconds = 10

const DefualtDisplayBufferFlushIntervalInMilliseconds = 500

type ScriptFile struct {
	Engine                                   string
	DaemonScript                             string
	ConnectScript                            string
	TimeoutInSeconds                         int
	DisplayBufferFlushIntervalInMilliseconds int
	ControlCommand                           []string
	ReservedCommands                         []string
	OnConnectStart                           string
	OnConnectClose                           string
	OnConnectUserInput                       string
	OnConnectUserCommand                     map[string]string
	OnConnectUserSubNegotiation              map[string]string
	OnConnectMainServerCommand               map[string]string
	OnConnectMainServerSubNegotiation        map[string]string
	OnConnectLocalServerCommand              map[string]string
	OnConnectLocalServerSubNegotiation       map[string]string
	OnDaemonStart                            string
	OnDaemonClose                            string
	OnDaemonConnectStarted                   string
	OnDaemonConnectClosed                    string
	OnDaemonCommand                          string
	OnDaemonInitConnect                      string
}

type ScriptConfig struct {
	Engine                                   string
	DaemonScript                             string
	ConnectScript                            string
	DisplayBufferFlushIntervalInMilliseconds int
	ControlCommand                           map[byte]bool
	ReservedCommands                         map[byte]bool
	OnConnectStart                           string
	OnConnectClose                           string
	OnConnectUserInput                       string
	OnConnectUserCommand                     map[int]string
	OnConnectUserSubNegotiation              map[int]string
	OnConnectMainServerCommand               map[int]string
	OnConnectMainServerSubNegotiation        map[int]string
	OnConnectLocalServerCommand              map[int]string
	OnConnectLocalServerSubNegotiation       map[int]string
	OnDaemonStart                            string
	OnDaemonClose                            string
	OnDaemonConnectStarted                   string
	OnDaemonConnectClosed                    string
	OnDaemonCommand                          string
	OnDaemonInitConnect                      string
}

func NewScriptConfig() *ScriptConfig {
	return &ScriptConfig{
		ControlCommand:                     map[byte]bool{},
		ReservedCommands:                   map[byte]bool{},
		OnConnectUserCommand:               map[int]string{},
		OnConnectUserSubNegotiation:        map[int]string{},
		OnConnectMainServerCommand:         map[int]string{},
		OnConnectMainServerSubNegotiation:  map[int]string{},
		OnConnectLocalServerCommand:        map[int]string{},
		OnConnectLocalServerSubNegotiation: map[int]string{},
	}
}
func MustConvert(c *ScriptFile) *ScriptConfig {
	config := &ScriptConfig{
		Engine:                                   c.Engine,
		DaemonScript:                             c.DaemonScript,
		ConnectScript:                            c.ConnectScript,
		DisplayBufferFlushIntervalInMilliseconds: c.DisplayBufferFlushIntervalInMilliseconds,
		ControlCommand:                           map[byte]bool{},
		ReservedCommands:                         map[byte]bool{},
		OnConnectStart:                           c.OnConnectStart,
		OnConnectClose:                           c.OnConnectClose,
		OnConnectUserInput:                       c.OnConnectUserInput,
		OnConnectUserCommand:                     map[int]string{},
		OnConnectUserSubNegotiation:              map[int]string{},
		OnConnectMainServerCommand:               map[int]string{},
		OnConnectMainServerSubNegotiation:        map[int]string{},
		OnConnectLocalServerCommand:              map[int]string{},
		OnConnectLocalServerSubNegotiation:       map[int]string{},
		OnDaemonStart:                            c.OnDaemonStart,
		OnDaemonClose:                            c.OnDaemonClose,
		OnDaemonConnectStarted:                   c.OnDaemonConnectStarted,
		OnDaemonConnectClosed:                    c.OnDaemonConnectClosed,
		OnDaemonCommand:                          c.OnDaemonCommand,
		OnDaemonInitConnect:                      c.OnDaemonInitConnect,
	}
	for _, v := range c.ControlCommand {
		i, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		config.ControlCommand[byte(i)] = true
	}
	for _, v := range c.ReservedCommands {
		i, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		config.ReservedCommands[byte(i)] = true
	}
	for k, v := range c.OnConnectUserCommand {
		i, err := strconv.Atoi(k)
		if err != nil {
			panic(err)
		}
		config.OnConnectUserCommand[i] = v
	}
	for k, v := range c.OnConnectMainServerCommand {
		i, err := strconv.Atoi(k)
		if err != nil {
			panic(err)
		}
		config.OnConnectMainServerCommand[i] = v
	}
	for k, v := range c.OnConnectLocalServerCommand {
		i, err := strconv.Atoi(k)
		if err != nil {
			panic(err)
		}
		config.OnConnectLocalServerCommand[i] = v
	}
	for k, v := range c.OnConnectUserSubNegotiation {
		i, err := strconv.Atoi(k)
		if err != nil {
			panic(err)
		}
		config.OnConnectUserSubNegotiation[i] = v
	}
	for k, v := range c.OnConnectMainServerSubNegotiation {
		i, err := strconv.Atoi(k)
		if err != nil {
			panic(err)
		}
		config.OnConnectMainServerSubNegotiation[i] = v
	}
	for k, v := range c.OnConnectLocalServerSubNegotiation {
		i, err := strconv.Atoi(k)
		if err != nil {
			panic(err)
		}
		config.OnConnectLocalServerSubNegotiation[i] = v
	}
	return config
}
