package main

import (
	"0xbs/premade-groups-filter-helper/logger"
	"fmt"
	"os"
)

type GroupFinderActivity struct {
	ID                       int
	FullName_lang            string
	ShortName_lang           string
	GroupFinderCategoryID    int
	OrderIndex               int
	GroupFinderActivityGrpID int
	Flags                    int
	MinGearLevelSuggestion   int
	PlayerConditionID        int
	MapID                    int
	DifficultyID             int
	AreaID                   int
	MaxPlayers               int
	DisplayType              int
	OverrideContentTuningID  int
	MapChallengeModeID       int
}

type MapChallengeMode struct {
	Name_lang      string
	ID             int
	MapID          int
	Flags          int
	ExpansionLevel int
}

// maps DifficultyID to PGF difficulty (0 = ignore, 1 = normal, 2 = heroic, 3 = mythic, 4 = mythic+)
var difficultyMap = map[int]int{
	1:  1, // DungeonNormal
	2:  2, // DungeonHeroic
	3:  1, // Raid10Normal
	4:  1, // Raid25Normal
	5:  2, // Raid10Heroic
	6:  2, // Raid25Heroic
	7:  0, // RaidLFR
	8:  4, // DungeonChallenge
	9:  1, // Raid40
	14: 1, // PrimaryRaidNormal
	15: 2, // PrimaryRaidHeroic
	16: 3, // PrimaryRaidMythic
	17: 0, // PrimaryRaidLFR
	23: 3, // DungeonMythic
	24: 0, // DungeonTimewalker
	33: 0, // RaidTimewalker
}

//goland:noinspection GoUnhandledErrorResult
func main() {
	gfa, err := os.Open("data/GroupFinderActivity.csv")
	if err != nil {
		panic(err)
	}
	defer gfa.Close()
	activities := parseGroupFinderActivity(gfa)

	mcm, err := os.Open("data/MapChallengeMode.csv")
	if err != nil {
		panic(err)
	}
	defer mcm.Close()
	challenges := parseMapChallengeMode(mcm)

	mapID2cmID := make(map[int]int)
	for _, challenge := range challenges {
		mapID2cmID[challenge.MapID] = challenge.ID
	}

	raids, err := os.Create("data/Raids.lua")
	if err != nil {
		panic(err)
	}
	defer raids.Close()

	fmt.Fprint(raids, "C.ACTIVITY_TO_MAP_ID = {\n")
	for _, activity := range activities {
		if activity.ID >= 1189 && (activity.DifficultyID == 14 || activity.DifficultyID == 15 || activity.DifficultyID == 16) {
			fmt.Fprintf(raids, "    [%4d] = %4d, -- %s\n", activity.ID, activity.MapID, activity.FullName_lang)
		}
	}
	fmt.Fprint(raids, "}\n")

	act, err := os.Create("data/Activity.lua")
	if err != nil {
		panic(err)
	}
	defer act.Close()

	fmt.Fprint(act, "C.ACTIVITY = {\n")
	for _, activity := range activities {
		if difficultyMap[activity.DifficultyID] == 0 {
			continue
		}
		cmID := activity.MapChallengeModeID
		if activity.DifficultyID == 8 {
			mapChallengeModeID := mapID2cmID[activity.MapID]
			if activity.MapChallengeModeID == 0 {
				logger.Infof("set missing cmID %d for activity %d %s",
					mapChallengeModeID, activity.ID, activity.FullName_lang)
				cmID = mapChallengeModeID
			} else if mapChallengeModeID != activity.MapChallengeModeID {
				logger.Warnf("different cmIDs %d (MapChallengeMode) and %d (GroupFinderActivity) for activity %d %s",
					mapChallengeModeID, activity.MapChallengeModeID, activity.ID, activity.FullName_lang)
			}
		}
		fmt.Fprintf(act, "    [%4d] = { d = %d, mapID = %4d, cmID = %3d, cat = %d }, -- %s\n",
			activity.ID,
			difficultyMap[activity.DifficultyID],
			activity.MapID,
			cmID,
			activity.GroupFinderCategoryID,
			activity.FullName_lang,
		)
	}
	fmt.Fprint(act, "}\n")
}
