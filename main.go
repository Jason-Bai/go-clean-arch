package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Jason-Bai/go-clean-arch/config"
	"github.com/Jason-Bai/go-clean-arch/model"
	v "github.com/Jason-Bai/go-clean-arch/pkg/version"
	"github.com/Jason-Bai/go-clean-arch/router"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	// custom config path, it's useful
	cfg = pflag.StringP("config", "c", "", "apiserver config file path.")
	// build a version command into main binary program
	version = pflag.BoolP("version", "v", false, "show version info.")
)

// 0. Keep main function more clear
func main() {
	pflag.Parse()

	// 1. Binary Program with version, to find bug quickly
	if *version {
		v := v.Get()
		marshalled, err := json.MarshalIndent(&v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshalled))
		return
	}

	// 2. Load configs
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// 3. Connect to DB
	model.DB.Init()
	defer model.DB.Close()

	gin.SetMode(viper.GetString("runmode"))

	// 4. Create a gin-gonic router
	g := router.Init()

	// 5. test avaliable in last
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}

		log.Info("The router has been deployed successfully")
	}()

	// 6. Start to listening the incoming requests with https
	cert := viper.GetString("tls.cert")
	key := viper.GetString("tls.key")

	if cert != "" && key != "" {
		go func() {
			log.Infof("Start to lisening the incoming requests on https address: %s", viper.GetString("tls.addr"))
			log.Info(http.ListenAndServeTLS(viper.GetString("tls.addr"), cert, key, g).Error())
		}()
	}

	// 7. Start to listening the incoming request with http
	log.Infof("Start to lisening the incoming requests on http address: %s", viper.GetString("addr"))
	log.Info(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

// last check avaliable
func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}

	return errors.New("Cannot connect to the router.")
}
