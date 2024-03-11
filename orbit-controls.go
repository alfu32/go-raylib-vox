package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Orbit struct {
	LastMousePos rl.Vector2
	Azimuth      float32 // Horizontal angle
	Elevation    float32 // Vertical angle
	Radius       float32 // Distance from the target
	IsOrbiting   bool
	IsPanning    bool
	Camera       *rl.Camera
}

func NewOrbit(camera *rl.Camera) *Orbit {
	return &Orbit{
		LastMousePos: rl.GetMousePosition(),
		Azimuth:      0.709,
		Elevation:    0.709,
		Radius:       17.3,
		IsOrbiting:   false,
		IsPanning:    false,
		Camera:       camera,
	}
}

func (o *Orbit) ControlCamera() int {
	mouseDelta := rl.Vector2Subtract(rl.GetMousePosition(), o.LastMousePos)
	o.LastMousePos = rl.GetMousePosition()
	currentCameraDist := rl.Vector3Length(rl.Vector3Subtract(o.Camera.Position, o.Camera.Target))

	// Check for right mouse button pressed for orbiting
	if rl.IsMouseButtonDown(rl.MouseLeftButton) && rl.IsKeyDown(rl.KeyLeftControl) {
		o.IsOrbiting = true
		o.IsPanning = false
	} else if rl.IsMouseButtonDown(rl.MouseLeftButton) && rl.IsKeyDown(rl.KeyLeftShift) {
		// Check for right mouse button pressed with Shift for panning
		o.IsPanning = true
		o.IsOrbiting = false
	} else {
		o.IsOrbiting = false
		o.IsPanning = false
	}

	if o.IsOrbiting {
		o.Azimuth -= mouseDelta.X * 0.003
		if math.Abs(float64(o.Elevation+mouseDelta.Y*0.003)) < 1.57 {
			o.Elevation += mouseDelta.Y * 0.003
		}
	} else if o.IsPanning {
		right := rl.Vector3Normalize(rl.Vector3CrossProduct(rl.Vector3Subtract(o.Camera.Position, o.Camera.Target), o.Camera.Up))
		up := rl.Vector3Normalize(rl.Vector3CrossProduct(right, rl.Vector3Subtract(o.Camera.Position, o.Camera.Target)))
		panSpeed := currentCameraDist * 0.00075
		o.Camera.Target = rl.Vector3Add(o.Camera.Target, rl.Vector3Scale(right, mouseDelta.X*panSpeed))
		o.Camera.Target = rl.Vector3Add(o.Camera.Target, rl.Vector3Scale(up, mouseDelta.Y*panSpeed))
		o.Camera.Position = rl.Vector3Add(o.Camera.Position, rl.Vector3Scale(right, mouseDelta.X*panSpeed))
		o.Camera.Position = rl.Vector3Add(o.Camera.Position, rl.Vector3Scale(up, mouseDelta.Y*panSpeed))
	}

	o.Radius -= rl.GetMouseWheelMove() * currentCameraDist * 0.08
	o.Radius = rl.Clamp(o.Radius, 1.0, 200000000.0)

	if o.IsOrbiting || !o.IsPanning { // Update position only if orbiting or not panning
		o.Camera.Position.X = o.Camera.Target.X + o.Radius*float32(math.Cos(float64(o.Elevation))*math.Sin(float64(o.Azimuth)))
		o.Camera.Position.Y = o.Camera.Target.Y + o.Radius*float32(math.Sin(float64(o.Elevation)))
		o.Camera.Position.Z = o.Camera.Target.Z + o.Radius*float32(math.Cos(float64(o.Elevation))*math.Cos(float64(o.Azimuth)))
	}
	return 0
}
