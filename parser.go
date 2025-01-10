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
			MaxPlayers:               atoi(record[12]),
			DisplayType:              atoi(record[13]),
			OverrideContentTuningID:  atoi(record[14]),
			MapChallengeModeID:       atoi(record[15]),
		}

		entries = append(entries, entry)
	}

	return entries
}

func parseGroupFinderActivityWrath(source io.Reader) []GroupFinderActivityWrath {
	csvReader := csv.NewReader(source)

	records, err := csvReader.ReadAll()
	if err != nil {
		logger.Fatalf("failed to read CSV records: %v", err)
	}

	var entries []GroupFinderActivityWrath

	for i, record := range records {
		if i == 0 {
			continue // ignore header
		}

		entry := GroupFinderActivityWrath{
			ID:                       atoi(record[0]),
			FullName_lang:            record[1],
			ShortName_lang:           record[2],
			GroupFinderCategoryID:    atoi(record[3]),
			OrderIndex:               atoi(record[4]),
			GroupFinderActivityGrpID: atoi(record[5]),
			//Field_3_4_0_43659_005:    atoi(record[6]),
			Flags:                    atoi(record[6]),
			MinGearLevelSuggestion:   atoi(record[7]),
			PlayerConditionID:        atoi(record[8]),
			MapID:                    atoi(record[9]),
			DifficultyID:             atoi(record[10]),
			AreaID:                   atoi(record[11]),
			MaxPlayers:               atoi(record[12]),
			DisplayType:              atoi(record[13]),
			MinLevel:                 atoi(record[14]),
			MaxLevelSuggestion:       atoi(record[15]),
			IconFileDataID:           atoi(record[16]),
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
