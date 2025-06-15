package simulation

import (
	"foxes-rabbits-simulation/internal/config"
	"math/rand"
)

type Animal interface {
	IsDead() bool
	GetPosition() Position
	CanEat(cooldown int) bool
}

type AnimalBase struct {
	Position               Position
	Energy                 int
	Config                 *config.Config
	TurnsSinceEaten        int
	TurnsSinceReproduction int
}

func (a *AnimalBase) IsDead() bool {
	return a.Energy <= 0
}

func (a *AnimalBase) GetPosition() Position {
	return a.Position
}

// FindEmptyAdjacentPosition finds an empty position nearby animal
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
	return a.TurnsSinceEaten >= cooldown
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
	minDistance := maxRange + 1

	for _, animal := range animals {
		dx := abs(from.GetPosition().X - animal.GetPosition().X)
		dy := abs(from.GetPosition().Y - animal.GetPosition().Y)

		distance := dx + dy

		if distance <= maxRange && (distance < minDistance || !foundAnimal) {
			nearest = animal
			minDistance = distance
			foundAnimal = true
		}
	}

	return nearest, foundAnimal
}

func (a *AnimalBase) MoveRandomly(world *World) bool {
	directions := []struct{ dx, dy int }{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1}, // Right, Left, Down, Up
	}

	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	// Try each direction
	for _, dir := range directions {
		newX := a.Position.X + dir.dx
		newY := a.Position.Y + dir.dy

		if newX >= 0 && newX < world.Width &&
			newY >= 0 && newY < world.Height &&
			!world.IsPositionOccupied(newX, newY) {
			a.Position.X, a.Position.Y = newX, newY
			return true
		}
	}

	return false // Couldn't move
}

// MoveDirectionally moves the animal toward or away from a target position
// If moveToward is true, animal moves toward the target, otherwise it moves away
func (a *AnimalBase) MoveDirectionally(targetPos Position, world *World, moveToward bool) bool {
	if !moveToward && rand.Float64() < a.Config.ChanceToStayStillWhenFleeing {
		return false
	}

	dx := targetPos.X - a.Position.X
	dy := targetPos.Y - a.Position.Y

	// Reverse direction if moving away
	if !moveToward {
		dx = -dx
		dy = -dy
	}

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

func (a *AnimalBase) MoveToward(targetPos Position, world *World) bool {
	return a.MoveDirectionally(targetPos, world, true)
}

func (a *AnimalBase) MoveAwayFrom(targetPos Position, world *World) bool {
	return a.MoveDirectionally(targetPos, world, false)
}
