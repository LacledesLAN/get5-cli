package get5

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Config represents a get5 configuration
type Config struct {
	MatchID        string `json:"matchid"`
	NumMaps        uint   `json:"num_maps"`
	PlayersPerTeam byte   `json:"players_per_team"`
	// MinPlayersToReady is the number of players a team must have ready to begin
	MinPlayersToReady byte `json:"min_player_to_ready"`
	// MinSpectatorsToReady is the number of spectators that must be ready to begin
	MinSpectatorsToReady byte   `json:"min_spectators_to_ready"`
	SkipVeto             bool   `json:"skip_veto"`
	VetoFirst            string `json:"vetofirst"`
	SideType             string `json:"side_type"`
	// Spectators contains players that are allow to spectate
	Spectators Spectators `json:"spectators"`
	MapList    []string   `json:"maplist"`
	// Team1 starts as Counter-Terrorists (mp_team1)
	Team1 Team `json:"team1"`
	// Team2 starts as Terrorists (mp_team2)
	Team2 Team `json:"team2"`
	// Cvars that will be executed on each map start or config load.
	Cvars map[string]string `json:"cvars"`
}

func sanitizeConfig(c *Config) {
	// MapList should have no empty elements or elements with whitespace
	var maps []string
	for _, m := range c.MapList {
		m = strings.TrimSpace(m)

		if len(m) > 0 {
			maps = append(maps, m)
		}
	}
	c.MapList = maps

	// can't have 0 maps; derive from number of elements in MapList
	if c.NumMaps < 1 {
		c.NumMaps = uint(len(c.MapList))
	}

	// filter out any duplicate or whitespace spectators
	if len(c.Spectators.Players) > 0 {
		keys := make(map[string]bool, len(c.Spectators.Players))
		spectators := []string{}

		for _, s := range c.Spectators.Players {
			s = strings.TrimSpace(s)
			if len(s) == 0 {
				continue
			}

			if _, found := keys[s]; !found {
				keys[s] = true
				spectators = append(spectators, s)
			}
		}

		c.Spectators.Players = spectators
	}

	c.Team1.Players = sanitizePlayers(c.Team1.Players)
	c.Team2.Players = sanitizePlayers(c.Team2.Players)
}

// Spectators are players who are allowed to spectate the server
type Spectators struct {
	Players []string `json:"players"`
}

// Team represents a CSGO side
type Team struct {
	Name    string  `json:"name"`
	Tag     string  `json:"tag"`
	Flag    string  `json:"flag"`
	Logo    string  `json:"logo"`
	Players Players `json:"players"`
}

// Players represents connected CSGO clients (including bots)
type Players map[string]string

func (p Players) len() int {
	if p == nil {
		return 0
	}

	return len((map[string]string)(p))
}

func sanitizePlayers(p Players) Players {
	buf := map[string]string{}

	// filter out whitespace-only keys and trim and values
	for k, v := range (map[string]string)(p) {
		k = strings.TrimSpace(k)

		if len(k) == 0 {
			continue
		}

		if _, found := buf[k]; !found {
			buf[k] = strings.TrimSpace(v)
		}
	}

	return Players(buf)
}

// UnmarshalJSON parses the JSON-encoded "Players" data into the parent struct
func (p *Players) UnmarshalJSON(data []byte) error {
	if *p == nil {
		*p = Players{}
	}

	data = bytes.TrimSpace(data)

	if len(data) == 0 {
		return nil
	}

	buf := map[string]string{}
	if err := json.Unmarshal(data, &buf); err == nil {
		for steamID, name := range buf {
			steamID = strings.ToUpper(strings.TrimSpace(steamID))

			if len(steamID) > 0 {
				(map[string]string)(*p)[steamID] = strings.TrimSpace(name)
			}
		}

		return nil
	}

	altBuf := []string{}
	if err := json.Unmarshal(data, &altBuf); err == nil {

		for _, steamID := range altBuf {
			steamID = strings.TrimSpace(steamID)

			if len(steamID) > 0 {
				(map[string]string)(*p)[steamID] = ""
			}
		}

		return nil
	}

	return fmt.Errorf("failed to unmarshal: %q", string(data))
}
