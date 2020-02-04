package viperConfig

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"ll-jsonConverter/pkg/config"
	"ll-jsonConverter/pkg/model"
	"log"
	"net/http"
	"strconv"
)

func InitViper(conf config.Config) error{

	SetDefaults() //if no config found or provided defaults can be set
	viper.SetConfigType("json")
	viper.AddConfigPath(conf.ConfigPath) //multiple paths can be added to watch, must be added before calling WatchConfig()

	//Initial Save
	SaveViper(conf.Name)

	return nil
}

func SaveViper(configName string) error {

	var match model.Match
	viper.SetConfigName(configName)

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(fmt.Errorf("InitViper Fatal error config file: %s \n", err))
	}

	err = viper.Unmarshal(&match)
	if err != nil {
		panic(fmt.Errorf("INITIAL unable to decode into struct: %s \n", err))
	}

	//config changes
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		log.Println(viper.Get("matchid"))

	})
	err = viper.WriteConfigAs(configName) // writes current config to predefined path set by 'viper.AddConfigPath()' and 'viper.SetConfigName'
	if err != nil {
		panic(fmt.Errorf("Failed to Save File %s \n", err))
	}

	//TODO Jordan
	log.Println(viper.Get("matchid"))
	log.Println(viper.Get("maplist"))
	//viper.Set("matchid","hello")

	return nil
}

func UpdateConfigWatcher() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		configName := vars["configName"] //Name of new config to start watching

		SaveViper(configName)

		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {
			panic(fmt.Errorf("UpdateConfigWater: Fatal error config file: %s \n", err))
		}
		return
	})
}

func UpdateParam(key string, value string, conf config.Config) {

		// Check if the value passed in is a string or int, if its in an int, pass an int
		if checkForInt, err := strconv.Atoi(value); err == nil {
			viper.Set(key,checkForInt)
			viper.WriteConfigAs(conf.ConfigPath + conf.Name)
			return
		}
		viper.Set(key,value)
		//TODO Update save viper and split out write vs watch function
		//SaveViper(conf.Name)
		viper.WriteConfigAs(conf.ConfigPath + conf.Name)
		return

}

func SetDefaults(){

	return
}