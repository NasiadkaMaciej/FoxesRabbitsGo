package simulation

import "foxes-rabbits-simulation/internal/config"

type Grass struct {
	GrowthRate    int
	Amount        int
	MaxAmount     int
	RegrowthTimer int
	Config        *config.Config
}

func NewGrass(cfg *config.Config) *Grass {
	return &Grass{
		GrowthRate:    cfg.GrassGrowthRate,
		Amount:        cfg.InitialGrass,
		MaxAmount:     cfg.GrassMaxAmount,
		RegrowthTimer: 0,
		Config:        cfg,
	}
}

func (g *Grass) Grow() {
	if g.Amount < g.MaxAmount {
		g.RegrowthTimer++

		// Start regrowing after timer reaches threshold
		if g.RegrowthTimer >= g.Config.GrassRegrowthTimer {
			g.Amount += g.GrowthRate
			g.RegrowthTimer = 0
		}
	}
}

func (g *Grass) Eat(amount int) {
	if g.Amount >= amount {
		g.Amount -= amount
	} else {
		g.Amount = 0
	}
}
