package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestParseGroupFinderActivityColumnOrder(t *testing.T) {
	want := []GroupFinderActivity{{
		ID: 42, FullName_lang: "Dungeon, the Great", ShortName_lang: "Dungeon",
		GroupFinderCategoryID: 2, OrderIndex: -1, GroupFinderActivityGrpID: 7,
		Flags: 98, MinGearLevelSuggestion: 100, PlayerConditionID: 123,
		MapID: 645, DifficultyID: 8, AreaID: 141, ExpansionID: 3,
		MaxPlayers: 5, DisplayType: 1, OverrideContentTuningID: 456, MapChallengeModeID: 789,
	}}
	tests := map[string]string{
		"classic": "ID,FullName_lang,ShortName_lang,GroupFinderCategoryID,OrderIndex,GroupFinderActivityGrpID,Flags,MinGearLevelSuggestion,PlayerConditionID,MapID,DifficultyID,AreaID,ExpansionID,MaxPlayers,DisplayType,MinLevel,OverrideContentTuningID,MapChallengeModeID\n42,\"Dungeon, the Great\",Dungeon,2,-1,7,98,100,123,645,8,141,3,5,1,80,456,789\n",
		"retail":  "FullName_lang,ShortName_lang,ID,GroupFinderCategoryID,OrderIndex,GroupFinderActivityGrpID,Flags,MinGearLevelSuggestion,PlayerConditionID,MapID,DifficultyID,AreaID,ExpansionID,MaxPlayers,DisplayType,MinLevel,OverrideContentTuningID,MapChallengeModeID\n\"Dungeon, the Great\",Dungeon,42,2,-1,7,98,100,123,645,8,141,3,5,1,80,456,789\n",
	}
	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			if got := parseGroupFinderActivity(strings.NewReader(input)); !reflect.DeepEqual(got, want) {
				t.Fatalf("got %+v, want %+v", got, want)
			}
		})
	}
}

func TestParseMapChallengeModeColumnOrder(t *testing.T) {
	input := "Extra,ExpansionLevel,ID,Name_lang,Field_12_0_0_63854_004,Flags,MapID\nignored,4,2,Temple of the Jade Serpent,3,32,960\n"
	want := []MapChallengeMode{{Name_lang: "Temple of the Jade Serpent", ID: 2, MapID: 960, Flags: 32, Field_12_0_0_63854_004: 3, ExpansionLevel: 4}}
	if got := parseMapChallengeMode(strings.NewReader(input)); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestParseMythicPlusSeasonTrackedMapColumnOrder(t *testing.T) {
	input := "DisplaySeasonID,Extra,MapChallengeModeID,ID\n32,ignored,2,235\n"
	want := []MythicPlusSeasonTrackedMap{{ID: 235, MapChallengeModeID: 2, DisplaySeasonID: 32}}
	if got := parseMythicPlusSeasonTrackedMap(strings.NewReader(input)); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestReadCSVMissingColumn(t *testing.T) {
	defer func() {
		if got := fmt.Sprint(recover()); got != `missing CSV column "ID"` {
			t.Fatalf("unexpected panic: %s", got)
		}
	}()
	readCSV(strings.NewReader("Other\n42\n"), "ID")
}

func TestParseLocalCSVFiles(t *testing.T) {
	parsers := []struct {
		name  string
		parse func(io.Reader) int
	}{
		{"GroupFinderActivity", func(r io.Reader) int { return len(parseGroupFinderActivity(r)) }},
		{"MapChallengeMode", func(r io.Reader) int { return len(parseMapChallengeMode(r)) }},
		{"MythicPlusSeasonTrackedMap", func(r io.Reader) int { return len(parseMythicPlusSeasonTrackedMap(r)) }},
	}
	for _, parser := range parsers {
		for _, suffix := range []string{"", "_Mists"} {
			t.Run(parser.name+suffix, func(t *testing.T) {
				f, err := os.Open("data/" + parser.name + suffix + ".csv")
				if err != nil {
					t.Fatal(err)
				}
				defer f.Close()
				records, err := csv.NewReader(f).ReadAll()
				if err != nil {
					t.Fatal(err)
				}
				if _, err := f.Seek(0, io.SeekStart); err != nil {
					t.Fatal(err)
				}
				if got, want := parser.parse(f), len(records)-1; got != want || got == 0 {
					t.Fatalf("got %d entries, want %d nonempty entries", got, want)
				}
			})
		}
	}
}
