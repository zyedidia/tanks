package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/zyedidia/turbotanks/input"
)

const (
	twidth  = 26
	theight = 30
	tmass   = 10
)

type Tank struct {
	body    *cp.Body
	control *cp.Body

	input input.Controller

	ltspeed     float64
	rtspeed     float64
	targetAngle float64

	lastShot time.Time

	img *ebiten.Image
}

func NewTank(space *cp.Space) *Tank {
	control := space.AddBody(cp.NewKinematicBody())
	body := addBox(space, twidth, theight, tmass)

	pivot := space.AddConstraint(cp.NewPivotJoint2(control, body, cp.Vector{}, cp.Vector{}))
	pivot.SetMaxBias(0)
	pivot.SetMaxForce(1000)

	gear := space.AddConstraint(cp.NewGearJoint(control, body, 0.0, 1.0))
	gear.SetErrorBias(0) // attempt to fully correct the joint each step
	gear.SetMaxBias(1.2)
	gear.SetMaxForce(50000)

	// img := ebiten.NewImage(twidth, theight)
	// img.Fill(color.White)
	f, err := os.Open("assets/img/tank.png")
	if err != nil {
		log.Fatal(err)
	}
	tankimg, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	img := ebiten.NewImageFromImage(tankimg)

	t := &Tank{
		body:    body,
		control: control,
		input:   input.NewGamepad(0, DefaultGamepad),
		img:     img,
	}
	t.body.UserData = t

	return t
}

func (t *Tank) Update(space *cp.Space) {
	ldrive := t.input.Get(ActionLeftDrive) - t.input.Get(ActionLeftReverse)
	rdrive := t.input.Get(ActionRightDrive) - t.input.Get(ActionRightReverse)
	t.ltspeed = 2.0 * ldrive
	t.rtspeed = 2.0 * rdrive

	switch {
	case t.input.Get(ActionShoot) != 0:
		if time.Since(t.lastShot) >= 500*time.Millisecond {
			NewBullet(space, 100, 100)
			fmt.Println("Pew pew")
			t.lastShot = time.Now()
		}
	case t.input.Get(ActionReload) != 0:
		fmt.Println("Reload")
	}

	pos := t.body.Position()
	angle := t.body.Angle()
	rtpos := pos.Add(cp.Vector{twidth / 2 * math.Cos(angle), twidth / 2 * math.Sin(angle)})
	ltpos := pos.Add(cp.Vector{-twidth / 2 * math.Cos(angle), -twidth / 2 * math.Sin(angle)})

	nrtpos := rtpos.Add(cp.Vector{t.rtspeed * math.Sin(angle), -t.rtspeed * math.Cos(angle)})
	nltpos := ltpos.Add(cp.Vector{t.ltspeed * math.Sin(angle), -t.ltspeed * math.Cos(angle)})
	targetPos := cp.Vector{(nrtpos.X + nltpos.X) / 2, (nrtpos.Y + nltpos.Y) / 2}
	targetAngle := math.Atan2(nrtpos.Y-nltpos.Y, nrtpos.X-nltpos.X)

	diffX := targetPos.X - pos.X
	diffY := targetPos.Y - pos.Y
	if math.Abs(diffX) < 0.1 {
		diffX = 0
	}
	if math.Abs(diffY) < 0.1 {
		diffY = 0
	}
	t.control.SetVelocityVector(cp.Vector{
		X: diffX * 50,
		Y: diffY * 50,
	})

	diff := math.Atan2(math.Sin(targetAngle-angle), math.Cos(targetAngle-angle))
	targetAngle = angle + diff

	t.control.SetAngle(targetAngle)
}

func (t *Tank) Draw(screen *ebiten.Image) {
	pos := t.body.Position()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-twidth/2, -theight/2)
	op.GeoM.Rotate(t.body.Angle())
	op.GeoM.Translate(pos.X, pos.Y)

	screen.DrawImage(t.img, op)
}

func addBox(space *cp.Space, width, height, mass float64) *cp.Body {
	body := space.AddBody(cp.NewBody(mass, cp.MomentForBox(mass, width, height)))
	body.SetPosition(cp.Vector{100, 100})

	shape := space.AddShape(cp.NewBox(body, width, height, 0))
	shape.SetElasticity(0)
	shape.SetFriction(0.7)
	shape.SetCollisionType(CollisionTank)
	return body
}

func clamp(val, min, max float64) float64 {
	if val > max {
		val = max
	}
	if val < min {
		val = min
	}
	return val
}
