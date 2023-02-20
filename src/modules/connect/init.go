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
		DefaultManager.Start(app.Gateway)
		for _, v := range app.Gateway.Servers {
			fmt.Println("Listnening " + v.Port + " -> " + v.Target)
		}
	})
}
