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

// ///////////////////   func updateMousePosition(previousMousePosition *rl.Vector2, dbgMoveNumber *uint8) (rl.Vector2, bool) {
// ///////////////////   	app.Mouse2 := rl.GetMousePosition()
// ///////////////////   	isMousePositionChanged := false
// ///////////////////   	if rl.Vector2Length(rl.Vector2Subtract(app.Mouse2, *previousMousePosition)) > 5 {
// ///////////////////   		*previousMousePosition = app.Mouse2
// ///////////////////   		isMousePositionChanged = true
// ///////////////////   		*dbgMoveNumber = (*dbgMoveNumber + 1) % 255
// ///////////////////   	}
// ///////////////////   	return app.Mouse2, isMousePositionChanged
// ///////////////////   }
// ///////////////////
// ///////////////////   func getCursorPosition(camera rl.Camera, scene *Scene, app *VxdiAppEditor, cursor *Voxel, cursor2 *Voxel) {
// ///////////////////   	app.Mouse3 := scene.GetIntersections(&camera)
// ///////////////////   	if app.Mouse3.CollisionHit == CollisionHitPlane || app.Mouse3.CollisionHit == CollisionHitNone {
// ///////////////////   		app.Mouse3 = app.Guides.GetIntersections(&camera)
// ///////////////////   		if app.Mouse3.CollisionHit == CollisionHitVoxel {
// ///////////////////   			app.Mouse3.CollisionHit = CollisionHitGuide
// ///////////////////   		}
// ///////////////////   	}
// ///////////////////
// ///////////////////   	app.Mouse3Int := Vector3Floor(app.Mouse3.Point)
// ///////////////////   	app.Mouse3IntNext := rl.Vector3Add(app.Mouse3Int, rl.NewVector3(0, 1, 0))
// ///////////////////
// ///////////////////   	switch app.Mouse3.CollisionHit {
// ///////////////////   	case CollisionHitNone:
// ///////////////////   		app.Mouse3IntNext = Vector3Round(app.Mouse3.Point)
// ///////////////////   		app.Mouse3Int = app.Mouse3IntNext
// ///////////////////   	case CollisionHitVoxel:
// ///////////////////   		app.Mouse3Int = app.Mouse3.Voxel.Position
// ///////////////////   		app.Mouse3IntNext = rl.Vector3Add(app.Mouse3.Voxel.Position, app.Mouse3.Normal)
// ///////////////////   	case CollisionHitPlane, CollisionHitGuide:
// ///////////////////   		app.Mouse3IntNext = Vector3Round(app.Mouse3.Point)
// ///////////////////   		app.Mouse3Int = app.Mouse3IntNext
// ///////////////////   	}
// ///////////////////
// ///////////////////   	cursor.Position = app.Mouse3Int
// ///////////////////   	cursor2.Position = app.Mouse3IntNext
// ///////////////////   }
func doNothing() {}
func (app *VxdiAppEditor) AddGuides(point rl.Vector3) {

	if rl.IsKeyPressed(rl.KeyG) {
		app.Guides.AddVoxel(InitVoxel(point.X, point.Y, point.Z, rl.White))
		for i := 1; i <= 10; i++ {
			app.Guides.AddVoxel(InitVoxel(point.X+float32(i), point.Y, point.Z, rl.Red))
			app.Guides.AddVoxel(InitVoxel(point.X-float32(i), point.Y, point.Z, rl.Fade(rl.Red, 0.5)))
			app.Guides.AddVoxel(InitVoxel(point.X, point.Y+float32(i), point.Z, rl.Green))
			app.Guides.AddVoxel(InitVoxel(point.X, point.Y-float32(i), point.Z, rl.Fade(rl.Green, 0.5)))
			app.Guides.AddVoxel(InitVoxel(point.X, point.Y, point.Z+float32(i), rl.Blue))
			app.Guides.AddVoxel(InitVoxel(point.X, point.Y, point.Z-float32(i), rl.Fade(rl.Blue, 0.5)))
		}
	}
}
func box_contains(box *rl.Rectangle, p *rl.Vector2) bool {
	return (box.X <= p.X) && (p.X <= box.X+box.Width) && (box.Y <= p.Y) && (p.Y <= box.Y+box.Height)
}
func main() {
	/// pointTool := MultistepTool2Points(func(a, b rl.Vector3, scene *Scene, material rl.Color, materialID uint, operator VoxelOperatorFn) {
	/// 	operator(scene, a, material, 0)
	/// })

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
	// scene.ImportScene("temp.exp.vxdi")
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
	isMousePositionChanged := true
	var dbgMoveNumber uint8 = 0 // Unsigned char in C is equivalent to uint8 in Go.
	app := NewVxdiAppEditor(&camera, light)
	app.Mouse2 = rl.GetMousePosition()  // Direct assignment; Go infers type from the function's return value.
	previousMousePosition := app.Mouse2 // Copy the value from app.Mouse2.
	app.ConstructionHints.AddVoxel(cursor1)
	app.ConstructionHints.AddVoxel(cursor2)

	rl.SetConfigFlags(rl.FlagMsaa4xHint | rl.FlagWindowResizable)

	rl.InitWindow(app.ScreenWidth, app.ScreenHeight, "raylib [shaders] example - basic lighting")

	app.Mouse3 = scene.GetIntersections(&camera)

	app.Mouse3Int = Vector3Floor(app.Mouse3.Point)
	app.Mouse3IntNext = rl.Vector3Add(app.Mouse3Int, rl.NewVector3(0, 1, 0))

	helpTool := NewTool(2, func(t *Tool) {
		app.Status = fmt.Sprintf("help tool not finished : points:%v, point:%v\n", t.Points, app.Mouse3IntNext)
	}, func(t *Tool) {
		app.Status = fmt.Sprintf("help tool finished : points:%v, point:%v\n", t.Points, app.Mouse3IntNext)
	})
	selectTool := NewTool(2, func(t *Tool) {
		app.Status = fmt.Sprintf("select tool not finished : points:%v, point:%v\n", t.Points, app.Mouse3IntNext)
	}, func(t *Tool) {
		app.Status = fmt.Sprintf("select tool finished : points:%v, point:%v\n", t.Points, app.Mouse3IntNext)
	})
	voxelTool := NewTool(1, func(t *Tool) {
		app.Status = fmt.Sprintf("voxel tool not finished : points:%v, point:%v\n", t.Points, app.Mouse3IntNext)
	}, func(t *Tool) {
		lastPoint := t.Points[len(t.Points)-1]
		if rl.IsKeyDown(rl.KeyLeftAlt) {
			app.Status = fmt.Sprintf("delete voxel : app.Mouse3Int:%v,app.Mouse3IntNext:%v,points:%v, point:%v\n", app.Mouse3Int, app.Mouse3IntNext, t.Points, lastPoint)
			scene.RemoveVoxel(app.Mouse3Int.X, app.Mouse3Int.Y, app.Mouse3Int.Z)
		} else {
			app.Status = fmt.Sprintf("add voxel : app.Mouse3Int:%v,app.Mouse3IntNext:%v,points:%v, point:%v\n", app.Mouse3Int, app.Mouse3IntNext, t.Points, lastPoint)
			scene.AddVoxel(InitVoxel(app.Mouse3IntNext.X, app.Mouse3IntNext.Y, app.Mouse3IntNext.Z, rl.Red))
		}
	})
	lineTool := NewTool(2, func(t *Tool) {
		app.Status = fmt.Sprintf("line tool not finished : points:%v, point:%v\n", t.Points, app.Mouse3IntNext)
	}, func(t *Tool) {
		app.Status = fmt.Sprintf("line tool finished : points:%v, point:%v\n", t.Points, app.Mouse3IntNext)
	})
	structureTool := NewTool(2, func(t *Tool) {
		app.Status = fmt.Sprintf("structure tool not finished : points:%v, point:%v\n", t.Points, app.Mouse3IntNext)
	}, func(t *Tool) {
		app.Status = fmt.Sprintf("structure tool finished : points:%v, point:%v\n", t.Points, app.Mouse3IntNext)
	})
	shellTool := NewTool(2, func(t *Tool) {
		app.Status = fmt.Sprintf("shell tool not finished : points:%v, point:%v\n", t.Points, app.Mouse3IntNext)
	}, func(t *Tool) {
		app.Status = fmt.Sprintf("shell tool finished : points:%v, point:%v\n", t.Points, app.Mouse3IntNext)
	})
	volumeTool := NewTool(2, func(t *Tool) {
		app.Status = fmt.Sprintf("volume tool not finished : points:%v, point:%v\n", t.Points, app.Mouse3IntNext)
	}, func(t *Tool) {
		app.Status = fmt.Sprintf("volume tool finished : points:%v, point:%v\n", t.Points, app.Mouse3IntNext)
	})
	app.RegisterTool("help", helpTool)
	app.RegisterTool("select", selectTool)
	app.RegisterTool("voxel", voxelTool)
	app.RegisterTool("line", lineTool)
	app.RegisterTool("structure", structureTool)
	app.RegisterTool("shell", shellTool)
	app.RegisterTool("volume", volumeTool)

	for !rl.WindowShouldClose() && !windowShouldClose {
		app.ScreenWidth = int32(rl.GetScreenWidth())
		app.ScreenHeight = int32(rl.GetScreenHeight())
		drawBox := rl.NewRectangle(300, 0, float32(app.ScreenWidth)-364, float32(app.ScreenHeight)-40)
		rightPanelBox := rl.NewRectangle(float32(app.ScreenWidth)-64, 0, 64, float32(64*len(app.ToolNames)))
		leftPanelBox := rl.NewRectangle(0, 0, 300, float32(app.ScreenHeight)-40)

		// update mouse model coordinates

		app.Mouse2 = (rl.GetMousePosition())
		isMousePositionChanged = false
		if rl.Vector2Length(rl.Vector2Subtract(app.Mouse2, previousMousePosition)) > 5 {
			previousMousePosition = app.Mouse2
			isMousePositionChanged = true
			dbgMoveNumber = (dbgMoveNumber + 1) % 255
		}
		app.Mouse3 = scene.GetIntersections(&camera)
		if app.Mouse3.CollisionHit == CollisionHitPlane || app.Mouse3.CollisionHit == CollisionHitNone {
			app.Mouse3 = app.Guides.GetIntersections(&camera)
			if app.Mouse3.CollisionHit == CollisionHitVoxel {
				app.Mouse3.CollisionHit = CollisionHitGuide
			}
		}

		app.Mouse3Int = Vector3Floor(app.Mouse3.Point)
		app.Mouse3IntNext = rl.Vector3Add(app.Mouse3Int, rl.NewVector3(0, 1, 0))

		switch app.Mouse3.CollisionHit {
		case CollisionHitNone:
			app.Mouse3IntNext = Vector3Round(app.Mouse3.Point)
			app.Mouse3Int = app.Mouse3IntNext
		case CollisionHitVoxel:
			app.Mouse3Int = app.Mouse3.Voxel.Position
			app.Mouse3IntNext = rl.Vector3Add(app.Mouse3.Voxel.Position, app.Mouse3.Normal)
		case CollisionHitPlane, CollisionHitGuide:
			app.Mouse3IntNext = Vector3Round(app.Mouse3.Point)
			app.Mouse3Int = app.Mouse3IntNext
		}
		if rl.IsKeyPressed(rl.KeyLeftAlt) {
			app.Mouse3IntNext = app.Mouse3Int
		}

		cursor1.Position = app.Mouse3Int
		cursor2.Position = app.Mouse3IntNext

		if isMousePositionChanged || showHelp {
			// fmt.Printf("mouseMoved [%d] : (%v) {%v} \n", dbgMoveNumber, app.Mouse2, app.Mouse3IntNext)
			doNothing()
		}
		if box_contains(&drawBox, &app.Mouse2) && rl.IsMouseButtonReleased(rl.MouseButtonLeft) && !rl.IsKeyDown(rl.KeyLeftControl) && !rl.IsKeyDown(rl.KeyLeftShift) {
			app.Tools[app.CurrentTool].Next(app.Mouse3IntNext)
		}
		orbit.ControlCamera()
		app.AddGuides(app.Mouse3Int)

		rl.UpdateCamera(&camera, rl.CameraCustom)

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawGrid(20, 1.0)
		scene.Draw(1, light)
		// You can then use Raylib functions to check for intersections, etc.
		if box_contains(&drawBox, &app.Mouse2) {
			/// for _, v := range app.Guides.Voxels {
			/// 	rl.DrawCubeWires(v.Position, VOXEL_SZ, VOXEL_SZ, VOXEL_SZ, rl.DarkGray)
			/// }
			app.Guides.Draw(2, light)
			//rl.DrawCubeWires(cursor1.Position, 1, 1, 1, cursor1.Material)
			//rl.DrawCubeWires(cursor2.Position, 1, 1, 1, cursor2.Material)
			app.ConstructionHints.Draw(2, light)

			rl.DrawLine3D(
				app.Mouse3.Point,
				rl.Vector3Add(
					app.Mouse3.Point,
					rl.NewVector3(
						app.Mouse3.Normal.X*1,
						app.Mouse3.Normal.Y*1,
						app.Mouse3.Normal.Z*1,
					),
				),
				rl.NewColor(
					uint8(app.Mouse3.Normal.X*255),
					uint8(app.Mouse3.Normal.Y*255),
					uint8(app.Mouse3.Normal.Z*255),
					255,
				))
			InitVoxel(app.Mouse3.Point.X, app.Mouse3.Point.Y, app.Mouse3.Point.Z, rl.Blue).DrawShaded(light, VOXEL_SZ2/2)
		}
		// for _, v := range app.ConstructionHints.Voxels {
		// 	rl.DrawCubeWires(v.Position, VOXEL_SZ, VOXEL_SZ, VOXEL_SZ, rl.DarkGray)
		// }
		// rl.DrawLine3D(rl.NewVector3(0, 0, 0), app.Mouse3Int, rl.Green)

		rl.EndMode3D()
		rl.DrawRectangleLinesEx(leftPanelBox, 2, rl.Red)
		rl.DrawRectangleLinesEx(rightPanelBox, 2, rl.Green)
		rl.DrawRectangleLinesEx(drawBox, 2, rl.Blue)

		posy64 := int32(app.Mouse2.Y/64) * 64
		if box_contains(&rightPanelBox, &app.Mouse2) {
			rl.DrawRectangle(int32(app.ScreenWidth)-64, posy64, 64, 64, rl.NewColor(127, 255, 0, 127))
		}
		for index, toolName := range app.ToolNames {
			func(i int, name string) {
				if app.CurrentTool == name {
					rl.DrawRectangle(int32(app.ScreenWidth)-64, 64*int32(i), 64, 64, rl.NewColor(255, 127, 127, 255))
				}
				if box_contains(&rightPanelBox, &app.Mouse2) && int32(posy64/64) == int32(i) && rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
					fmt.Printf("setting mode to %s\n", name)
					app.CurrentTool = name
					for _, t := range app.Tools {
						t.Reset()
					}
				}
				rl.DrawRectangleLines(int32(app.ScreenWidth)-64, int32(i)*64, 64, 64, rl.NewColor(127, 255, 0, 255))
				rl.DrawText(name, int32(app.ScreenWidth)-64, int32(i)*64+48, 12, rl.NewColor(0, 0, 0, 255))
			}(index, toolName)
		}

		rl.DrawFPS(10, 10)

		status = fmt.Sprintf("current tool : %s, scene: %d, guides: %d, helpers: %d, mouse [%d] : (%v) {%v} \n", app.CurrentTool, len(scene.Voxels), len(app.Guides.Voxels), len(app.ConstructionHints.Voxels), dbgMoveNumber, app.Mouse2, app.Mouse3IntNext)
		rl.DrawText(status, 10, 40, 20, rl.Black)

		rl.DrawText(app.Status, 0, app.ScreenHeight-20, 20, rl.Black)

		rl.EndDrawing()
		// time.Sleep(1000000000)
	}
	// scene.SaveScene("temp.scene.vxdi")
	// scene.ExportScene("temp.exp.vxdi")
	rl.CloseWindow()
}
