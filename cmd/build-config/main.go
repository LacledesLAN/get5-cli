package main

import (
	"net/http"
)

func main() {

	// conf, err := config.GetConfig("cmd/svr/config.json")
	// if err != nil {
	// 	log.Fatalln("Error getting configuration, shutting down...", err)
	// }

	// log.Println("listening at", conf.Port)

	// err = viperConfig.InitViper(conf)
	// if err != nil {
	// 	log.Fatalln("Viper Error: ", err)
	// }

	// args := os.Args[1:] //1: gets all command line args, normal indexing works for single args

	// for index := 0; index < len(args); index++ {
	// 	println("Index ", index, " Arguments", args[index], args[index+1])
	// 	key := args[index]
	// 	index++
	// 	value := args[index]
	// 	//TODO Need to fix updating ints inside config to ints instead of strings
	// 	viperConfig.UpdateParam(key, value, conf)
	// }

	//TODO WEB APP OPTION
	//Start serving app and listen
	//server := createHTTPServer(conf)

	//if err := server.ListenAndServe(); err != nil {
	//	log.Fatal("fatal server error: ", err)
	//}
}

func createHTTPServer() *http.Server {

	//r := mux.NewRouter().StrictSlash(true)
	//server := &http.Server{Addr: ":" + config.Port, Handler: r}

	// //r.Handle("/update-config/{configName}", viperConfig.UpdateConfigWatcher()).Methods(http.MethodPost)
	// r.HandleFunc("/update-param/", func(w http.ResponseWriter, r *http.Request) {
	// 	v := r.URL.Query()
	// 	key := v.Get("key")
	// 	value := v.Get("value")

	// 	println("key passed", key)
	// 	println("value passed", value)

	// 	//viperConfig.UpdateParam(key, value, config)
	// })

	return nil
}
