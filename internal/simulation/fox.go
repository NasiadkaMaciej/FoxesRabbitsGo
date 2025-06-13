package simulation

import (
	"foxes-rabbits-simulation/internal/config"
)

type Fox struct {
	AnimalBase
}

func NewFox(x, y int, cfg *config.Config) *Fox {
	return &Fox{
		AnimalBase: AnimalBase{
			Energy:   cfg.FoxInitialEnergy,
			Position: Position{X: x, Y: y},
			Config:   cfg,
		},
	}
}

func (f *Fox) Move(world *World) {
	// Try to find nearby rabbit
	nearest, found := FindNearestAnimal(f, world.Rabbits, f.Config.FoxFollowRabbitRange)

	// If found a rabbit within range, try to move toward it
	if found && f.MoveToward(nearest.GetPosition(), world) {
		// Successfully moved toward rabbit
	} else {
		// No nearby rabbit or couldn't move toward it, move randomly
		f.MoveRandomly(world)
	}

	// Energy decreases regardless of movement
	f.Energy -= f.Config.FoxEnergyLossPerMove
}

func (f *Fox) Eat(world *World) {
	// Increment turns since last eaten
	f.TurnsSinceEaten++

	// Check if fox can eat based on cooldown
	if !f.CanEat(f.Config.FoxEatingCooldown) {
		return
	}

	nearestRabbit, found := FindNearestAnimal(f, world.Rabbits, f.Config.FoxEatingRange)

	if found {
		f.Energy += f.Config.FoxEnergyGainFromRabbit
		f.TurnsSinceEaten = 0

		// Remove the eaten rabbit - filter it out of the slice
		newRabbits := make([]*Rabbit, 0, len(world.Rabbits)-1)
		for _, rabbit := range world.Rabbits {
			if rabbit != nearestRabbit {
				newRabbits = append(newRabbits, rabbit)
			}
		}
		world.Rabbits = newRabbits
	}
}

func (f *Fox) Reproduce(world *World) *Fox {
	// Increment turns since last reproduction
	f.TurnsSinceReproduction++

	// Check if fox can reproduce based on cooldown
	if !f.CanReproduce(f.Config.FoxReproductionCooldown) {
		return nil
	}

	if f.Energy >= f.Config.FoxReproductionCost && f.hasNearbyFox(world) {
		f.Energy -= f.Config.FoxReproductionCost
		f.TurnsSinceReproduction = 0

		// Find an empty adjacent position for the new fox
		if newX, newY, found := FindEmptyAdjacentPosition(f.Position, world, 8); found {
			return NewFox(newX, newY, f.Config)
		}
	}
	return nil
}

func (f *Fox) hasNearbyFox(world *World) bool {
	return IsNearbyAnimal(f, world.Foxes, f.Config.FoxReproductionRange)
}
