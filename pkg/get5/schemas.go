package get5

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

// Match represents a get5 match
type Match struct {
	// MatchID is a unique string matchid used to identify the match
	MatchID string `json:"matchid,omitempty"`
	// MatchTitle is the name of the match
	MatchTitle string `json:"match_title,omitempty"`
	// ClinchSeries If false, the entire map list will be played, regardless of score. If true, a series will be won when the series score for a team exceeds the number of maps divided by two.
	// ClinchSeries bool `json:"clinch_series"`
	// NumberOfMaps in the series; must be positive, odd number
	NumberOfMaps int `json:"num_maps"`
	// MapList is the maps in use for the match; should be an odd-sized list
	MapList []string `json:"maplist"`
	// SkipVeto determines whether the veto will be skipped and all the maps will come from the maplist (in the given order)
	SkipVeto bool `json:"skip_veto"`
	// VetoFirst either "team1", or "team2". If not set, or set to any other value, team 1 will veto first.
	VetoFirst string `json:"vetofirst"`
	// SideType either "standard", "never_knife", or "always_knife"; standard means the team that doesn't pick a map gets the side choice, never_knife
	// 	means team1 is always on CT first, and always knife means there is always a knife round.
	SideType string `json:"side_type"`
	// PlayersPerTeam maximum players per team (doesn't include a coach spot, default: 5)
	PlayersPerTeam byte `json:"players_per_team"`
	// MinPlayersToReady is the minimum players a team needs to be able to ready up (default: 1)
	MinPlayersToReady byte `json:"min_player_to_ready"`
	// MinSpectatorsToReady is the number of spectators that must be ready to begin
	MinSpectatorsToReady byte `json:"min_spectators_to_ready"`
	// Spectators contains players that are allow to spectate
	Spectators Spectators `json:"spectators"`
	// Team1 starts as Counter-Terrorists (mp_team1)
	Team1 MatchTeam `json:"team1,omitempty"`
	// Team2 starts as Terrorists (mp_team2)
	Team2 MatchTeam `json:"team2,omitempty"`
	// Cvars that will be executed during match warmup/knife round/live state
	Cvars map[string]string `json:"cvars"`
}

func sanitizeMatch(c *Match) {
	c.MatchID = strings.TrimSpace(c.MatchID)

	c.VetoFirst = strings.TrimSpace(strings.ToLower(c.VetoFirst))
	if c.VetoFirst != "team2" {
		c.VetoFirst = "team1"
	}

	c.SideType = strings.TrimSpace(strings.ToLower(c.SideType))
	if c.SideType != "always_knife" && c.SideType != "never_knife" {
		c.SideType = "standard"
	}

	if c.PlayersPerTeam < 1 || c.PlayersPerTeam >= math.MaxUint8-1 {
		c.PlayersPerTeam = 5
	}

	if c.MinPlayersToReady < 1 || c.MinPlayersToReady > 48 {
		c.MinPlayersToReady = 1
	}

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
	if c.NumberOfMaps < 1 {
		c.NumberOfMaps = len(c.MapList)
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

	// filter out empty/whitespace cvars (both key and value)
	if len(c.Cvars) > 0 {
		buf := map[string]string{}

		for key, value := range c.Cvars {
			key = strings.TrimSpace(key)
			value = strings.TrimSpace(value)

			if len(key) > 0 && len(value) > 0 {
				buf[key] = value
			}
		}

		c.Cvars = buf
	}

	c.Team1 = sanitizeMatchTeam(c.Team1)
	c.Team2 = sanitizeMatchTeam(c.Team2)
}

// Spectators are players who are allowed to spectate the server
type Spectators struct {
	Players []string `json:"players,omitempty"`
}

// Team represents a CSGO side
type MatchTeam struct {
	// Name (wraps mp_teamname_# and is displayed often in chat messages)
	Name string `json:"name"`
	// Tag (or short name), this replaces client "clan tags"
	Tag string `json:"tag"`
	// Flag team flag (2 letter country code, wraps mp_teamflag_#)
	Flag string `json:"flag"`
	// Logo (wraps mp_teamlogo_#)
	Logo string `json:"logo"`
	// Players list of Steam id's for users on the team (not used if get5_check_auths is set to 0). You can also force player names in here; in JSON you may use either an array of steamids or a dictionary of steamids to names.
	Players Players `json:"players"`
	// current score in the series, this can be used to give a team a map advantage or used as a manual backup method, defaults to 0
	SeriesScore int `json:"series_score"`
}

func sanitizeMatchTeam(t MatchTeam) MatchTeam {
	if t.SeriesScore < 0 {
		t.SeriesScore = 0
	}

	t.Players = sanitizePlayers(t.Players)

	return t
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
