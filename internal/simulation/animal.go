package simulation

import (
	"foxes-rabbits-simulation/internal/config"
	"math/rand"
)

// Animal defines the interface for all animals in the simulation
type Animal interface {
	IsDead() bool
	GetPosition() Position
	GetConfig() *config.Config
	GetEnergy() int
	CanEat(cooldown int) bool
}

// AnimalBase contains common fields and functionality for all animals
type AnimalBase struct {
	Position               Position
	Energy                 int
	Config                 *config.Config
	TurnsSinceEaten        int
	TurnsSinceReproduction int
}

// IsDead returns true if the animal's energy is <= 0
func (a *AnimalBase) IsDead() bool {
	return a.Energy <= 0
}

// GetPosition returns the animal's position
func (a *AnimalBase) GetPosition() Position {
	return a.Position
}

// GetConfig returns the animal's configuration
func (a *AnimalBase) GetConfig() *config.Config {
	return a.Config
}

// GetEnergy returns the animal's current energy
func (a *AnimalBase) GetEnergy() int {
	return a.Energy
}

// FindEmptyAdjacentPosition finds an empty position adjacent to the animal
func FindEmptyAdjacentPosition(pos Position, world *World, maxAttempts int) (int, int, bool) {
	for attempts := 0; attempts < maxAttempts; attempts++ {
		dx := rand.Intn(3) - 1 // -1, 0, or 1
		dy := rand.Intn(3) - 1 // -1, 0, or 1

		newX := pos.X + dx
		newY := pos.Y + dy

		// Check if the position is valid and empty
		if newX >= 0 && newX < world.Width &&
			newY >= 0 && newY < world.Height &&
			!world.IsPositionOccupied(newX, newY) {
			return newX, newY, true
		}
	}
	return 0, 0, false
}

// IsNearbyAnimal checks if there's a nearby animal of the same type
func IsNearbyAnimal[T Animal](animal T, others []T, range_ int) bool {
	for _, other := range others {
		// Skip checking against itself
		if animal.GetPosition() == other.GetPosition() {
			continue
		}

		// Calculate distance between animals
		dx := abs(animal.GetPosition().X - other.GetPosition().X)
		dy := abs(animal.GetPosition().Y - other.GetPosition().Y)

		// Check if animal is within the configured reproduction range
		if dx <= range_ && dy <= range_ {
			return true
		}
	}
	return false
}

// CanEat returns true if enough turns have passed since last eating
func (a *AnimalBase) CanEat(cooldown int) bool {
	return a.TurnsSinceEaten >= cooldown // We'll check specific animal type in Eat methods
}

// CanReproduce returns true if enough turns have passed since last reproduction
func (a *AnimalBase) CanReproduce(cooldown int) bool {
	return a.TurnsSinceReproduction >= cooldown
}

// Helper function for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func FindNearestAnimal[T Animal](from Animal, animals []T, maxRange int) (T, bool) {
	var nearest T
	foundAnimal := false
	minDistance := maxRange * 2 // Initialize with a value larger than maxRange

	for _, animal := range animals {
		dx := abs(from.GetPosition().X - animal.GetPosition().X)
		dy := abs(from.GetPosition().Y - animal.GetPosition().Y)

		// Use Manhattan distance
		distance := dx + dy

		if distance < minDistance && distance <= maxRange {
			minDistance = distance
			nearest = animal
			foundAnimal = true
		}
	}

	return nearest, foundAnimal
}

func (a *AnimalBase) MoveRandomly(world *World) {
	// Try up to 4 times to find a valid move
	for attempts := 0; attempts < 4; attempts++ {
		newX, newY := a.Position.X, a.Position.Y
		direction := rand.Intn(4)

		switch direction {
		case 0:
			if newX < world.Width-1 {
				newX++ // Move right
			}
		case 1:
			if newX > 0 {
				newX-- // Move left
			}
		case 2:
			if newY < world.Height-1 {
				newY++ // Move down
			}
		case 3:
			if newY > 0 {
				newY-- // Move up
			}
		}

		// If position is not occupied, move there
		if !world.IsPositionOccupied(newX, newY) || (newX == a.Position.X && newY == a.Position.Y) {
			a.Position.X, a.Position.Y = newX, newY
			break
		}
	}
}

func (a *AnimalBase) MoveToward(targetPos Position, world *World) bool {
	dx := targetPos.X - a.Position.X
	dy := targetPos.Y - a.Position.Y

	// Try to move horizontally first if dx is larger
	if abs(dx) >= abs(dy) {
		if dx > 0 && a.Position.X < world.Width-1 && !world.IsPositionOccupied(a.Position.X+1, a.Position.Y) {
			a.Position.X++
			return true
		} else if dx < 0 && a.Position.X > 0 && !world.IsPositionOccupied(a.Position.X-1, a.Position.Y) {
			a.Position.X--
			return true
		}
	}

	// Try to move vertically if horizontal movement not possible
	if dy > 0 && a.Position.Y < world.Height-1 && !world.IsPositionOccupied(a.Position.X, a.Position.Y+1) {
		a.Position.Y++
		return true
	} else if dy < 0 && a.Position.Y > 0 && !world.IsPositionOccupied(a.Position.X, a.Position.Y-1) {
		a.Position.Y--
		return true
	}

	return false
}

func (a *AnimalBase) MoveAwayFrom(targetPos Position, world *World) bool {
	// 1 in 5 chance of not moving
	if rand.Float64() < a.Config.ChanceForNotMoveWhenMovingAway {
		return false
	}

	dx := targetPos.X - a.Position.X
	dy := targetPos.Y - a.Position.Y

	// Try to move horizontally first if dx is larger
	if abs(dx) >= abs(dy) {
		if dx > 0 && a.Position.X > 0 && !world.IsPositionOccupied(a.Position.X-1, a.Position.Y) {
			a.Position.X--
			return true
		} else if dx < 0 && a.Position.X < world.Width-1 && !world.IsPositionOccupied(a.Position.X+1, a.Position.Y) {
			a.Position.X++
			return true
		}
	}

	// Try to move vertically if horizontal movement not possible
	if dy > 0 && a.Position.Y > 0 && !world.IsPositionOccupied(a.Position.X, a.Position.Y-1) {
		a.Position.Y--
		return true
	} else if dy < 0 && a.Position.Y < world.Height-1 && !world.IsPositionOccupied(a.Position.X, a.Position.Y+1) {
		a.Position.Y++
		return true
	}

	return false
}
