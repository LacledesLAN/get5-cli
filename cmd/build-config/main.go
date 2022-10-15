package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/LacledesLAN/get5-cli/pkg/get5"
	"github.com/jessevdk/go-flags"
)

// GitCommitHash is used for storing the git commit hash that built this instance of `get5-cli`.
var GitCommitHash string

// opts contains the command-line arguments for `get5-cli`.
var opts struct {
	// Required Arguments
	Maplist   []string `long:"map-list" short:"m" description:"list of maps to use for the get5 match; must be an odd number" required:"true"`
	Team1Name string   `long:"team1-name" short:"1" description:"the name for team1" required:"true"`
	Team2Name string   `long:"team2-name" short:"2" description:"The name for team2" required:"true"`

	// Meta Optional Arguments
	BaseSchemaPath string `long:"base-schema" description:"path to the base get5 configuration file to load (default: ./csgo/base-schema.json)"`
	DestSchemaPath string `long:"dest-schema" description:"path to the destination schema to generate (default './csgo/automatch.json')"`

	// Optional Arguments
	Cvars             map[string]string `long:"cvar" short:"v" description:"a CSGO cvar"`
	MatchTitle        string            `long:"match-title" description:"Commands to execute when the match configuration is loaded (cvars and commands)"`
	MinPlayersToReady byte              `long:"min-ready" description:"The minimum players a team needs to be able to ready up"`
	NumberOfMaps      byte              `long:"map-count" short:"c" description:"the number of maps to include in the series (default 3)"`
	PlayersPerTeam    byte              `long:"team-size" description:" The maximum players per team (doesn't include a coach spot)"`
	Team1Score        byte              `long:"team1-score" description:"the team's current score in the series"`
	Team2Score        byte              `long:"team2-score" description:"the team's current score in the series"`
}

func main() {
	if len(GitCommitHash) > 0 {
		fmt.Printf("get5-cli version: %s\n", GitCommitHash)
	}

	//
	// VALIDATE COMMAND LINE ARGUMENTS
	//
	if _, err := flags.Parse(&opts); err != nil {
		fmt.Printf("Encountered error while parsing arguments from the command line: %s\n", err)
		os.Exit(8)
	}

	opts.BaseSchemaPath = strings.TrimSpace(opts.BaseSchemaPath)
	if len(opts.BaseSchemaPath) == 0 {
		opts.BaseSchemaPath = "./csgo/base-schema.json"
	}

	opts.DestSchemaPath = strings.TrimSpace(opts.DestSchemaPath)
	if len(opts.DestSchemaPath) == 0 {
		opts.DestSchemaPath = "./csgo/automatch.json"
	}

	if opts.NumberOfMaps > 0 && opts.NumberOfMaps%2 == 0 {
		fmt.Printf("Number of Maps must be an odd number")
		os.Exit(1)
	}

	if opts.Maplist == nil || len(opts.Maplist) < 1 {
		fmt.Printf("A list of maps must be provided")
		os.Exit(1)
	} else if opts.NumberOfMaps > 0 && len(opts.Maplist) < int(opts.NumberOfMaps) {
		fmt.Printf("List of maps must contain as least %d maps in it", opts.NumberOfMaps)
		os.Exit(1)
	}

	if len(opts.Maplist)%2 == 0 {
		fmt.Printf("Must provide a positive, odd number of maps; got: %v\n", opts.Maplist)
		os.Exit(61)
	}

	opts.Team1Name = strings.TrimSpace(opts.Team1Name)
	if len(opts.Team1Name) < 1 {
		fmt.Printf("Team 1 Name must be provided")
		os.Exit(1)
	}

	opts.Team2Name = strings.TrimSpace(opts.Team2Name)
	if len(opts.Team2Name) < 1 {
		fmt.Printf("Team 2 Name must be provided")
		os.Exit(1)
	}

	//
	// LOAD BASE SCHEMA
	//
	wipSchema := get5.Match{}
	if err := get5.FromFile(opts.BaseSchemaPath, &wipSchema); err != nil {
		fmt.Printf("Encountered error loading base schema file: %s\n", err)
		os.Exit(22)
	}
	fmt.Printf("Loaded base schema file: %s\n", opts.BaseSchemaPath)

	//
	// MODIFY WIP SCHEMA
	//
	fmt.Println("\nModifying base get5 configuration")

	// From Required Arguments
	fmt.Printf("\t• Setting maplist to: [%s]\n", opts.Maplist)
	wipSchema.MapList = opts.Maplist

	fmt.Printf("\t• Setting team 1 name to %s\n", opts.Team1Name)
	wipSchema.Team1.Name = opts.Team1Name

	fmt.Printf("\t• Setting team 2 name to %s\n", opts.Team2Name)
	wipSchema.Team2.Name = opts.Team2Name

	// From Optional Arguments
	if opts.Cvars != nil && len(opts.Cvars) > 1 {
		if wipSchema.Cvars == nil {
			wipSchema.Cvars = make(map[string]string, len(opts.Cvars))
		}

		for n, v := range opts.Cvars {
			fmt.Printf("\t• Setting match cvar '%s' to '%s'\n", n, v)
			wipSchema.Cvars[n] = v
		}
	}

	opts.MatchTitle = strings.TrimSpace(opts.MatchTitle)
	if len(opts.MatchTitle) > 0 {
		fmt.Printf("\t• Setting match title to %s\n", opts.MatchTitle)
		wipSchema.MatchTitle = opts.MatchTitle
	}

	if opts.MinPlayersToReady > 0 {
		fmt.Printf("\t• Setting min players ready to %d\n", opts.MinPlayersToReady)
		wipSchema.MinPlayersToReady = &opts.MinPlayersToReady
	}

	if opts.NumberOfMaps > 0 {
		fmt.Printf("\t• Setting number of maps to %d\n", opts.NumberOfMaps)
		*wipSchema.NumberOfMaps = int(opts.NumberOfMaps)
	}

	if opts.PlayersPerTeam > 0 {
		fmt.Printf("\t• Setting players per team to %d\n", opts.PlayersPerTeam)
		wipSchema.PlayersPerTeam = &opts.PlayersPerTeam
	}

	if opts.Team1Score > 0 {
		fmt.Printf("\t• Setting team 1 series score to %d\n", opts.Team1Score)
		t := int(opts.Team1Score)
		wipSchema.Team1.SeriesScore = &t
	}

	if opts.Team2Score > 0 {
		fmt.Printf("\t• Setting team 2 series score to %d\n", opts.Team2Score)
		t := int(opts.Team2Score)
		wipSchema.Team2.SeriesScore = &t
	}

	if err := get5.SaveFile(wipSchema, opts.DestSchemaPath); err != nil {
		fmt.Printf("Encountered error saving get5 configuration file to %q: %s", opts.DestSchemaPath, err)
		os.Exit(43)
	}

	fmt.Printf("\nSaved get5 configuration file to %q\n", opts.DestSchemaPath)
}
