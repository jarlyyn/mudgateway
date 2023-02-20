package gatewayconfig

type Config struct {
	Script           string
	TimeoutInSeconds int
	Servers          []*Server
}

type Server struct {
	Port   string
	Target string
	Tags   []string
}
