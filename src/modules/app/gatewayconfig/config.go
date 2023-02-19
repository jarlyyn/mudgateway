package gatewayconfig

type Config struct {
	Engine                                   string
	DaemonScript                             string
	ConnectScript                            string
	Servers                                  []*Server
	TimeoutInSeconds                         int
	DisplayBufferFlustIntervalInMilliseconds int
	ControlCommand                           []string
	ReservedCommands                         []string
	OnConnectStart                           string
	OnConnectClose                           string
	OnConnectUserInput                       string
	OnConnectUserCommand                     map[string]string
	OnConnectServerCommand                   map[string]string
	OnConnectUserSubNegotiation              map[string]string
	OnConnectServerSubNegotiation            map[string]string
	OnDaemonStart                            string
	OnDaemonClose                            string
	OnTick                                   string
	OnDaemonConnectStarted                   string
	OnDaemonConnectClosed                    string
	OnDaemonCommand                          string
	OnDaemonInitConnect                      string
}

type Server struct {
	Port   string
	Target string
	Tags   []string
}
