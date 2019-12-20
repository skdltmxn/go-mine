package server

const (
	GameModeSurvival = iota
	GameModeCreative
	GameModeAdventure
	GameModeSpectator
)

const (
	GameDimensionNether = iota - 1
	GameDimensionOverworld
	GameDimensionEnd
)

const (
	GameLevelDefault     = "default"
	GameLevelFlat        = "flat"
	GameLevelLargeBiomes = "largeBiomes"
	GameLevelAmplified   = "amplified"
	GameLevelCustomized  = "customized"
	GameLevelBuffet      = "buffet"
	GameLevelDefault_1_1 = "default_1_1"
)
