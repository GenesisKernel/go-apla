package modes

import (
	"github.com/GenesisKernel/go-genesis/packages/api"
	"github.com/GenesisKernel/go-genesis/packages/conf"
	"github.com/GenesisKernel/go-genesis/packages/config/syspar"
	"github.com/GenesisKernel/go-genesis/packages/daemons"
	"github.com/GenesisKernel/go-genesis/packages/model"
	"github.com/GenesisKernel/go-genesis/packages/smart"
	"github.com/GenesisKernel/go-genesis/packages/utils"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func InitVDEMode(config *conf.VDEConfig) *VDE {
	mode := &VDE{
		VDEConfig: config,
		api:       api.CreateDefaultRouter(),
		daemonList: []string{
			"Notificator",
			"Scheduler",
		},
	}

	return mode
}

type VDE struct {
	*conf.VDEConfig
	api        *httprouter.Router
	daemonList []string
}

func (mode *VDE) Start(exitFunc func(int), gormInit func(conf.DBConfig), listenerFunc func(string, *httprouter.Router)) {
	gormInit(mode.DB)

	listenerFunc(mode.VDEConfig.HTTP.Str(), mode.api)

	if model.DBConn != nil {
		// The installation process is already finished (where user has specified DB and where wallet has been restarted)
		err := syspar.SysUpdate(nil)
		if err != nil {
			log.Errorf("can't read system parameters: %s", utils.ErrInfo(err))
			exitFunc(1)
		}

		log.Info("load contracts")
		if err := smart.LoadContracts(nil); err != nil {
			log.Errorf("Load Contracts error: %s", err)
			exitFunc(1)
		}

		log.Info("start daemons")
		daemons.StartDaemons(mode.daemonList)

		log.Info("Daemons started")
	}
}

func (mode *VDE) DaemonList() []string {
	return mode.daemonList
}

func (mode *VDE) Stop() {
	log.Infoln("VDE mode stopped")
}
