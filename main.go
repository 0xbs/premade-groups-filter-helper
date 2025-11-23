package main

import (
	"0xbs/premade-groups-filter-helper/logger"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
)

const (
	runDownload   = true
	classicBuild  = "5.5.3.64530"
	classicSuffix = "_Mists"
	headerPath    = "data/Header.lua"
)

func main() {
	if runDownload {
		download("https://wago.tools/db2/GroupFinderActivity/csv", "data/GroupFinderActivity.csv")
		download("https://wago.tools/db2/MapChallengeMode/csv", "data/MapChallengeMode.csv")
		download("https://wago.tools/db2/MythicPlusSeasonTrackedMap/csv", "data/MythicPlusSeasonTrackedMap.csv")
		download("https://wago.tools/db2/GroupFinderActivity/csv?build="+classicBuild, "data/GroupFinderActivity"+classicSuffix+".csv")
		download("https://wago.tools/db2/MapChallengeMode/csv?build="+classicBuild, "data/MapChallengeMode"+classicSuffix+".csv")
		download("https://wago.tools/db2/MythicPlusSeasonTrackedMap/csv?build="+classicBuild, "data/MythicPlusSeasonTrackedMap"+classicSuffix+".csv")
	}

	process("")
	process(classicSuffix)
	logger.Infof("✅ all done")
}

//goland:noinspection GoUnhandledErrorResult
func download(url, dest string) {
	logger.Infof("⬇️ downloading %s to %s", url, dest)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}

//goland:noinspection GoUnhandledErrorResult
func process(suffix string) {
	logger.Infof("➡️ processing activities with suffix '%s'", suffix)

	gfa, err := os.Open("data/GroupFinderActivity" + suffix + ".csv")
	if err != nil {
		panic(err)
	}
	defer gfa.Close()
	activities := parseGroupFinderActivity(gfa)

	mcm, err := os.Open("data/MapChallengeMode" + suffix + ".csv")
	if err != nil {
		panic(err)
	}
	defer mcm.Close()
	challenges := parseMapChallengeMode(mcm)

	// look-up map for the season processing
	cmID2Activity := make(map[int]GroupFinderActivity)

	mapID2cmID := make(map[int]int)
	for _, challenge := range challenges {
		mapID2cmID[challenge.MapID] = challenge.ID
	}

	fmtActivities := make([]string, len(activities))

	for _, activity := range activities {
		if difficultyMap[activity.DifficultyID] == 0 {
			//logger.Infof("unknown difficultyID %d", activity.DifficultyID)
			continue
		}
		// check for missing mapID
		if activity.MapID == 0 {
			fixedMapID := activity2MapID[activity.ID]
			if fixedMapID != 0 {
				activity.MapID = fixedMapID
				logger.Infof("set fixed mapID %d for activity %d %s",
					fixedMapID, activity.ID, activity.FullName_lang)
			}
		}
		// handle challenge modes
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

			cmID2Activity[cmID] = activity
		}
		fmtActivity := fmt.Sprintf("    [%4d] = { difficulty = %d, category = %3d, mapID = %4d, cmID = %3d }, -- %s\n",
			activity.ID,
			difficultyMap[activity.DifficultyID],
			activity.GroupFinderCategoryID,
			activity.MapID,
			cmID,
			activity.FullName_lang,
		)
		fmtActivities = append(fmtActivities, fmtActivity)
	}

	sort.Strings(fmtActivities)

	header, err := os.ReadFile(headerPath)
	if err != nil {
		panic(err)
	}

	// write the complete file at the end (because we have to sort first)
	actFile, err := os.Create("data/Activity" + suffix + ".lua")
	if err != nil {
		panic(err)
	}
	defer actFile.Close()
	actFile.Write(header)
	actFile.WriteString("C.ACTIVITY = {\n")
	actFile.WriteString(strings.Join(fmtActivities, ""))
	actFile.WriteString("}\n\n")
	actFile.WriteString("-- Return a default set if activity not found\n")
	actFile.WriteString("setmetatable(C.ACTIVITY, { __index = function() return { difficulty = 0, category = 0, mapID = 0, cmID = 0 } end })\n")

	logger.Infof("✅ processing activities done")
	processMythicPlusSeasons(cmID2Activity, suffix)
}

//goland:noinspection GoUnhandledErrorResult
func processMythicPlusSeasons(cmID2Activity map[int]GroupFinderActivity, suffix string) {
	logger.Infof("➡️ processing mythic plus seasons with suffix '%s'", suffix)
	mps, err := os.Open("data/MythicPlusSeasonTrackedMap" + suffix + ".csv")
	if err != nil {
		panic(err)
	}
	defer mps.Close()
	seasonTable := parseMythicPlusSeasonTrackedMap(mps)

	// group by seasons first: map DisplaySeasonID -> []MapChallengeModeID
	seasonMap := make(map[int][]int)
	for _, s := range seasonTable {
		seasonMap[s.DisplaySeasonID] = append(seasonMap[s.DisplaySeasonID], s.MapChallengeModeID)
	}

	// sort by DisplaySeasonID desc
	seasonMapKeys := make([]int, 0)
	for k := range seasonMap {
		seasonMapKeys = append(seasonMapKeys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(seasonMapKeys)))

	// open file for writing
	csFile, err := os.Create("data/MythicPlusSeasons" + suffix + ".lua")
	if err != nil {
		panic(err)
	}
	defer csFile.Close()
	csFile.WriteString("local CMID_MAP = {\n")

	// loop through seasons and their dungeons
	for _, displaySeasonID := range seasonMapKeys {
		// sort by MapChallengeModeID desc // TODO would name be better?
		seasonChallengeModeIDs := seasonMap[displaySeasonID]
		sort.Sort(sort.Reverse(sort.IntSlice(seasonChallengeModeIDs)))

		fmt.Fprintf(csFile, "\n    -- Season #%d\n", displaySeasonID)
		for i, cmID := range seasonChallengeModeIDs {
			activity := cmID2Activity[cmID]
			fmt.Fprintf(csFile, "    [%3d] = { order = %d, activityGroupID = %3d, keyword = \"\" }, -- %s\n",
				cmID, i+1, activity.GroupFinderActivityGrpID, activity.FullName_lang)
		}
	}

	csFile.WriteString("}\n")
	logger.Infof("✅ processing mythic plus seasons done")
}
