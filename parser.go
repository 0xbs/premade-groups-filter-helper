package main

import (
	"0xbs/premade-groups-filter-helper/logger"
	"encoding/csv"
	"io"
	"strconv"
)

func parseGroupFinderActivity(source io.Reader) []GroupFinderActivity {
	csvReader := csv.NewReader(source)

	records, err := csvReader.ReadAll()
	if err != nil {
		logger.Fatalf("failed to read CSV records: %v", err)
	}

	var entries []GroupFinderActivity

	for i, record := range records {
		if i == 0 {
			continue // ignore header
		}

		entry := GroupFinderActivity{
			ID:                       atoi(record[0]),
			FullName_lang:            record[1],
			ShortName_lang:           record[2],
			GroupFinderCategoryID:    atoi(record[3]),
			OrderIndex:               atoi(record[4]),
			GroupFinderActivityGrpID: atoi(record[5]),
			Flags:                    atoi(record[6]),
			MinGearLevelSuggestion:   atoi(record[7]),
			PlayerConditionID:        atoi(record[8]),
			MapID:                    atoi(record[9]),
			DifficultyID:             atoi(record[10]),
			AreaID:                   atoi(record[11]),
			MaxPlayers:               atoi(record[12]), // 0-40
			DisplayType:              atoi(record[13]), // 0 (Role Counts), 1 (Role Enumeration), 2 (Class Enumeration), 4 (Player Count)
			// unknown fields inbetween
			OverrideContentTuningID: atoi(record[18]), // only set for Atal'Dazar
			MapChallengeModeID:      atoi(record[19]), // only set for a few dungeons
		}

		entries = append(entries, entry)
	}

	return entries
}

func parseMapChallengeMode(source io.Reader) []MapChallengeMode {
	csvReader := csv.NewReader(source)

	records, err := csvReader.ReadAll()
	if err != nil {
		logger.Fatalf("failed to read CSV records: %v", err)
	}

	var entries []MapChallengeMode

	for i, record := range records {
		if i == 0 {
			continue // ignore header
		}

		entry := MapChallengeMode{
			Name_lang:      record[0],
			ID:             atoi(record[1]),
			MapID:          atoi(record[2]),
			Flags:          atoi(record[3]),
			ExpansionLevel: atoi(record[4]),
			// unknown fields following
		}

		entries = append(entries, entry)
	}

	return entries
}

func parseMythicPlusSeasonTrackedMap(source io.Reader) []MythicPlusSeasonTrackedMap {
	csvReader := csv.NewReader(source)

	records, err := csvReader.ReadAll()
	if err != nil {
		logger.Fatalf("failed to read CSV records: %v", err)
	}

	var entries []MythicPlusSeasonTrackedMap

	for i, record := range records {
		if i == 0 {
			continue // ignore header
		}

		entry := MythicPlusSeasonTrackedMap{
			ID:                 atoi(record[0]),
			MapChallengeModeID: atoi(record[1]),
			DisplaySeasonID:    atoi(record[2]),
		}

		entries = append(entries, entry)
	}

	return entries
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
