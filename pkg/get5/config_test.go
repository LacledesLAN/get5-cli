package get5

import "testing"

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
