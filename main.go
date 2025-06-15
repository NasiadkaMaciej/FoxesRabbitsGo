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

	chartWindow, err := chart.NewChartWindow("Population Chart", cfg.WorldWidth*cfg.AnimalSize, cfg.WorldHeight*cfg.AnimalSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize chart window: %s\n", err)
		os.Exit(1)
	}

	chartWindow.AddDataPoint(len(world.Foxes), len(world.Rabbits))

	frameDelay := cfg.FrameTime * time.Millisecond

	for {
		// Handle events in both windows
		mouseAction := renderer.HandleEvents()

		// Process mouse actions
		if mouseAction.Action != "" && !world.IsPositionOccupied(mouseAction.X, mouseAction.Y) {
			switch mouseAction.Action {
			case "AddRabbit":
				world.Rabbits = append(world.Rabbits, simulation.NewRabbit(mouseAction.X, mouseAction.Y, cfg))
			case "AddFox":
				world.Foxes = append(world.Foxes, simulation.NewFox(mouseAction.X, mouseAction.Y, cfg))
			}
		}

		// Update simulation and UI
		world.Update()

		// Update titles
		foxCount, rabbitCount := len(world.Foxes), len(world.Rabbits)
		renderer.SetTitle(fmt.Sprintf("Foxes and Rabbits Simulation - Foxes: %d | Rabbits: %d", foxCount, rabbitCount))
		chartWindow.SetTitle(fmt.Sprintf("Population Chart - Foxes: %d | Rabbits: %d", foxCount, rabbitCount))

		// Render windows
		renderer.Render(world)
		chartWindow.AddDataPoint(foxCount, rabbitCount)
		chartWindow.Render()

		time.Sleep(frameDelay)
	}
}
