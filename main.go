package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Vector3Floor returns a new Vector3 with each component rounded down to the nearest integer.
func Vector3Floor(v rl.Vector3) rl.Vector3 {
	return rl.NewVector3(
		float32(math.Floor(float64(v.X))),
		float32(math.Floor(float64(v.Y))),
		float32(math.Floor(float64(v.Z))),
	)
}

// Vector3Round returns a new Vector3 with each component rounded to the nearest integer.
func Vector3Round(v rl.Vector3) rl.Vector3 {
	return rl.NewVector3(
		float32(math.Round(float64(v.X))),
		float32(math.Round(float64(v.Y))),
		float32(math.Round(float64(v.Z))),
	)
}

func updateMousePosition(previousMousePosition *rl.Vector2, dbgMoveNumber *uint8) (rl.Vector2, bool) {
	currentMousePosition := rl.GetMousePosition()
	isMousePositionChanged := false
	if rl.Vector2Length(rl.Vector2Subtract(currentMousePosition, *previousMousePosition)) > 5 {
		*previousMousePosition = currentMousePosition
		isMousePositionChanged = true
		*dbgMoveNumber = (*dbgMoveNumber + 1) % 255
	}
	return currentMousePosition, isMousePositionChanged
}

func getCursorPosition(camera rl.Camera, scene *Scene, app *VxdiAppEditor, cursor *Voxel, cursor2 *Voxel) {
	mouseModel := scene.GetIntersections(camera)
	if mouseModel.CollisionHit == CollisionHitPlane || mouseModel.CollisionHit == CollisionHitNone {
		mouseModel = app.Guides.GetIntersections(camera)
		if mouseModel.CollisionHit == CollisionHitVoxel {
			mouseModel.CollisionHit = CollisionHitGuide
		}
	}

	modelPointInt := Vector3Floor(mouseModel.Point)
	modelPointNextInt := rl.Vector3Add(modelPointInt, rl.NewVector3(0, 1, 0))

	switch mouseModel.CollisionHit {
	case CollisionHitNone:
		modelPointNextInt = Vector3Round(mouseModel.Point)
		modelPointInt = modelPointNextInt
	case CollisionHitVoxel:
		modelPointInt = mouseModel.Voxel.Position
		modelPointNextInt = rl.Vector3Add(mouseModel.Voxel.Position, mouseModel.Normal)
	case CollisionHitPlane, CollisionHitGuide:
		modelPointNextInt = Vector3Round(mouseModel.Point)
		modelPointInt = modelPointNextInt
	}

	cursor.Position = modelPointInt
	cursor2.Position = modelPointNextInt
}
func doNothing() {}
func main() {
	screenWidth := int32(1280)
	screenHeight := int32(720)

	rl.SetConfigFlags(rl.FlagMsaa4xHint) //ENABLE 4X MSAA IF AVAILABLE

	rl.InitWindow(screenWidth, screenHeight, "raylib [shaders] example - basic lighting")

	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(2.0, 4.0, 6.0)
	camera.Target = rl.NewVector3(0.0, 0.5, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective
	light := VxdiLight{
		Direction:      rl.NewVector3(-1, -4, -2),
		ShadowStrength: .5,
		LightStrength:  .5,
	}
	scene := NewScene(true, light)
	scene.AddVoxel(&Voxel{Position: rl.Vector3{X: 0, Y: 1, Z: 2}, Material: rl.Red})
	cursor1 := NewVoxel()
	cursor1.Material = rl.Red
	cursor2 := NewVoxel()
	cursor2.Material = rl.Green

	orbit := NewOrbit(&camera)

	// lights[1] = NewLight(LightTypeDirectional, rl.NewVector3(2, 1, 2), rl.NewVector3(0, 0, 0), rl.Red, shader)

	// lights[2] = NewLight(LightTypeDirectional, rl.NewVector3(-2, 1, 2), rl.NewVector3(0, 0, 0), rl.Green, shader)

	// lights[3] = NewLight(LightTypeDirectional, rl.NewVector3(2, 1, -2), rl.NewVector3(0, 0, 0), rl.Blue, shader)

	rl.SetTargetFPS(60)

	// loop globals

	status := "KEYS [Y] [R] [G] [B] TURN LIGHTS ON/OFF" // In Go, a string is usually used for text, but a byte array can hold arbitrary binary data or ASCII text.
	/// var ctrl, leftCtrl int
	showHelp := false // Using byte instead of char, since Go doesn't have a char type. Byte is an alias for uint8.
	windowShouldClose := false
	currentMousePosition := rl.GetMousePosition() // Direct assignment; Go infers type from the function's return value.
	previousMousePosition := currentMousePosition // Copy the value from currentMousePosition.
	isMousePositionChanged := true
	var dbgMoveNumber uint8 = 0 // Unsigned char in C is equivalent to uint8 in Go.
	app := NewVxdiAppEditor(&camera, light)
	app.ConstructionHints.AddVoxel(cursor1)
	app.ConstructionHints.AddVoxel(cursor2)

	for !rl.WindowShouldClose() && !windowShouldClose {
		// update mouse model coordinates

		currentMousePosition := rl.GetMousePosition()
		isMousePositionChanged = false
		if rl.Vector2Length(rl.Vector2Subtract(currentMousePosition, previousMousePosition)) > 5 {
			previousMousePosition = currentMousePosition
			isMousePositionChanged = true
			dbgMoveNumber = (dbgMoveNumber + 1) % 255
		}
		mouseModel := scene.GetIntersections(camera)
		if mouseModel.CollisionHit == CollisionHitPlane || mouseModel.CollisionHit == CollisionHitNone {
			mouseModel = app.Guides.GetIntersections(camera)
			if mouseModel.CollisionHit == CollisionHitVoxel {
				mouseModel.CollisionHit = CollisionHitGuide
			}
		}

		modelPointInt := Vector3Floor(mouseModel.Point)
		modelPointNextInt := rl.Vector3Add(modelPointInt, rl.NewVector3(0, 1, 0))

		switch mouseModel.CollisionHit {
		case CollisionHitNone:
			modelPointNextInt = Vector3Round(mouseModel.Point)
			modelPointInt = modelPointNextInt
		case CollisionHitVoxel:
			modelPointInt = mouseModel.Voxel.Position
			modelPointNextInt = rl.Vector3Add(mouseModel.Voxel.Position, mouseModel.Normal)
		case CollisionHitPlane, CollisionHitGuide:
			modelPointNextInt = Vector3Round(mouseModel.Point)
			modelPointInt = modelPointNextInt
		}

		cursor1.Position = modelPointInt
		cursor2.Position = modelPointNextInt

		if isMousePositionChanged || showHelp {
			fmt.Printf("mouseMoved [%d] : (%v) {%v} \n", dbgMoveNumber, currentMousePosition, modelPointNextInt)
			doNothing()
		}
		orbit.ControlCamera()

		rl.UpdateCamera(&camera, rl.CameraCustom)

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawGrid(20, 1.0)
		scene.Draw(1, light)
		app.Guides.Draw(2, light)
		rl.DrawCubeWires(cursor1.Position, 1, 1, 1, cursor1.Material)
		rl.DrawCubeWires(cursor2.Position, 1, 1, 1, cursor2.Material)
		app.ConstructionHints.Draw(2, light)

		rl.EndMode3D()

		rl.DrawFPS(10, 10)

		status = fmt.Sprintf("mouseMoved [%d] : (%v) {%v} \n", dbgMoveNumber, currentMousePosition, modelPointNextInt)
		rl.DrawText(status, 10, 40, 20, rl.Black)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
