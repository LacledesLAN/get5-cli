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
	CfgFile           string   `long:"cfg" description:"full path the get5-cli configuration file"`
	MatchID           string   `long:"id" description:"A unique ID to identify the get5 match" required:"true"`
	Maplist           []string `long:"map" description:"list of maps to use for the get5 match; must be an odd number" required:"true"`
	MinPlayersToReady byte     `long:"playerstoready" description:"Minimum players a team needs to be able to ready up"`
}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		fmt.Println("Couldn't parse arguments from command line")
		os.Exit(1)
	}

	if len(strings.TrimSpace(opts.CfgFile)) == 0 {
		path, err := os.Getwd()

		if err != nil {
			fmt.Println("Couldn't determine current working directory")
			os.Exit(1)
		}

		opts.CfgFile = filepath.Join(path, "get5-cli.json")
	}

	cfg := &Config{}
	if err := LoadConfig(opts.CfgFile, cfg); err != nil {
		fmt.Printf("Error loading get5-cli configuration file %q: %s\n", opts.CfgFile, err)
		os.Exit(1)
	}

	if len(opts.Maplist)%2 == 0 {
		fmt.Println("Must provide an odd number of maps")
		os.Exit(1)
	}

	c := &get5.Config{}
	if err := get5.FromFile(cfg.Paths.Input, c); err != nil {
		fmt.Printf("Error loading input get5-cli configuration file %q: %s\n", cfg.Paths.Input, err)
		os.Exit(1)
	}

	c.MapList = opts.Maplist

	//c.Save("")
}
