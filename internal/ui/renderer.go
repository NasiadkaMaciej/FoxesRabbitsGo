package ui

import (
	"foxes-rabbits-simulation/internal/config"
	"foxes-rabbits-simulation/internal/simulation"

	"github.com/veandco/go-sdl2/sdl"
)

type Renderer struct {
	window     *sdl.Window
	renderer   *sdl.Renderer
	running    bool
	config     *config.Config
	MouseEvent MouseEvent
}

type MouseEvent struct {
	Type      int8 // 1 for left click, 2 for right click
	GridX     int
	GridY     int
	Triggered bool
}

func NewRenderer(title string, width, height int, cfg *config.Config) (*Renderer, error) {
	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(width), int32(height), sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		window.Destroy()
		return nil, err
	}

	return &Renderer{
		window:   window,
		renderer: renderer,
		running:  true,
		config:   cfg,
	}, nil
}

func (r *Renderer) HandleEvents() {
	// Reset mouse event
	r.MouseEvent.Triggered = false

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			r.running = false
		case *sdl.MouseButtonEvent:
			if e.Type == sdl.MOUSEBUTTONDOWN {
				// Convert pixel coordinates to grid coordinates
				gridX := int(e.X) / r.config.AnimalSize
				gridY := int(e.Y) / r.config.AnimalSize

				// Check if coordinates are within bounds
				if gridX >= 0 && gridX < r.config.WorldWidth &&
					gridY >= 0 && gridY < r.config.WorldHeight {

					r.MouseEvent.GridX = gridX
					r.MouseEvent.GridY = gridY
					r.MouseEvent.Triggered = true

					if e.Button == sdl.BUTTON_LEFT {
						r.MouseEvent.Type = 1 // Left click
					} else if e.Button == sdl.BUTTON_RIGHT {
						r.MouseEvent.Type = 2 // Right click
					}
				}
			}
		}
	}
}

func (r *Renderer) GetMouseEvent() MouseEvent {
	event := r.MouseEvent
	r.MouseEvent.Triggered = false
	return event
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
		r.drawRabbit(rabbit.Position.X, rabbit.Position.Y)
	}

	for _, fox := range world.Foxes {
		r.drawFox(fox.Position.X, fox.Position.Y)
	}

	r.renderer.Present()
}

func (r *Renderer) drawFox(x, y int) {
	color := r.config.FoxColor
	r.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	size := r.config.AnimalSize
	rect := sdl.Rect{X: int32(x * size), Y: int32(y * size), W: int32(size), H: int32(size)}
	r.renderer.FillRect(&rect)
}

func (r *Renderer) drawRabbit(x, y int) {
	color := r.config.RabbitColor
	r.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	size := r.config.AnimalSize
	rect := sdl.Rect{X: int32(x * size), Y: int32(y * size), W: int32(size), H: int32(size)}
	r.renderer.FillRect(&rect)
}

func (r *Renderer) drawGrass(x, y int, amount int) {
	// Clamp amount between 0 and max
	amount = max(0, min(amount, r.config.GrassMaxAmount))

	// Calculate green intensity - fix the type mismatch
	baseColor := r.config.GrassBaseColor
	ratio := float64(amount) / float64(r.config.GrassMaxAmount)
	greenValue := uint8(float64(baseColor.G) + ratio*(255.0-float64(baseColor.G)))

	// Draw grass
	r.renderer.SetDrawColor(baseColor.R, greenValue, baseColor.B, baseColor.A)
	size := r.config.AnimalSize
	rect := sdl.Rect{X: int32(x * size), Y: int32(y * size), W: int32(size), H: int32(size)}
	r.renderer.FillRect(&rect)
}

func (r *Renderer) IsRunning() bool {
	return r.running
}

func (r *Renderer) Destroy() {
	r.renderer.Destroy()
	r.window.Destroy()
}

func (r *Renderer) SetTitle(title string) {
	r.window.SetTitle(title)
}
