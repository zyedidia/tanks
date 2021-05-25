package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
)

const (
	CollisionBullet cp.CollisionType = iota + 1
	CollisionTank
	CollisionWall
)

type GameObject interface {
	Update(space *cp.Space)
	Draw(screen *ebiten.Image)
}

type Match struct {
	space *cp.Space
}

func NewMatch() *Match {
	space := cp.NewSpace()
	space.Iterations = 1
	space.SleepTimeThreshold = 0.5

	width := float64(screenWidth)
	height := float64(screenHeight)

	sides := []cp.Vector{
		{0, 0}, {0, height},
		{width, 0}, {width, height},
		{0, 0}, {width, 0},
		{0, height}, {width, height},
	}

	for i := 0; i < len(sides); i += 2 {
		var seg *cp.Shape
		seg = space.AddShape(cp.NewSegment(space.StaticBody, sides[i], sides[i+1], 0))
		seg.SetCollisionType(CollisionWall)
		fmt.Println(sides[i], sides[i+1])
		seg.SetElasticity(1)
		seg.SetFriction(1)
	}

	NewBullet(space, 10, 200)
	NewTank(space)

	handler := space.NewCollisionHandler(CollisionBullet, CollisionTank)
	handler.BeginFunc = bulletTankCollision

	return &Match{
		space: space,
	}
}

func (m *Match) Update() (GameState, error) {
	m.space.EachBody(func(body *cp.Body) {
		if g, ok := body.UserData.(GameObject); ok {
			g.Update(m.space)
		}
	})

	m.space.Step(1.0 / float64(ebiten.MaxTPS()))

	return nil, nil
}

func (m *Match) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	m.space.EachBody(func(body *cp.Body) {
		if g, ok := body.UserData.(GameObject); ok {
			g.Draw(screen)
		}
	})
}

func removeBullet(space *cp.Space, key, data interface{}) {
	bullet, ok := key.(*cp.Shape)
	if !ok {
		return
	}

	space.RemoveShape(bullet)
	space.RemoveBody(bullet.Body())
}

func bulletTankCollision(arb *cp.Arbiter, space *cp.Space, userdata interface{}) bool {
	bullet, _ := arb.Shapes()

	space.AddPostStepCallback(removeBullet, bullet, nil)
	return false
}
