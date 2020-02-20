package get5

import "testing"

func Test_sanitizeConfig(t *testing.T) {
	t.Parallel()

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

		if cfg.NumMaps != 3 {
			t.Errorf("Number of maps should have been reported as %d NOT %d", 3, cfg.NumMaps)
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
