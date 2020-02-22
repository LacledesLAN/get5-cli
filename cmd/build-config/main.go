package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LacledesLAN/get5-cli/pkg/get5"
	"github.com/jessevdk/go-flags"
)

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

var opts struct {
	CfgFile   string   `long:"cfg" description:"full path the get5-cli configuration file"`
	MatchID   string   `long:"id" description:"A unique ID to identify the get5 match" required:"true"`
	Maplist   []string `long:"map" description:"list of maps to use for the get5 match; must be an odd number" required:"true"`
	Team1Name string   `long:"team1" description:"The name for team1" required:"true"`
	Team2Name string   `long:"team2" description:"The name for team2" required:"true"`
}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		fmt.Println("Couldn't parse arguments from command line")
		os.Exit(1)
	}

	wrapperCfg := Config{}
	if len(strings.TrimSpace(opts.CfgFile)) == 0 {
		path, err := os.Getwd()

		if err != nil {
			fmt.Println("Couldn't determine current working directory")
			os.Exit(1)
		}

		f := filepath.Join(path, "get5-wrapper.json")
		fmt.Printf("Loading from: %s\n", f)
		if err := LoadConfig(f, &wrapperCfg); err != nil {
			panic(err)
		}
	} else {
		if err := LoadConfig(opts.CfgFile, &wrapperCfg); err != nil {
			fmt.Printf("Error loading get5-cli configuration file %q: %s\n", opts.CfgFile, err)
			os.Exit(87)
		}
	}

	//Removing duplicate maps
	fmt.Println(opts.Maplist)
	opts.Maplist = checkUniqueMaps(opts.Maplist)
	fmt.Println(opts.Maplist)

	if len(opts.Maplist)%2 == 0 {
		fmt.Println("Must provide an odd number of maps")
		os.Exit(1)
	}

	get5Cfg := &get5.Config{}
	if err := get5.FromFile(wrapperCfg.Paths.Input, get5Cfg); err != nil {
		fmt.Printf("Error loading input get5-cli configuration file %q: %s\n", wrapperCfg.Paths.Input, err)
		os.Exit(42)
	}

	if len(strings.TrimSpace(opts.MatchID)) == 0 {
		get5Cfg.MatchID = strings.TrimSpace(opts.MatchID)
	}

	if len(opts.Maplist) == 0 {
		get5Cfg.MapList = opts.Maplist
	}

	if len(strings.TrimSpace(opts.Team1Name)) == 0 {
		get5Cfg.Team1.Name = strings.TrimSpace(opts.Team1Name)
	}

	if len(strings.TrimSpace(opts.Team2Name)) == 0 {
		get5Cfg.Team2.Name = strings.TrimSpace(opts.Team2Name)
	}

	get5Cfg.SaveFile(wrapperCfg.Paths.Output)
}

func checkUniqueMaps(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key, _ := range encountered {
		result = append(result, key)
	}
	return result
}
