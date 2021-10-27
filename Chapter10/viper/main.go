package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"log"
	"net/http"
	"time"
)

func init() {
	rootCmd.Flags().IntP("number", "n", 42, "an integer")
	viper.BindPFlag("number", rootCmd.Flags().Lookup("number"))
}

var rootCmd = &cobra.Command{}
var runtime_viper = viper.New()

type Config struct {
	Host string
	Port uint16
	Tags map[string]string
}

func getTagHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(runtime_viper.GetString("host")))
}
func main() {
	fmt.Println(viper.GetInt("number"))
	viper.BindEnv("id")
	viper.BindEnv("port", "SERVICE_PORT")
	viper.SetDefault("id", "13")

	fmt.Println(viper.GetInt("id"))
	fmt.Println(viper.GetInt("port"))

	//viper.SetConfigName("sample")

	// Optional if the config has a file extension
	runtime_viper.SetConfigType("json")
	//viper.AddConfigPath("CloudNativeGo/Chapter10/config-yaml/decode/")
	//viper.AddConfigPath("CloudNativeGo/Chapter10/config-json/decode")
	//viper.AddConfigPath("CloudNativeGo/Chapter10/viper/")
	runtime_viper.AddRemoteProvider("consul", "consul-server1:8500", "/config/service.json")

	//if err := viper.ReadInConfig(); err != nil {
	//	panic(fmt.Errorf("fatal error reading config: %w", err))
	//}

	if err := runtime_viper.ReadRemoteConfig(); err != nil {
		panic(fmt.Errorf("fatal error reading config: %w", err))
	}
	var runtime_conf Config
	// unmarshal config
	runtime_viper.Unmarshal(&runtime_conf)

	go func() {

		for {
			time.Sleep(time.Second * 5) // delay after each request

			// currently, only tested with consul support
			err := runtime_viper.WatchRemoteConfig()
			if err != nil {
				log.Printf("unable to read remote config: %v", err)
				continue
			}
			// unmarshal new config into our runtime config struct. you can also use channel
			// to implement a signal to notify the system of the changes
			runtime_viper.Unmarshal(&runtime_conf)
		}
	}()
	//viper.WatchRemoteConfig()
	//viper.OnConfigChange(func(e fsnotify.Event) {
	//	fmt.Println("Config file changed:", e.Name)
	//})
	fmt.Println(runtime_viper.GetString("host"))

	r := mux.NewRouter()
	r.HandleFunc("/", getTagHandler)
	log.Fatal(http.ListenAndServe(":8080", r))
}
