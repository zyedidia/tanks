package main

import (
	_ "image/png"
	"math"
	"time"

	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/zyedidia/turbotanks/input"
)

const (
	twidth  = 26
	theight = 22
	tmass   = 10
)

type Tank struct {
	body    *cp.Body
	control *cp.Body

	input input.Controller

	ltspeed     float64
	rtspeed     float64
	targetAngle float64

	lastShot   time.Time
	lastReload time.Time

	health  int
	bullets int32

	ltrack  *AnimImage
	rtrack  *AnimImage
	chassis *ebiten.Image
	turret  *ebiten.Image
}

func NewTank(space *cp.Space, x, y, angle float64, input input.Controller) *Tank {
	control := space.AddBody(cp.NewKinematicBody())
	body := addBox(space, x, y, twidth, theight, tmass)
	body.SetAngle(angle)

	pivot := space.AddConstraint(cp.NewPivotJoint2(control, body, cp.Vector{}, cp.Vector{}))
	pivot.SetMaxBias(0)
	pivot.SetMaxForce(800)

	gear := space.AddConstraint(cp.NewGearJoint(control, body, 0.0, 1.0))
	gear.SetErrorBias(0) // attempt to fully correct the joint each step
	gear.SetMaxBias(1.2)
	gear.SetMaxForce(50000)

	t := &Tank{
		body:    body,
		control: control,
		input:   input,
		chassis: assets.images["chassis.png"],
		turret:  assets.images["turret.png"],
		health:  10,
		bullets: 5,
		ltrack: &AnimImage{
			anim:  assets.anims["ltrack"],
			count: 0,
		},
		rtrack: &AnimImage{
			anim:  assets.anims["ltrack"],
			count: 0,
		},
	}
	t.body.UserData = t

	return t
}

const reloadTime = 3 * time.Second

func (t *Tank) Reload() {
	go func() {
		time.Sleep(reloadTime)
		atomic.StoreInt32(&t.bullets, 5)
	}()
}

func (t *Tank) Update(space *cp.Space) {
	pos := t.body.Position()
	angle := t.body.Angle()

	if t.health > 0 {
		ldrive := math.Max(t.input.Get(ActionDrive), t.input.Get(ActionRight)) - math.Max(t.input.Get(ActionReverse), t.input.Get(ActionLeft))
		rdrive := math.Max(t.input.Get(ActionDrive), t.input.Get(ActionLeft)) - math.Max(t.input.Get(ActionReverse), t.input.Get(ActionRight))
		t.ltspeed = 2.0 * ldrive
		t.rtspeed = 2.0 * rdrive

		switch {
		case t.input.Get(ActionShoot) != 0:
			if time.Since(t.lastShot) >= 500*time.Millisecond {
				t.lastShot = time.Now()
				if atomic.LoadInt32(&t.bullets) <= 0 {
					assets.sounds["klick.ogg"].Rewind()
					assets.sounds["klick.ogg"].Play()
					break
				}
				atomic.AddInt32(&t.bullets, -1)
				angle := angle - math.Pi/2
				NewBullet(space, pos.X+twidth/4*math.Cos(angle), pos.Y+theight/4*math.Sin(angle), 150, angle)
				assets.sounds["shoot.ogg"].Rewind()
				assets.sounds["shoot.ogg"].Play()

				t.body.ApplyImpulseAtWorldPoint(cp.Vector{-200 * math.Cos(angle), -200 * math.Sin(angle)}, t.body.Position())

				if atomic.LoadInt32(&t.bullets) <= 0 {
					t.Reload()
				}
			}
		}
	} else {
		t.ltspeed = 0
		t.rtspeed = 0
	}

	if t.ltspeed > 0 {
		t.ltrack.Forward()
	} else if t.ltspeed < 0 {
		t.ltrack.Backward()
	}
	if t.rtspeed > 0 {
		t.rtrack.Forward()
	} else if t.rtspeed < 0 {
		t.rtrack.Backward()
	}

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
	velocity := t.body.Velocity()
	if (diffX != 0 && sign(velocity.X) == -sign(diffX)) || (diffY != 0 && sign(velocity.Y) == -sign(diffY)) {
		if !assets.sounds["brake.ogg"].IsPlaying() {
			assets.sounds["brake.ogg"].Rewind()
			assets.sounds["brake.ogg"].Play()
		}
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
	screen.DrawImage(t.chassis, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-14/2, -21)
	op.GeoM.Rotate(t.body.Angle())
	op.GeoM.Translate(pos.X, pos.Y)
	screen.DrawImage(t.turret, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-3, -10)
	op.GeoM.Rotate(t.body.Angle())
	op.GeoM.Translate(pos.X-9*math.Cos(t.body.Angle()), pos.Y-9*math.Sin(t.body.Angle()))
	screen.DrawImage(t.ltrack.Image(), op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-3, -10)
	op.GeoM.Scale(-1, 1)
	op.GeoM.Rotate(t.body.Angle())
	op.GeoM.Translate(pos.X+9*math.Cos(t.body.Angle()), pos.Y+9*math.Sin(t.body.Angle()))
	screen.DrawImage(t.rtrack.Image(), op)
}

func addBox(space *cp.Space, x, y, width, height, mass float64) *cp.Body {
	body := space.AddBody(cp.NewBody(mass, cp.MomentForBox(mass, width, height)))
	body.SetPosition(cp.Vector{x, y})

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

func sign(a float64) int {
	if a < 0 {
		return -1
	} else if a > 0 {
		return 1
	}
	return 0
}
