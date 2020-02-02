package get5

type Match struct {
	MatchId                string     `json:"matchid"`
	NumMaps                int        `json:"num_maps"`
	PlayersPerTeam         int        `json:"players_per_team"`
	MinPlayersToReady      int        `json:"min_player_to_ready"`
	MinSpectatorsToReady   int        `json:"min_spectators_to_ready"`
	SkipVeto               string     `json:"skip_veto"`
	VetoFirst              string     `json:"vetofirst"`
	SideType               string     `json:"side_type"`
	Spectators             Spectators `json:"spectators"`
	MapList                []string   `json:"maplist"`
	FavoredPercentageTeam1 int        `json:"favored_percentage_team1"`
	FavoredPercentageText  string     `json:"favored_percentage_text"`
	Team1                  Team1      `json:"team1"`
	Team2                  Team2      `json:"team2"`
	Cvars                  Cvars      `json:"cvars"`
}

type Spectators struct {
	Players []Players
}

type Team1 struct {
	Name    string
	Tag     string
	Flag    string
	Logo    string
	Players Players
}

type Team2 struct {
	Name    string
	Tag     string
	Flag    string
	Logo    string
	Players Players
}

type Cvars struct {
	HostName string
}

type Players struct {
	Name string
}
