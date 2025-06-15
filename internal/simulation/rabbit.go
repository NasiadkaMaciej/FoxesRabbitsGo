package simulation

import (
	"foxes-rabbits-simulation/internal/config"
)

type Rabbit struct {
	AnimalBase
}

func NewRabbit(x, y int, cfg *config.Config) *Rabbit {
	return &Rabbit{
		AnimalBase: AnimalBase{
			Position: Position{x, y},
			Energy:   cfg.RabbitInitialEnergy,
			Config:   cfg,
		},
	}
}

func (r *Rabbit) Move(world *World) {
	nearestFox, foundFox := FindNearestAnimal(r, world.Foxes, r.Config.RabbitEscapeRange)

	// If found a fox within range, try to move away from it
	if foundFox && r.MoveAwayFrom(nearestFox.GetPosition(), world) {
		// Successfully moved away from fox
	} else {
		// No nearby fox or couldn't move away, just move randomly
		r.MoveRandomly(world)
	}

	r.Energy -= r.Config.RabbitEnergyLossPerMove
}

func (r *Rabbit) Eat(grass *Grass) {
	r.TurnsSinceEaten++

	if !r.CanEat(r.Config.RabbitEatingCooldown) {
		return
	}

	if grass.Amount > 0 {
		grass.Eat(1)
		r.Energy += r.Config.RabbitEnergyGainFromGrass
		r.TurnsSinceEaten = 0
	}
}

func (r *Rabbit) Reproduce(world *World) *Rabbit {
	r.TurnsSinceReproduction++

	if !r.CanReproduce(r.Config.RabbitReproductionCooldown) {
		return nil
	}

	if r.Energy >= r.Config.RabbitReproductionCost && r.hasNearbyRabbit(world) {
		r.Energy -= r.Config.RabbitReproductionCost
		r.TurnsSinceReproduction = 0

		if newX, newY, found := FindEmptyAdjacentPosition(r.Position, world, 8); found {
			return NewRabbit(newX, newY, r.Config)
		}
	}
	return nil
}

func (r *Rabbit) hasNearbyRabbit(world *World) bool {
	return IsNearbyAnimal(r, world.Rabbits, r.Config.RabbitReproductionRange)
}
