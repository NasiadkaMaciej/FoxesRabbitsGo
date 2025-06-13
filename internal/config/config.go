package config

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Config struct {
	// World configuration
	WorldWidth                     int
	WorldHeight                    int
	FrameTime                      time.Duration
	InitialFoxes                   int
	InitialRabbits                 int
	InitialGrass                   int
	AnimalSize                     int
	ChanceForNotMoveWhenMovingAway float64

	// Fox parameters
	FoxInitialEnergy        int
	FoxEnergyLossPerMove    int
	FoxEnergyGainFromRabbit int
	FoxReproductionCost     int
	FoxColor                sdl.Color
	FoxReproductionRange    int
	FoxEatingRange          int
	FoxFollowRabbitRange    int
	FoxEatingCooldown       int
	FoxReproductionCooldown int

	// Rabbit parameters
	RabbitInitialEnergy        int
	RabbitEnergyLossPerMove    int
	RabbitEnergyGainFromGrass  int
	RabbitReproductionCost     int
	RabbitColor                sdl.Color
	RabbitReproductionRange    int
	RabbitEscapeRange          int
	RabbitEatingCooldown       int
	RabbitReproductionCooldown int

	// Grass parameters
	GrassGrowthRate    int
	GrassMaxAmount     int
	GrassRegrowthTimer int
	GrassBaseColor     sdl.Color
}

func NewConfig() *Config {
	return &Config{
		// World configuration
		WorldWidth:                     120,
		WorldHeight:                    80,
		FrameTime:                      100,
		InitialFoxes:                   10,
		InitialRabbits:                 40,
		InitialGrass:                   2,
		AnimalSize:                     8,
		ChanceForNotMoveWhenMovingAway: 0.2,

		// Fox parameters
		FoxInitialEnergy:        100,
		FoxEnergyLossPerMove:    2,
		FoxEnergyGainFromRabbit: 90,
		FoxReproductionCost:     200,
		FoxColor:                sdl.Color{R: 255, G: 0, B: 0, A: 255},
		FoxReproductionRange:    2,
		FoxEatingRange:          2,
		FoxFollowRabbitRange:    30,
		FoxEatingCooldown:       5,
		FoxReproductionCooldown: 15,

		// Rabbit parameters
		RabbitInitialEnergy:        15,
		RabbitEnergyLossPerMove:    1,
		RabbitEnergyGainFromGrass:  3,
		RabbitReproductionCost:     30,
		RabbitColor:                sdl.Color{R: 0, G: 0, B: 255, A: 255},
		RabbitReproductionRange:    1,
		RabbitEscapeRange:          10,
		RabbitEatingCooldown:       2,
		RabbitReproductionCooldown: 5,

		// Grass parameters
		GrassGrowthRate:    3,
		GrassMaxAmount:     2,
		GrassRegrowthTimer: 50,
		GrassBaseColor:     sdl.Color{R: 0, G: 100, B: 0, A: 255},
	}
}
