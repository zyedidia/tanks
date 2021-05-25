package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jakecoffman/cp"
	"github.com/zyedidia/turbotanks/input"
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
	done  bool
	space *cp.Space
}

func NewMatch() *Match {
	space := cp.NewSpace()
	space.Iterations = 1
	space.SleepTimeThreshold = 0.5

	width := float64(screenWidth)
	height := float64(screenHeight)

	sides := []cp.Vector{
		{1, 0}, {1, height - 1},
		{width, 0}, {width, height - 1},
		{1, 0}, {width, 0},
		{1, height - 1}, {width, height - 1},

		{width / 3, height / 3}, {width / 3, 2 * height / 3},
		{2 * width / 3, height / 3}, {2 * width / 3, 2 * height / 3},
		{width / 3, height / 3}, {2 * width / 3, height / 3},
		{width / 3, 2 * height / 3}, {2 * width / 3, 2 * height / 3},
	}

	for i := 0; i < len(sides); i += 2 {
		var seg *cp.Shape
		seg = space.AddShape(cp.NewSegment(space.StaticBody, sides[i], sides[i+1], 0))
		seg.UserData = &Line{sides[i], sides[i+1]}
		seg.SetCollisionType(CollisionWall)
		fmt.Println(sides[i], sides[i+1])
		seg.SetElasticity(1)
		seg.SetFriction(1)
	}

	NewTank(space, width/8, height/8, math.Pi, input.NewKeyboard(DefaultKeyboard))
	NewTank(space, 7*width/8, 7*height/8, 0, input.NewGamepad(0, DefaultGamepad))

	m := &Match{
		done:  false,
		space: space,
	}

	handler := space.NewCollisionHandler(CollisionBullet, CollisionTank)
	handler.PreSolveFunc = bulletTankCollision
	handler.SeparateFunc = func(arb *cp.Arbiter, space *cp.Space, userdata interface{}) {
		body, _ := arb.Bodies()
		if bullet, ok := body.UserData.(*Bullet); ok {
			bullet.spawning = false
		}
	}
	handler.UserData = m

	bulletHandler := space.NewCollisionHandler(CollisionBullet, CollisionBullet)
	bulletHandler.BeginFunc = func(arb *cp.Arbiter, space *cp.Space, userdata interface{}) bool {
		return false
	}

	return m
}

func (m *Match) Update() (GameState, error) {
	m.space.EachBody(func(body *cp.Body) {
		if g, ok := body.UserData.(GameObject); ok {
			g.Update(m.space)
		}
	})

	m.space.Step(1.0 / float64(ebiten.MaxTPS()))

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return NewMatch(), nil
	}

	return nil, nil
}

func (m *Match) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	m.space.EachShape(func(shape *cp.Shape) {
		body := shape.Body()
		if g, ok := body.UserData.(GameObject); ok {
			g.Draw(screen)
		}

		if line, ok := shape.UserData.(*Line); ok {
			line.Draw(screen)
		}
	})

	if m.done {
		ebitenutil.DebugPrintAt(screen, "Game Over", screenWidth/2-30, screenHeight/2-3)
	}
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
	bullet, tank := arb.Shapes()

	if bullet, ok := bullet.Body().UserData.(*Bullet); ok {
		if bullet.spawning {
			return false
		}

	}

	tank.Body().ApplyImpulseAtWorldPoint(bullet.Body().Velocity().Mult(5), bullet.Body().Position())

	assets.sounds["explode.ogg"].Rewind()
	assets.sounds["explode.ogg"].Play()

	if tank, ok := tank.Body().UserData.(*Tank); ok {
		tank.health--
		if tank.health <= 0 {
			userdata.(*Match).done = true
		}
	}

	space.AddPostStepCallback(removeBullet, bullet, nil)
	return false
}

func bulletSpawn(arb *cp.Arbiter, space *cp.Space, userdata interface{}) {
	bullet, _ := arb.Shapes()

	bullet.SetCollisionType(CollisionBullet)
}
