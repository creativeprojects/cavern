// +build !prod

package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var Debug = true

func (g *Game) displayDebug(screen *ebiten.Image) {
	template := " TPS: %0.2f \n Level %d - Colour %d \n Fruits %d - Pops %d - Orbs %d - Robots %d - Bolts %d \n%s"
	msg := fmt.Sprintf(template,
		ebiten.CurrentTPS(),
		g.level.id,
		g.level.colour,
		len(g.fruits),
		len(g.pops),
		len(g.orbs),
		len(g.robots),
		len(g.bolts),
		g.player,
	)
	fruitTemplate := " Fruit %d: ttl: %d coordinates: %s\n"
	for i, fruit := range g.fruits {
		msg += fmt.Sprintf(fruitTemplate, i, fruit.TTL, fruit.Sprite.String())
	}
	ebitenutil.DebugPrint(screen, msg)

	ebitenutil.DrawLine(screen, 70, 75, 70, 400, color.White)
	ebitenutil.DrawLine(screen, 730, 75, 730, 400, color.White)
	ebitenutil.DrawLine(screen, 70, 75, 730, 75, color.White)
	ebitenutil.DrawLine(screen, 70, 400, 730, 400, color.White)
}

// String returns a debug string
func (p *Player) String() string {
	return fmt.Sprintf(" Player score %d - health %d - lives %d - blow timer %d - hurt timer %d\n Player coordinates: %s\n",
		p.score,
		p.health,
		p.lives,
		p.blowTimer,
		p.hurtTimer,
		p.sprite.String(),
	)
}
