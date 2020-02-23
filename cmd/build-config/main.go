package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LacledesLAN/get5-cli/pkg/get5"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	BaseFile       string   `long:"basefile" description:"full path to the base get5 configuration file to load"`
	CfgFile        string   `long:"cfgfile" description:"full path to the get5-cli configuration file"`
	MatchID        string   `long:"id" description:"A unique ID to identify the get5 match" required:"true"`
	Maplist        []string `long:"map" description:"list of maps to use for the get5 match; must be an odd number" required:"true"`
	MinPlayers     byte     `long:"minready" description:"The minimum players a team needs to be able to ready up" default:"5"`
	PlayersPerTeam byte     `long:"teamsize" description:" The maximum players per team (doesn't include a coach spot)" default:"5"`
	Team1Name      string   `long:"team1" description:"The name for team1" required:"true"`
	Team2Name      string   `long:"team2" description:"The name for team2" required:"true"`
}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		fmt.Printf("Encountered error while parsing arguments from the command line: %s\n", err)
		os.Exit(8)
	}

	if len(opts.Maplist)%2 < 1 && len(opts.Maplist)%2 == 0 {
		fmt.Printf("Must provide a positive, odd number of maps; got: %v\n", opts.Maplist)
		os.Exit(61)
	}

	if len(strings.TrimSpace(opts.CfgFile)) == 0 {
		path, err := os.Getwd()

		if err != nil {
			fmt.Printf("Encountered error while determining the current working directory: %s\n", err)
			os.Exit(2)
		}

		opts.CfgFile = filepath.Join(path, "get5-wrapper.json")
	}

	wrapperCfg := Config{}
	if err := LoadConfig(opts.CfgFile, &wrapperCfg); err != nil {
		fmt.Printf("Encountered error loading cli configuration file: %s\n", err)
		os.Exit(22)
	}
	fmt.Printf("Loaded cli configuration file: %s\n", opts.CfgFile)

	if len(strings.TrimSpace(opts.BaseFile)) > 0 {
		wrapperCfg.Paths.Input = strings.TrimSpace(opts.BaseFile)
	}

	get5Cfg := &get5.Config{}
	if err := get5.FromFile(wrapperCfg.Paths.Input, get5Cfg); err != nil {
		fmt.Printf("Encountered error loading base get5 configuration file %q: %s\n", wrapperCfg.Paths.Input, err)
		os.Exit(42)
	}
	fmt.Printf("Loaded base get5 configuration file: %s\n", wrapperCfg.Paths.Input)

	fmt.Print("\nModifying base get5 configuration\n")

	fmt.Printf("\t• Match ID to %q\n", opts.MatchID)
	get5Cfg.MatchID = strings.TrimSpace(opts.MatchID)

	fmt.Printf("\t• Map list to: %v\n", opts.Maplist)
	get5Cfg.MapList = opts.Maplist
	get5Cfg.NumberOfMaps = len(opts.Maplist)

	fmt.Printf("\t• Minimum players per team to ready up to %d\n", opts.MinPlayers)
	get5Cfg.MinPlayersToReady = opts.MinPlayers

	fmt.Printf("\t• Team 1's name to %q\n", opts.Team1Name)
	get5Cfg.Team1.Name = strings.TrimSpace(opts.Team1Name)

	fmt.Printf("\t• Team 2's name to %q\n", opts.Team2Name)
	get5Cfg.Team2.Name = strings.TrimSpace(opts.Team2Name)

	fmt.Printf("\t• Team size to: %d\n", opts.PlayersPerTeam)
	get5Cfg.PlayersPerTeam = opts.PlayersPerTeam

	if ok, issues := get5Cfg.Validate(); !ok {
		fmt.Print("\nget5 configuration failed validation:\n")

		for _, issue := range issues {
			fmt.Printf("\t• %s\n", issue)
		}

		os.Exit(124)
	}

	if err := get5Cfg.SaveFile(wrapperCfg.Paths.Output); err != nil {
		fmt.Printf("Encountered error saving get5 configuration file to %q: %s", wrapperCfg.Paths.Output, err)
		os.Exit(43)
	}

	fmt.Printf("\nSaved get5 configuration file to %q\n", wrapperCfg.Paths.Output)
}
