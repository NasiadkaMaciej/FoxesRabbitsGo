package ui

import (
	"foxes-rabbits-simulation/internal/config"
	"foxes-rabbits-simulation/internal/simulation"

	"github.com/veandco/go-sdl2/sdl"
)

type MouseAction struct {
	Action string // "AddFox", "AddRabbit", or "" // ToDo: "RemoveAnimal"?
	X      int
	Y      int
}

type Renderer struct {
	window         *sdl.Window
	renderer       *sdl.Renderer
	config         *config.Config
	leftMouseDown  bool
	rightMouseDown bool
}

func NewRenderer(title string, width, height int, cfg *config.Config) (*Renderer, error) {
	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(width), int32(height), sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}

	return &Renderer{
		window:   window,
		renderer: renderer,
		config:   cfg,
	}, nil
}

func (r *Renderer) HandleEvents() MouseAction {
	var action MouseAction

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.MouseButtonEvent:
			if e.Type == sdl.MOUSEBUTTONDOWN {
				r.leftMouseDown = e.Button == sdl.BUTTON_LEFT
				r.rightMouseDown = e.Button == sdl.BUTTON_RIGHT
			} else if e.Type == sdl.MOUSEBUTTONUP {
				if e.Button == sdl.BUTTON_LEFT {
					r.leftMouseDown = false
				} else if e.Button == sdl.BUTTON_RIGHT {
					r.rightMouseDown = false
				}
			}
		}
	}

	if r.leftMouseDown || r.rightMouseDown {
		mouseX, mouseY, _ := sdl.GetMouseState()
		gridX := int(mouseX) / r.config.AnimalSize
		gridY := int(mouseY) / r.config.AnimalSize

		if r.leftMouseDown {
			action = MouseAction{Action: "AddRabbit", X: gridX, Y: gridY}
		} else if r.rightMouseDown {
			action = MouseAction{Action: "AddFox", X: gridX, Y: gridY}
		}
	}

	return action
}

func (r *Renderer) Render(world *simulation.World) {
	// Clear screen
	r.renderer.SetDrawColor(255, 255, 255, 255)
	r.renderer.Clear()

	// Draw grass
	for x := 0; x < world.Width; x++ {
		for y := 0; y < world.Height; y++ {
			r.drawGrass(x, y, world.GrassGrid[x][y].Amount)
		}
	}

	// Draw rabbits and foxes
	for _, rabbit := range world.Rabbits {
		r.drawAnimal(rabbit.Position.X, rabbit.Position.Y, rabbit.Config.RabbitColor)
	}

	for _, fox := range world.Foxes {
		r.drawAnimal(fox.Position.X, fox.Position.Y, fox.Config.FoxColor)
	}

	r.renderer.Present()
}

func (r *Renderer) drawAnimal(x, y int, color sdl.Color) {
	r.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	size := r.config.AnimalSize
	rect := sdl.Rect{X: int32(x * size), Y: int32(y * size), W: int32(size), H: int32(size)}
	r.renderer.FillRect(&rect)
}

func (r *Renderer) drawGrass(x, y int, amount int) {
	// Clamp amount between 0 and max
	amount = max(0, min(amount, r.config.GrassMaxAmount))

	// Calculate green intensity
	baseColor := r.config.GrassBaseColor
	ratio := float64(amount) / float64(r.config.GrassMaxAmount)
	greenValue := uint8(float64(baseColor.G) + ratio*(255.0-float64(baseColor.G)))

	// Draw grass
	r.renderer.SetDrawColor(baseColor.R, greenValue, baseColor.B, baseColor.A)
	size := r.config.AnimalSize
	rect := sdl.Rect{X: int32(x * size), Y: int32(y * size), W: int32(size), H: int32(size)}
	r.renderer.FillRect(&rect)
}

func (r *Renderer) SetTitle(title string) {
	r.window.SetTitle(title)
}
