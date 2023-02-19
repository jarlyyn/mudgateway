package app

import (
	"mudgateway/modules/app/gatewayconfig"
	"sync/atomic"

	"github.com/herb-go/herbconfig/source"
	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
	"github.com/herb-go/util/config/tomlconfig"
)

//Gateway config instance of gateway.
var Gateway = &gatewayconfig.Config{}

var syncGateway atomic.Value

//StoreGateway atomically store gateway config
func (a *appSync) StoreGateway(c *gatewayconfig.Config) {
	syncGateway.Store(c)
}

//LoadGateway atomically load gateway config
func (a *appSync) LoadGateway() *gatewayconfig.Config {
	v := syncGateway.Load()
	if v == nil {
		return nil
	}
	return v.(*gatewayconfig.Config)
}

func init() {
	//Register loader which will be execute when Config.LoadAll func be called.
	//You can put your init code after load.
	//You must panic if any error rasied when init.
	config.RegisterLoader(util.ConfigFile("/gateway.toml"), func(configpath source.Source) {
		util.Must(tomlconfig.Load(configpath, Gateway))
		Sync.StoreGateway(Gateway)
	})
}
