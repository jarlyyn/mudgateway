package connect

import (
	"mudgateway/modules/app/gatewayconfig"
	"strconv"
)

const DefaultTimeoutInSeconds = 10

const DefualtDisplayBufferFlustIntervalInMilliseconds = 500

type Config struct {
	Engine                                   string
	DaemonScript                             string
	ConnectScript                            string
	Servers                                  []*gatewayconfig.Server
	TimeoutInSeconds                         int
	DisplayBufferFlustIntervalInMilliseconds int
	ControlCommand                           map[byte]bool
	ReservedCommands                         map[byte]bool
	OnConnectStart                           string
	OnConnectClose                           string
	OnConnectUserInput                       string
	OnConnectUserCommand                     map[int]string
	OnConnectServerCommand                   map[int]string
	OnConnectUserSubNegotiation              map[int]string
	OnConnectServerSubNegotiation            map[int]string
	OnDaemonStart                            string
	OnDaemonClose                            string
	OnTick                                   string
	OnDaemonConnectStarted                   string
	OnDaemonConnectClosed                    string
	OnDaemonCommand                          string
	OnDaemonInitConnect                      string
}

func MustConvert(c *gatewayconfig.Config) *Config {
	config := &Config{
		Engine:                                   c.Engine,
		DaemonScript:                             c.DaemonScript,
		ConnectScript:                            c.ConnectScript,
		Servers:                                  c.Servers,
		DisplayBufferFlustIntervalInMilliseconds: c.DisplayBufferFlustIntervalInMilliseconds,
		TimeoutInSeconds:                         c.TimeoutInSeconds,
		ControlCommand:                           map[byte]bool{},
		ReservedCommands:                         map[byte]bool{},
		OnConnectStart:                           c.OnConnectStart,
		OnConnectClose:                           c.OnConnectClose,
		OnConnectUserInput:                       c.OnConnectUserInput,
		OnConnectUserCommand:                     map[int]string{},
		OnConnectServerCommand:                   map[int]string{},
		OnConnectUserSubNegotiation:              map[int]string{},
		OnConnectServerSubNegotiation:            map[int]string{},
		OnDaemonStart:                            c.OnDaemonStart,
		OnDaemonClose:                            c.OnDaemonClose,
		OnTick:                                   c.OnTick,
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
	for k, v := range c.OnConnectServerCommand {
		i, err := strconv.Atoi(k)
		if err != nil {
			panic(err)
		}
		config.OnConnectServerCommand[i] = v
	}
	for k, v := range c.OnConnectUserSubNegotiation {
		i, err := strconv.Atoi(k)
		if err != nil {
			panic(err)
		}
		config.OnConnectUserSubNegotiation[i] = v
	}
	for k, v := range c.OnConnectServerSubNegotiation {
		i, err := strconv.Atoi(k)
		if err != nil {
			panic(err)
		}
		config.OnConnectServerSubNegotiation[i] = v
	}
	return config
}
