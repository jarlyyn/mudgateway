package connect

import (
	"fmt"
	"mudgateway/modules/app"

	"github.com/herb-go/util"
)

//ModuleName module name
const ModuleName = "900connect"

func init() {
	util.RegisterModule(ModuleName, func() {
		//Init registered initator which registered by RegisterInitiator
		//util.RegisterInitiator(ModuleName, "func", func(){})
		util.InitOrderByName(ModuleName)
		DefaultManager = NewManager()
		DefaultManager.Config = MustConvert(app.Gateway)
		switch DefaultManager.Config.Engine {
		case "":
			DefaultManager.Daemon = NopDaemon{}
		default:
			f := Factories[DefaultManager.Config.Engine]
			if f == nil {
				panic(fmt.Errorf("未知的 脚本引擎:[" + DefaultManager.Config.Engine + "]"))
			}
			d, err := f.Create(DefaultManager)
			if err != nil {
				panic(err)
			}
			DefaultManager.Daemon = d
		}
		util.Must(DefaultManager.Start())
		for _, v := range app.Gateway.Servers {
			fmt.Println("Listnening " + v.Port + " -> " + v.Target)
		}
	})
}
