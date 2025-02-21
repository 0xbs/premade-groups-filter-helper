package main

import (
	"0xbs/premade-groups-filter-helper/logger"
	"fmt"
	"os"
	"sort"
	"strings"
)

//goland:noinspection GoUnhandledErrorResult
func main() {
	writeWrath()

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

	writeRaids(activities)

	fmtActivities := []string{}
	for _, activity := range activities {
		if difficultyMap[activity.DifficultyID] == 0 {
			//logger.Infof("unknown difficultyID %d", activity.DifficultyID)
			continue
		}
		cmID := activity.MapChallengeModeID
		if activity.DifficultyID == 8 {
			mapChallengeModeID := mapID2cmID[activity.MapID]
			fixedChallengeModeID := fixedChallengeModeIDs[activity.ID]
			if fixedChallengeModeID != 0 {
				logger.Infof("set fixed cmID %d for activity %d %s",
					fixedChallengeModeID, activity.ID, activity.FullName_lang)
				cmID = fixedChallengeModeID
			} else if activity.MapChallengeModeID == 0 {
				logger.Infof("set missing cmID %d for activity %d %s",
					mapChallengeModeID, activity.ID, activity.FullName_lang)
				cmID = mapChallengeModeID
			} else if mapChallengeModeID != activity.MapChallengeModeID {
				logger.Warnf("different cmIDs %d (MapChallengeMode) and %d (GroupFinderActivity) for activity %d %s",
					mapChallengeModeID, activity.MapChallengeModeID, activity.ID, activity.FullName_lang)
			}
		}
		fmtActivity := fmt.Sprintf("    [%4d] = { difficulty = %d, category = %d, mapID = %4d, cmID = %3d }, -- %s\n",
			activity.ID,
			difficultyMap[activity.DifficultyID],
			activity.GroupFinderCategoryID,
			activity.MapID,
			cmID,
			activity.FullName_lang,
		)
		fmtActivities = append(fmtActivities, fmtActivity)
	}
	sort.Sort(sort.StringSlice(fmtActivities))

	act, err := os.Create("data/Activity.lua")
	if err != nil {
		panic(err)
	}
	defer act.Close()
	fmt.Fprint(act, "C.ACTIVITY = {\n")
	fmt.Fprint(act, strings.Join(fmtActivities, ""))
	fmt.Fprint(act, "}\n\n")
	fmt.Fprint(act, "-- Return a default set if activity not found\n")
	fmt.Fprint(act, "setmetatable(C.ACTIVITY, { __index = function() return { difficulty = 0, category = 0, mapID = 0, cmID = 0 } end })\n")
}

func writeWrath() {
	writeLegacy("data/GroupFinderActivity_Wrath.csv", "data/Activity_Wrath.lua")
}

func writeLegacy(source, target string) {
	gfaWrath, err := os.Open(source)
	if err != nil {
		panic(err)
	}
	defer gfaWrath.Close()
	activitiesWrath := parseGroupFinderActivityWrath(gfaWrath)

	fmtActivities := []string{}
	for _, activity := range activitiesWrath {
		if difficultyMap[activity.DifficultyID] == 0 {
			continue
		}
		fmtActivity := fmt.Sprintf("    [%4d] = { difficulty = %d, category = %3d, mapID = %4d }, -- %s\n",
			activity.ID,
			difficultyMap[activity.DifficultyID],
			activity.GroupFinderCategoryID,
			activity.MapID,
			activity.FullName_lang,
		)
		fmtActivities = append(fmtActivities, fmtActivity)
	}
	sort.Sort(sort.StringSlice(fmtActivities))

	actWrath, err := os.Create(target)
	if err != nil {
		panic(err)
	}
	defer actWrath.Close()
	fmt.Fprint(actWrath, "C.ACTIVITY = {\n")
	fmt.Fprint(actWrath, strings.Join(fmtActivities, ""))
	fmt.Fprint(actWrath, "}\n\n")
	fmt.Fprint(actWrath, "-- Return a default set if activity not found\n")
	fmt.Fprint(actWrath, "setmetatable(C.ACTIVITY, { __index = function() return { difficulty = 0, category = 0, mapID = 0 } end })\n")
}

func writeRaids(activities []GroupFinderActivity) {
	fmtActivities := []string{}
	for _, activity := range activities {
		if activity.ID >= 1189 && (activity.DifficultyID == 14 || activity.DifficultyID == 15 || activity.DifficultyID == 16) {
			fmtActivity := fmt.Sprintf("    [%4d] = %4d, -- %s\n", activity.ID, activity.MapID, activity.FullName_lang)
			fmtActivities = append(fmtActivities, fmtActivity)
		}
	}
	sort.Sort(sort.StringSlice(fmtActivities))

	raids, err := os.Create("data/Raids.lua")
	if err != nil {
		panic(err)
	}
	defer raids.Close()
	fmt.Fprint(raids, "C.ACTIVITY_TO_MAP_ID = {\n")
	fmt.Fprint(raids, strings.Join(fmtActivities, ""))
	fmt.Fprint(raids, "}\n")
}
