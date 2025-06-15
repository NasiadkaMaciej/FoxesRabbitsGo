package chart

import (
	"github.com/veandco/go-sdl2/sdl"
)

type ChartWindow struct {
	window    *sdl.Window
	renderer  *sdl.Renderer
	data      []struct{ Foxes, Rabbits int }
	maxPoints int
	width     int32
	height    int32
	padding   int32
}

func NewChartWindow(title string, width, height int) (*ChartWindow, error) {
	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(width), int32(height), sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}

	return &ChartWindow{
		window:    window,
		renderer:  renderer,
		data:      make([]struct{ Foxes, Rabbits int }, 0),
		maxPoints: 300,
		width:     int32(width),
		height:    int32(height),
		padding:   20,
	}, nil
}

func (c *ChartWindow) AddDataPoint(foxes, rabbits int) {
	c.data = append(c.data, struct{ Foxes, Rabbits int }{foxes, rabbits})
	if len(c.data) > c.maxPoints {
		c.data = c.data[1:]
	}
}

func (c *ChartWindow) Render() {
	c.renderer.SetDrawColor(255, 255, 255, 255)
	c.renderer.Clear()

	if len(c.data) < 2 {
		c.renderer.Present()
		return
	}

	// Find the maximum population for scaling
	maxValue := 10 // Minimum scale
	for _, point := range c.data {
		if point.Foxes > maxValue {
			maxValue = point.Foxes
		}
		if point.Rabbits > maxValue {
			maxValue = point.Rabbits
		}
	}

	maxValue = ((maxValue + 9) / 10) * 10

	chartWidth := c.width - 2*c.padding
	chartHeight := c.height - 2*c.padding

	c.renderer.SetDrawColor(200, 200, 200, 255)

	// Horizontal grid lines
	step := maxValue / 5
	if step == 0 {
		step = 1
	}

	for y := step; y <= maxValue; y += step {
		yPos := c.height - c.padding - int32(float64(y)/float64(maxValue)*float64(chartHeight))
		c.renderer.DrawLine(c.padding, yPos, c.width-c.padding, yPos)
	}

	// Vertical grid lines
	timeStep := 5
	if len(c.data) <= timeStep {
		timeStep = 1
	}

	for i := 0; i <= c.maxPoints; i += timeStep {
		xPos := c.padding + int32(float64(i)*float64(chartWidth)/float64(c.maxPoints))
		if xPos <= c.padding {
			continue
		}
		c.renderer.DrawLine(xPos, c.padding, xPos, c.height-c.padding)
	}

	c.renderer.SetDrawColor(0, 0, 0, 255)
	c.renderer.DrawLine(c.padding, c.padding, c.padding, c.height-c.padding)                  // Y axis
	c.renderer.DrawLine(c.padding, c.height-c.padding, c.width-c.padding, c.height-c.padding) // X axis

	xStep := float64(chartWidth) / float64(c.maxPoints-1)

	// Draw fox population (red)
	c.renderer.SetDrawColor(255, 0, 0, 255)
	for i := 0; i < len(c.data)-1; i++ {
		x1 := c.padding + int32(float64(i)*xStep)
		y1 := c.height - c.padding - int32(float64(c.data[i].Foxes)/float64(maxValue)*float64(chartHeight))
		x2 := c.padding + int32(float64(i+1)*xStep)
		y2 := c.height - c.padding - int32(float64(c.data[i+1].Foxes)/float64(maxValue)*float64(chartHeight))
		c.renderer.DrawLine(x1, y1, x2, y2)
	}

	// Draw rabbit population (blue)
	c.renderer.SetDrawColor(0, 0, 255, 255)
	for i := 0; i < len(c.data)-1; i++ {
		x1 := c.padding + int32(float64(i)*xStep)
		y1 := c.height - c.padding - int32(float64(c.data[i].Rabbits)/float64(maxValue)*float64(chartHeight))
		x2 := c.padding + int32(float64(i+1)*xStep)
		y2 := c.height - c.padding - int32(float64(c.data[i+1].Rabbits)/float64(maxValue)*float64(chartHeight))
		c.renderer.DrawLine(x1, y1, x2, y2)
	}

	c.renderer.Present()
}

func (c *ChartWindow) SetTitle(title string) {
	c.window.SetTitle(title)
}
