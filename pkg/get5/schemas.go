package get5

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/exp/slices"
)

func sanitizePrintable(raw string, maxLength int) string {
	raw = strings.Join(strings.Fields(strings.TrimSpace(raw)), "_")

	raw = strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}

		return -1
	}, raw)

	if maxLength < 1 || maxLength >= len(raw) {
		return raw
	}

	return raw[:maxLength]
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// MatchTeam represents a CSGO side (CT/T).
type MatchTeam struct {
	// Players contains the players on the team.
	Players json.RawMessage `json:"players"`
	// Coaches, similarly to Players, this object maps coaches using their Steam ID and name, locking them to the coach slot unless removed using
	// get5_removeplayer. Setting a Steam ID as coach takes precedence over being set as a player.
	Coaches json.RawMessage `json:"coaches"`
	// Name is the team's name. Sets mp_teamname_1 or mp_teamname_2. Printed frequently in chat. If you don't define a team name, it will be set to team_
	// followed by the name of the captain.
	Name string `json:"name"`
	// Tag is a short version of the team name, used in clan tags in-game (requires that get5_set_client_clan_tags is disabled).
	Tag string `json:"tag"`
	// Flag is the ISO-code to use for the in-game flag of the team. Must be a supported country, i.e. FR,UK,SE etc.
	Flag string `json:"flag"`
	// Logo The team logo (wraps mp_teamlogo_1 or mp_teamlogo_2), which requires to be on a FastDL in order for clients to see.
	Logo string `json:"logo"`
	// SeriesScore is the current score in the series, this can be used to give a team a map advantage or used as a manual backup method.
	SeriesScore int `json:"series_score"`
	// MatchText wraps mp_teammatchstat_1 and mp_teammatchstat_2. You probably don't want to set this, in BoX series, mp_teamscore cvars are automatically
	// set and take the place of the mp_teammatchstat_x cvars.
	MatchText string `json:"matchtext"`
}

func sanitizeMatchTeam(mt MatchTeam) MatchTeam {
	mt.Name = sanitizePrintable(mt.Name, 30)
	mt.Tag = sanitizePrintable(mt.Tag, 0)
	mt.Flag = sanitizePrintable(mt.Flag, 2)
	mt.Logo = sanitizePrintable(mt.Logo, 0)

	if mt.SeriesScore < 0 {
		mt.SeriesScore = 0
	}

	mt.MatchText = sanitizePrintable(mt.MatchText, 0)

	return mt
}

