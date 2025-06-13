package simulation

import (
	"foxes-rabbits-simulation/internal/config"
	"math/rand"
)

type World struct {
	Width     int
	Height    int
	Foxes     []*Fox
	Rabbits   []*Rabbit
	GrassGrid [][]*Grass
	Config    *config.Config
}

func NewWorld(cfg *config.Config) *World {
	world := &World{
		Width:   cfg.WorldWidth,
		Height:  cfg.WorldHeight,
		Foxes:   make([]*Fox, 0),
		Rabbits: make([]*Rabbit, 0),
		Config:  cfg,
	}

	// Initialize grass grid
	world.GrassGrid = make([][]*Grass, world.Width)
	for x := 0; x < world.Width; x++ {
		world.GrassGrid[x] = make([]*Grass, world.Height)
		for y := 0; y < world.Height; y++ {
			world.GrassGrid[x][y] = NewGrass(cfg)
		}
	}

	return world
}

// Update advances the simulation one step
func (w *World) Update() {
	// Update and collect new animals
	var newFoxes []*Fox
	var newRabbits []*Rabbit

	// Process foxes
	for i := 0; i < len(w.Foxes); i++ {
		if !w.Foxes[i].IsDead() {
			w.Foxes[i].Move(w)
			w.Foxes[i].Eat(w)

			if newFox := w.Foxes[i].Reproduce(w); newFox != nil {
				newFoxes = append(newFoxes, newFox)
			}
		}
	}

	// Process rabbits
	for i := 0; i < len(w.Rabbits); i++ {
		if !w.Rabbits[i].IsDead() {
			w.Rabbits[i].Move(w)
			x, y := w.Rabbits[i].Position.X, w.Rabbits[i].Position.Y
			w.Rabbits[i].Eat(w.GrassGrid[x][y])

			if newRabbit := w.Rabbits[i].Reproduce(w); newRabbit != nil {
				newRabbits = append(newRabbits, newRabbit)
			}
		}
	}

	// Add new animals
	w.Foxes = append(w.Foxes, newFoxes...)
	w.Rabbits = append(w.Rabbits, newRabbits...)

	// Remove dead animals using filter pattern
	w.Foxes = filterAlive(w.Foxes)
	w.Rabbits = filterAlive(w.Rabbits)

	// Grow grass
	for x := 0; x < w.Width; x++ {
		for y := 0; y < w.Height; y++ {
			w.GrassGrid[x][y].Grow()
		}
	}
}

// Helper functions to remove dead animals
func filterAlive[T Animal](animals []T) []T {
	alive := animals[:0]
	for _, animal := range animals {
		if !animal.IsDead() {
			alive = append(alive, animal)
		}
	}
	return alive
}

// Initialize populates the world with animals
func (w *World) Initialize(numFoxes, numRabbits int) {
	for i := 0; i < numFoxes; i++ {
		x, y := w.getRandomEmptyPosition()
		w.Foxes = append(w.Foxes, NewFox(x, y, w.Config))
	}

	for i := 0; i < numRabbits; i++ {
		x, y := w.getRandomEmptyPosition()
		w.Rabbits = append(w.Rabbits, NewRabbit(x, y, w.Config))
	}
}

// getRandomEmptyPosition finds an unoccupied position
func (w *World) getRandomEmptyPosition() (int, int) {
	for {
		x := rand.Intn(w.Width)
		y := rand.Intn(w.Height)

		if !w.IsPositionOccupied(x, y) {
			return x, y
		}
	}
}

// IsPositionOccupied checks if a position is occupied by any animal
func (w *World) IsPositionOccupied(x, y int) bool {
	// Check if position is outside world boundaries
	if x < 0 || x >= w.Width || y < 0 || y >= w.Height {
		return true
	}

	// Check if any fox is at this position
	for _, fox := range w.Foxes {
		if fox.Position.X == x && fox.Position.Y == y {
			return true
		}
	}

	// Check if any rabbit is at this position
	for _, rabbit := range w.Rabbits {
		if rabbit.Position.X == x && rabbit.Position.Y == y {
			return true
		}
	}

	return false
}
