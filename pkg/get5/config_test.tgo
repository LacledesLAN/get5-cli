package get5

import (
	"math"
	"strings"
	"testing"
)

func Test_sanitizeConfig(t *testing.T) {
	t.Parallel()

	t.Run("Empty and whitespace MatchID should get replaced with generated ID", func(t *testing.T) {
		for _, input := range []string{"", " ", " \t\r\n\v "} {
			sut := Config{MatchID: input}
			sanitizeConfig(&sut)

			if len(sut.MatchID) == 0 || sut.MatchID == input {
				t.Errorf("input `%#v` should have been replaced with a generated one", input)
			}
		}

		for _, input := range []string{"a", " test ", "match 250\t"} {
			sut := Config{MatchID: input}
			sanitizeConfig(&sut)

			if sut.MatchID != strings.TrimSpace(input) {
				t.Errorf("input %#v should have been sanitized to %#v NOT %#v", input, strings.TrimSpace(input), sut.MatchID)
			}
		}
	})

	t.Run("VetoFirst should be `team2` or default to `team1`", func(t *testing.T) {
		tests := map[string]string{"": "team1", " \t\r\n": "team1", "TeAM1": "team1", "hello": "team1", " team2": "team2", "tEaM2\t": "team2"}

		for input, expected := range tests {
			sut := Config{VetoFirst: input}
			sanitizeConfig(&sut)

			if sut.VetoFirst != expected {
				t.Errorf("Input `%#v` should have been sanitized to `%s` not `%s`", input, expected, sut.VetoFirst)
			}
		}
	})

	t.Run("PlayersPerTeam should default to 5", func(t *testing.T) {
		for _, input := range []byte{0, 254, 255} {
			sut := Config{PlayersPerTeam: input}
			sanitizeConfig(&sut)

			if sut.PlayersPerTeam != 5 {
				t.Errorf("Input %d should have been set to default of 5 NOT %d", input, sut.PlayersPerTeam)
			}
		}

		for _, input := range []byte{1, 3, 5, 200} {
			sut := Config{PlayersPerTeam: input}
			sanitizeConfig(&sut)

			if sut.PlayersPerTeam != input {
				t.Errorf("Input %d should have been left alone but got changed to %d", input, sut.PlayersPerTeam)
			}
		}
	})

	t.Run("MinPlayersToReady should default to 1", func(t *testing.T) {
		for _, input := range []byte{0, 49, 254, 255} {
			sut := Config{MinPlayersToReady: input}
			sanitizeConfig(&sut)

			if sut.MinPlayersToReady != 1 {
				t.Errorf("Input %d should have been set to default of 5 NOT %d", input, sut.MinPlayersToReady)
			}
		}

		for _, input := range []byte{1, 2, 5, 12, 48} {
			sut := Config{MinPlayersToReady: input}
			sanitizeConfig(&sut)

			if sut.MinPlayersToReady != input {
				t.Errorf("Input %d should have been left alone but got changed to %d", input, sut.MinPlayersToReady)
			}
		}
	})

	t.Run("SideType should be `always_knife`, `never_knife`, or default to `standard`", func(t *testing.T) {
		tests := map[string]string{
			"always_knife": "always_knife", "always_knife\t": "always_knife", "ALWAYS_knIfe": "always_knife",
			"never_knife": "never_knife", " never_knife\r\n ": "never_knife", "nEVER_KNIFE ": "never_knife",
			"": "standard", " \r\n\t \v": "standard", "standard": "standard", "hello": "standard", "42": "standard", "STAnDARD": "standard",
		}

		for input, expected := range tests {
			sut := Config{SideType: input}
			sanitizeConfig(&sut)

			if sut.SideType != expected {
				t.Errorf("Input `%#v` should have been sanitized to `%s` not `%s`", input, expected, sut.SideType)
			}
		}
	})

	t.Run("MapList should have no empty elements or elements with whitespace", func(t *testing.T) {
		cfg := &Config{
			MapList: []string{"", "  ", " \r\n \t", "  de_depot\t", "", "   \v "},
		}

		sanitizeConfig(cfg)

		if len(cfg.MapList) != 1 {
			t.Errorf("After filtering empty and whitespace values %d map should have remained, but had %d", 1, len(cfg.MapList))
		}

		if cfg.MapList[0] != "de_depot" {
			t.Errorf("After removing padding whitespace the map should have been %q not %q", "de_depot", cfg.MapList[0])
		}
	})

	t.Run("Filtered MapList should maintain the same order", func(t *testing.T) {
		cfg := &Config{
			MapList: []string{"", "  one  ", " \r\n \t", " two\t", "", "   \v "},
		}

		sanitizeConfig(cfg)

		if cfg.MapList[0] != "one" {
			t.Errorf("After filtering the first map should have been %q not %q", "one", cfg.MapList[0])
		}

		if cfg.MapList[1] != "two" {
			t.Errorf("After filtering the first map should have been %q not %q", "two", cfg.MapList[1])
		}
	})

	t.Run("can't have 0 number of maps when maps exist in the MapList", func(t *testing.T) {
		cfg := &Config{
			MapList: []string{"", "one", " ", "two", "\t", "three"},
		}

		sanitizeConfig(cfg)

		if cfg.NumberOfMaps != 3 {
			t.Errorf("Number of maps should have been reported as %d NOT %d", 3, cfg.NumberOfMaps)
		}
	})

	t.Run("filter out any duplicate or whitespace spectators", func(t *testing.T) {
		cfg := &Config{
			Spectators: Spectators{
				Players: []string{"", " \r\n \t \v ", "one", "  one", " one ", "one  ", "two"},
			},
		}

		sanitizeConfig(cfg)

		if len(cfg.Spectators.Players) != 2 {
			t.Errorf("After filtering spectators %d spectators should have remained, but had %d", 2, len(cfg.Spectators.Players))
		}
	})

	t.Run("cvars can't have empty or whitespace keys or values", func(t *testing.T) {
		cfg := &Config{
			Cvars: map[string]string{"": "", " ": "test", "test": "\r\n ", " good ": "123.456", "\t\v": "test", "  ": "test", "ok": " yes",
				"test2": "\r\n\t\v"},
		}
		sanitizeConfig(cfg)

		if len(cfg.Cvars) != 2 {
			t.Error("All cvars with empty and/or whitespace keys and values should have been removed")
		}

		if cfg.Cvars["good"] != "123.456" {
			t.Error("Key `good` should have had value `123.456`")
		}

		if cfg.Cvars["ok"] != "yes" {
			t.Error("Key `ok` should have had value `yes`")
		}
	})
}

