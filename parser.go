package main

import (
	"0xbs/premade-groups-filter-helper/logger"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

func parseGroupFinderActivity(source io.Reader) []GroupFinderActivity {
	records := readCSV(source,
		"ID",
		"FullName_lang",
		"ShortName_lang",
		"GroupFinderCategoryID",
		"OrderIndex",
		"GroupFinderActivityGrpID",
		"Flags",
		"MinGearLevelSuggestion",
		"PlayerConditionID",
		"MapID",
		"DifficultyID",
		"ExpansionID",
		"AreaID",
		"MaxPlayers",
		"DisplayType",
		"OverrideContentTuningID",
		"MapChallengeModeID",
	)

	var entries []GroupFinderActivity

	for _, record := range records {
		entry := GroupFinderActivity{
			ID:                       atoi(record["ID"]),
			FullName_lang:            record["FullName_lang"],
			ShortName_lang:           record["ShortName_lang"],
			GroupFinderCategoryID:    atoi(record["GroupFinderCategoryID"]),
			OrderIndex:               atoi(record["OrderIndex"]),
			GroupFinderActivityGrpID: atoi(record["GroupFinderActivityGrpID"]),
			Flags:                    atoi(record["Flags"]),
			MinGearLevelSuggestion:   atoi(record["MinGearLevelSuggestion"]),
			PlayerConditionID:        atoi(record["PlayerConditionID"]),
			MapID:                    atoi(record["MapID"]),
			DifficultyID:             atoi(record["DifficultyID"]),
			ExpansionID:              atoi(record["ExpansionID"]),
			AreaID:                   atoi(record["AreaID"]),
			MaxPlayers:               atoi(record["MaxPlayers"]),              // 0-40
			DisplayType:              atoi(record["DisplayType"]),             // 0 (Role Counts), 1 (Role Enumeration), 2 (Class Enumeration), 4 (Player Count)
			OverrideContentTuningID:  atoi(record["OverrideContentTuningID"]), // only set for Atal'Dazar
			MapChallengeModeID:       atoi(record["MapChallengeModeID"]),      // only set for a few dungeons
		}

		entries = append(entries, entry)
	}

	return entries
}

func parseMapChallengeMode(source io.Reader) []MapChallengeMode {
	records := readCSV(source,
		"Name_lang",
		"ID",
		"MapID",
		"Flags",
		"Field_12_0_0_63854_004",
		"ExpansionLevel",
	)

	var entries []MapChallengeMode

	for _, record := range records {
		entry := MapChallengeMode{
			Name_lang:              record["Name_lang"],
			ID:                     atoi(record["ID"]),
			MapID:                  atoi(record["MapID"]),
			Flags:                  atoi(record["Flags"]),
			Field_12_0_0_63854_004: atoi(record["Field_12_0_0_63854_004"]),
			ExpansionLevel:         atoi(record["ExpansionLevel"]),
		}

		entries = append(entries, entry)
	}

	return entries
}

func parseMythicPlusSeasonTrackedMap(source io.Reader) []MythicPlusSeasonTrackedMap {
	records := readCSV(source,
		"ID",
		"MapChallengeModeID",
		"DisplaySeasonID",
	)

	var entries []MythicPlusSeasonTrackedMap

	for _, record := range records {
		entry := MythicPlusSeasonTrackedMap{
			ID:                 atoi(record["ID"]),
			MapChallengeModeID: atoi(record["MapChallengeModeID"]),
			DisplaySeasonID:    atoi(record["DisplaySeasonID"]),
		}

		entries = append(entries, entry)
	}

	return entries
}

// readCSV resolves fields by name because column order differs between builds.
func readCSV(source io.Reader, fields ...string) []map[string]string {
	csvReader := csv.NewReader(source)
	records, err := csvReader.ReadAll()
	if err != nil {
		logger.Fatalf("failed to read CSV records: %v", err)
	}
	if len(records) == 0 {
		return nil
	}

	columns := make(map[string]int, len(records[0]))
	for i, name := range records[0] {
		columns[name] = i
	}
	for _, field := range fields {
		if _, ok := columns[field]; !ok {
			panic(fmt.Errorf("missing CSV column %q", field))
		}
	}

	rows := make([]map[string]string, 0, len(records)-1)
	for _, record := range records[1:] {
		row := make(map[string]string, len(fields))
		for _, field := range fields {
			row[field] = record[columns[field]]
		}
		rows = append(rows, row)
	}
	return rows
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
