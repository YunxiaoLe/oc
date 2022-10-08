package main

import (
	"flag"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/client/dbclt"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/client/emailclt"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/client/redisclt"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/conf"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/router"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/util"
	"github.com/gin-gonic/gin"
)

var port *string

func SettingUpEnvironment() {
	c := conf.ReadSettingsFromFile("Config.json")
	initArgs(c.Version)
	util.InitUtilSettings(c.Key)
	dbclt.InitDb(c.DbSettings)
	redisclt.InitRedis(c.RedisSettings)
	//结巴分词
	//business.InitSegmentation()
	emailclt.InitEmailCtl(c.EmailSenderSettings)
}
func initArgs(version string) {
	port = flag.String("port", "8080", "Listen port")
	//enableEs = *flag.Bool("enables", false, "Decide whether to enable esclt")
	flag.Parse()
	//Check whether json file version is match to server version,avoid using
	//Develop Server json file on deployed server
	if *port != version {
		panic(any("Input json doesn't match server!! Pay attention to its version"))
	}
}

func main() {
	SettingUpEnvironment()
	r := gin.Default()
	router.UseMyRouter(r)
	des := ":" + *port
	_ = r.Run(des)
}