func Test_sanitizeTeam(t *testing.T) {
	t.Parallel()

	t.Run("SeriesScore stays in range", func(t *testing.T) {
		tests := map[int]int{math.MinInt32: 0, -1: 0, 0: 0, 1: 1, 3: 3}

		for input, expected := range tests {
			sut := sanitizeTeam(Team{SeriesScore: input})

			if sut.SeriesScore != expected {
				t.Errorf("With an input of `%d` expected `%d` but got `%d`", input, expected, sut.SeriesScore)
			}
		}
	})
}

func Test_sanitizePlayers(t *testing.T) {
	t.Parallel()

	testdata := map[string]string{
		"":       "this should get filtered out",
		" one  ": "\r\n\t\v ",
		" two ":  "  my nickname\r\n",
	}

	sut := sanitizePlayers(Players(testdata))

	if sut.len() != 2 {
		t.Errorf("After sanitization there should have been %d players NOT %d", 2, sut.len())
	}

	if _, ok := (map[string]string)(sut)[""]; ok {
		t.Error("Element with empty key should have been removed")
	}

	if value, ok := (map[string]string)(sut)["one"]; !ok || value != "" {
		t.Errorf("Player with steam id %q should have remained with an empty nickname", "one")
	}

	if value, ok := (map[string]string)(sut)["two"]; !ok || value != "my nickname" {
		t.Errorf("Player with steam id %q should have remained with a nickname of %q", "two", "my nickname")
	}
}

func Test_Players_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("Valid Tests", func(t *testing.T) {
		t.Parallel()

		tests := map[string]struct {
			input       []byte
			expectedLen int
		}{
			"empty":                   {input: []byte{}, expectedLen: 0},
			"empty array":             {input: []byte(`[]`), expectedLen: 0},
			"array with empty string": {input: []byte(`[""]`), expectedLen: 0},
			"empty object":            {input: []byte(`{}`), expectedLen: 0},
			"whitespace":              {input: []byte(" \n  \r\n \t \v"), expectedLen: 0},
			"null literal":            {input: []byte(`null`), expectedLen: 0},
			"lone steam id":           {input: []byte(`["STEAM_1:6:12345678"]`), expectedLen: 1},
			"steam ids only": {
				input: []byte(`[
					"STEAM_1:0:12345678",
					"STEAM_1:2:12345678",
					"STEAM_1:3:12345678",
					"STEAM_1:4:12345678",
					"STEAM_1:5:12345678",
					"STEAM_1:6:12345678"
				]`),
				expectedLen: 6,
			},
			"steam id without alias": {input: []byte(`{"STEAM_1:0:12345678": ""}`), expectedLen: 1},
			"steam id with alias":    {input: []byte(`{"STEAM_1:0:12345678": "Ava"}`), expectedLen: 1},
			"steam ids with aliases": {
				input: []byte(`		{
					"STEAM_1:0:12345678" : "Ava",
					"STEAM_1:1:12345678" : "Oliver",
					"STEAM_1:2:12345678" : "Mia",
					"STEAM_1:3:12345678" : "Mason",
					"STEAM_1:4:12345678" : "Daniel",
					"STEAM_1:5:12345678" : "Olivia"
				}`),
				expectedLen: 6,
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				var sut Players

				if err := (&sut).UnmarshalJSON(test.input); err != nil {
					t.Fatalf("Error unmarshalling: %s", err.Error())
				}

				if sut.len() != test.expectedLen {
					t.Fatalf("Expected a len of %d but got %d", test.expectedLen, sut.len())
				}
			})
		}
	})

	t.Run("Invalid Tests", func(t *testing.T) {
		t.Parallel()

		tests := [][]byte{
			[]byte("aa"), []byte("dW5kZWZbAmVk"), []byte("Tm9uZQ=="),
			[]byte("OTk5OTk5OTk5OTk5OTk5OTk5OTk5OTk5OTk5OTk5OTk5OTk5OTk5OTk5OTk5OTk5OTk5OTk5OTk5"),
			[]byte("77u/"), []byte("undefined"), []byte("true"), []byte("false"), []byte("\\"), []byte("\\\\"),
			[]byte("0/0"), []byte("사회과학원 어학연구소"), []byte("・(￣∀￣)・:*:"), []byte("👾 🙇 💁 🙅 🙆 🙋 🙎 🙍"),
		}

		for _, test := range tests {
			var sut Players

			if err := (&sut).UnmarshalJSON(test); err == nil {
				t.Fatalf("Unmarshalling %q input should have resulted in an error but got nil", string(test))
			}

			if sut.len() != 0 {
				t.Fatalf("A failed unmarshalling should have resulted in a len of 0 not %d", sut.len())
			}
		}
	})
}
