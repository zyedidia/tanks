package main

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
)

type Bullet struct {
	body  *cp.Body
	shape *cp.Shape

	spawning  bool
	spawnTime time.Time

	img *ebiten.Image
}

func NewBullet(space *cp.Space, x, y, speed, angle float64) *Bullet {
	body := space.AddBody(cp.NewBody(5, 10))
	body.SetPosition(cp.Vector{x, y})
	body.SetVelocity(speed*math.Cos(angle), speed*math.Sin(angle))

	shape := space.AddShape(cp.NewCircle(body, 2.5, cp.Vector{}))
	shape.SetElasticity(1)
	shape.SetFriction(0)
	shape.SetCollisionType(CollisionBullet)

	img := ebiten.NewImage(5, 5)
	img.Fill(color.White)

	b := &Bullet{
		shape:     shape,
		body:      body,
		spawning:  true,
		spawnTime: time.Now(),
		img:       assets.images["bullet.png"],
	}
	b.body.UserData = b
	return b
}

const bulletLifetime = time.Second * 4

func (b *Bullet) Update(space *cp.Space) {
	if time.Since(b.spawnTime) >= bulletLifetime {
		space.RemoveShape(b.shape)
		space.RemoveBody(b.body)
	}
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	pos := b.body.Position()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-2.5, -2.5)
	op.GeoM.Rotate(b.body.Angle())
	op.GeoM.Translate(pos.X, pos.Y)

	screen.DrawImage(b.img, op)
}