// Match represents a get5 match
type Match struct {
	// MatchTitle is a wrapper of the server's mp_teammatchstat_txt cvar, but can use {MAPNUMBER} and {MAXMAPS} as variables that get replaced with their
	// integer values. In a BoX series, you probably don't want to set this since Get5 automatically sets mp_teamscore cvars for the current series score,
	// and take the place of the mp_teammatchstat cvars. Default: "Map {MAPNUMBER} of {MAXMAPS}"
	MatchTitle string `json:"match_title,omitempty"`
	// MatchID is the ID of the match. This determines the matchid parameter in all forwards and events. If you use the MySQL extension, you should leave
	// this field blank (or omit it), as match IDs will be assigned automatically. If you do want to assign match IDs from another source, they must be
	// integers (in a string) and must increment between matches. Default: ""
	MatchID string `json:"matchid,omitempty"`
	// ClinchSeries If false, the entire map list will be played, regardless of score. If true, a series will be won when the series score for a team
	// exceeds the number of maps divided by two. Default: true
	ClinchSeries *bool `json:"clinch_series,omitempty"`
	// NumberOfMaps The number of maps to play in the series; must be positive, odd number
	NumberOfMaps *int `json:"num_maps,omitempty"`
	// PlayersPerTeam is the number of players per team. You should never set this to a value higher than the number of players you want to actually play in
	// a game, excluding coaches.
	PlayersPerTeam *byte `json:"players_per_team,omitempty"`
	// CoachesPerTeam is the maximum number of coaches per team.
	CoachesPerTeam *byte `json:"coaches_per_team,omitempty"`
	// MinPlayersToReady is the minimum number of players that must be present for the !forceready command to succeed. If not forcing a team ready, all
	// players must !ready up themselves. Default: 0
	MinPlayersToReady *byte `json:"min_player_to_ready,omitempty"`
	// MinSpectatorsToReady is the minimum number of spectators that must be !ready for the game to begin. Default: 0
	MinSpectatorsToReady *byte `json:"min_spectators_to_ready,omitempty"`
	// SkipVeto determines whether to skip the veto phase. When skipping veto, map_sides determines sides, and if map_sides is not set, sides are determined
	// by side_type. Default: false
	SkipVeto *bool `json:"skip_veto,omitempty"`
	// VetoFirst is The team that vetoes first. Default: "team1". Allowed values are "team1", "team2", and "random".
	VetoFirst string `json:"vetofirst,omitempty"`
	// SideType is the method used to determine sides when vetoing or if veto is disabled and map_sides are not set. This parameter is ignored if map_sides
	// is set for all maps. standard and always_knife behave similarly when skip_veto is true.
	//	"standard" means that the team that doesn't pick a map gets the side choice (only if skip_veto is false).
	//	"always_knife" means that sides are always determined by a knife-round.
	//	"never_knife" means that team1 always starts on CT."
	SideType string `json:"side_type,omitempty"`
	// MapSides Determines the starting sides for each map. If this array is shorter than num_maps, side_type will determine the side-behavior of the
	// remaining maps. Ignored if skip_veto is false. Allowed values are "team1_ct", "team1_t", "knife".
	MapSides []string `json:"map_sides,omitempty"`
	// Spectators is the spectators to allow into the game. If not defined, spectators cannot join the game.
	Spectators json.RawMessage `json:"spectators,omitempty"`
	// MapList is the map pool to pick from, as an array of strings (["de_dust2", "de_nuke"] etc.), or if skip_veto is true, the order of maps played
	// (limited by num_maps). This should always be odd-sized if using the in-game veto system.
	MapList []string `json:"maplist"`
	// FavoredPercentageTeam1 is a wrapper for the server's mp_teamprediction_pct. This determines the chances of team1 winning. Default: 0
	FavoredPercentageTeam1 *byte `json:"favored_percentage_team1"`
	// FavoredPercentageText is a wrapper for the server's mp_teamprediction_txt. Default: "".
	FavoredPercentageText string `json:"favored_percentage_text"`
	// Team1 starts as Counter-Terrorists (mp_team1)
	Team1 MatchTeam `json:"team1,omitempty"`
	// Team2 starts as Terrorists (mp_team2)
	Team2 MatchTeam `json:"team2,omitempty"`
	// Cvars contains various commands to execute on the server when loading the match configuration. This can be both regular server-commands and any Get5
	// configuration parameter, i.e. {"hostname": "Match #3123 - Red vs. Blu"}
	Cvars map[string]string `json:"cvars,omitempty"`
}

func sanitizeMatch(m *Match) error {
	m.MatchTitle = sanitizePrintable(m.MatchTitle, 36)
	m.MatchID = sanitizePrintable(m.MatchID, 0)

	if m.NumberOfMaps != nil && (*m.NumberOfMaps < 1 || *m.NumberOfMaps%2 == 0) {
		return errors.New("NumberOfMaps must be a positive, odd-valued integer")
	}

	m.VetoFirst = strings.ToLower(m.VetoFirst)
	if !slices.Contains([]string{"team1", "team2", "random"}, m.VetoFirst) {
		m.VetoFirst = ""
	}

	m.SideType = strings.ToLower(m.SideType)
	if !slices.Contains([]string{"standard", "always_knife", ""}, "never_knife") {
		m.SideType = ""
	}

	if m.MapSides != nil {
		if len(m.MapSides) == 0 {
			m.MapSides = nil
		} else {
			for i, v := range m.MapSides {
				m.MapSides[i] = strings.ToLower(v)

				if !slices.Contains([]string{"team1_ct", "team1_t", "knife"}, v) {
					return fmt.Errorf("MapSides: '%s' is not a valid value", v)
				}
			}
		}
	}

	if m.MapList == nil || len(m.MapList) == 0 {
		return errors.New("maplist cannot be empty")
	}

	if m.SkipVeto == nil || !*m.SkipVeto {
		if len(m.MapList)%2 == 0 {
			return errors.New("maplist must contain an odd number of maps")
		}
	}

	m.FavoredPercentageText = sanitizePrintable(m.FavoredPercentageText, 0)

	m.Team1 = sanitizeMatchTeam(m.Team1)
	m.Team2 = sanitizeMatchTeam(m.Team2)

	return nil
}
