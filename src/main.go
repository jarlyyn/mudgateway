package main

import (
	_ "mudgateway/modules"
	"mudgateway/modules/app"
	_ "mudgateway/modules/drivers"
	"mudgateway/modules/overseers"

	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
)

func loadConfigs() {
	//Uncomment next line to print config loading log .
	//config.Debug = true
	config.Lock.RLock()
	config.LoadAll()
	config.Lock.RUnlock()
}

func initModules() {
	util.InitModulesOrderByName()
	//Put Your own init code here.
}

//Main app run func.
var run = func() {
	//Put your run code here
	util.WaitingQuit()
	//Delay util.QuitDelayDuration for modules quit.
	util.DelayAndQuit()

}

//Init init app
func Init() {
	defer util.RecoverAndExit()
	util.ApplieationLock.Lock()
	defer util.ApplieationLock.Unlock()
	util.UpdatePaths()
	util.MustChRoot()
	loadConfigs()
	overseers.MustInitOverseers()
	initModules()
	app.Development.NotTestingOrPanic()
	//Auto created appdata folder if not exists
	util.RegisterDataFolder("scripts")
	util.MustLoadRegisteredFolders()
	app.Development.InitializeAndPanicIfNeeded()
	overseers.MustTrainWorkers()
}

func main() {
	// Set app root path.
	//Default rootpah is "src/../"
	//You can set os env  "HerbRoot" to overwrite this setting while starting app.
	// util.RootPath = ""
	Init()
	run()
}
