package main

// maps DifficultyID to PGF difficulty (0 = ignore, 1 = normal, 2 = heroic, 3 = mythic, 4 = mythic+)
var difficultyMap = map[int]int{
	1:   1, // DungeonNormal
	2:   2, // DungeonHeroic
	3:   1, // Raid10Normal
	4:   1, // Raid25Normal
	5:   2, // Raid10Heroic
	6:   2, // Raid25Heroic
	7:   0, // RaidLFR
	8:   4, // DungeonChallenge
	9:   1, // Raid40
	14:  1, // PrimaryRaidNormal
	15:  2, // PrimaryRaidHeroic
	16:  3, // PrimaryRaidMythic
	17:  0, // PrimaryRaidLFR
	23:  3, // DungeonMythic
	24:  0, // DungeonTimewalker
	33:  0, // RaidTimewalker
	167: 0, // Torghast
}

// contains split dungeons that have the same mapID
var fixedChallengeModeIDs = map[int]int{
	471:  227, // Return to Karazhan: Lower
	473:  234, // Return to Karazhan: Upper
	679:  369, // Operation: Mechagon - Junkyard
	683:  370, // Operation: Mechagon - Workshop
	1016: 391, // Tazavesh: Streets of Wonder
	1017: 392, // Tazavesh: So'leah's Gambit
	1247: 463, // Dawn of the Infinite: Galakrond's Fall
	1248: 464, // Dawn of the Infinite: Murozond's Rise
	1250: 463, // Galakrond's Fall (Mythic Keystone) // no idea why Fall exists twice
}
