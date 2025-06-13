package main

import (
	"fmt"
	"foxes-rabbits-simulation/internal/chart"
	"foxes-rabbits-simulation/internal/config"
	"foxes-rabbits-simulation/internal/simulation"
	"foxes-rabbits-simulation/internal/ui"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDL: %s\n", err)
		os.Exit(1)
	}
	defer sdl.Quit()

	cfg := config.NewConfig()
	world := simulation.NewWorld(cfg)
	world.Initialize(cfg.InitialFoxes, cfg.InitialRabbits)

	renderer, err := ui.NewRenderer("Foxes and Rabbits Simulation",
		cfg.WorldWidth*cfg.AnimalSize, cfg.WorldHeight*cfg.AnimalSize, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize renderer: %s\n", err)
		os.Exit(1)
	}
	defer renderer.Destroy()

	// Create chart window
	chartWindow, err := chart.NewChartWindow("Population Chart", cfg.WorldWidth*cfg.AnimalSize, cfg.WorldHeight*cfg.AnimalSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize chart window: %s\n", err)
		os.Exit(1)
	}
	defer chartWindow.Destroy()

	// Add initial data point
	chartWindow.AddDataPoint(len(world.Foxes), len(world.Rabbits))

	frameDelay := cfg.FrameTime * time.Millisecond

	for renderer.IsRunning() && chartWindow.IsRunning() {
		renderer.HandleEvents()
		chartWindow.HandleEvents()

		// Check for mouse events to add animals
		mouseEvent := renderer.GetMouseEvent()
		if mouseEvent.Triggered {
			x, y := mouseEvent.GridX, mouseEvent.GridY

			// Only add if position is not occupied
			if !world.IsPositionOccupied(x, y) {
				switch mouseEvent.Type {
				case 1: // Left click - add rabbit
					world.Rabbits = append(world.Rabbits, simulation.NewRabbit(x, y, cfg))
				case 2: // Right click - add fox
					world.Foxes = append(world.Foxes, simulation.NewFox(x, y, cfg))
				}
			}
		}

		// Update window title with current counts
		title := fmt.Sprintf("Foxes and Rabbits Simulation - Foxes: %d | Rabbits: %d",
			len(world.Foxes), len(world.Rabbits))
		renderer.SetTitle(title)

		// Update chart window title
		chartWindow.SetTitle(fmt.Sprintf("Population Chart - Foxes: %d | Rabbits: %d",
			len(world.Foxes), len(world.Rabbits)))

		renderer.Render(world)
		world.Update()

		// Update chart with new data
		chartWindow.AddDataPoint(len(world.Foxes), len(world.Rabbits))
		chartWindow.Render()

		time.Sleep(frameDelay)
	}
}
