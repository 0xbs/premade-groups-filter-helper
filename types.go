package main

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
	ExpansionID              int
	AreaID                   int
	MaxPlayers               int
	DisplayType              int
	OverrideContentTuningID  int // diff
	MapChallengeModeID       int // diff
}

type MapChallengeMode struct {
	Name_lang              string
	ID                     int
	MapID                  int
	Flags                  int
	Field_12_0_0_63854_004 int
	ExpansionLevel         int
	RequiredWorldStateID   int
}

type MythicPlusSeasonTrackedMap struct {
	ID                 int
	MapChallengeModeID int
	DisplaySeasonID    int
}
