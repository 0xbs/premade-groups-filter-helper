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
	AreaID                   int
	MaxPlayers               int
	DisplayType              int
	OverrideContentTuningID  int // diff
	MapChallengeModeID       int // diff
}

type GroupFinderActivityWrath struct {
	ID                       int
	FullName_lang            string
	ShortName_lang           string
	GroupFinderCategoryID    int
	OrderIndex               int
	GroupFinderActivityGrpID int
	Field_3_4_0_43659_005    int // diff
	Flags                    int
	MinGearLevelSuggestion   int
	PlayerConditionID        int
	MapID                    int
	DifficultyID             int
	AreaID                   int
	MaxPlayers               int
	DisplayType              int
	MinLevel                 int // diff
	MaxLevelSuggestion       int // diff
	IconFileDataID           int // diff
}

type MapChallengeMode struct {
	Name_lang      string
	ID             int
	MapID          int
	Flags          int
	ExpansionLevel int
}
