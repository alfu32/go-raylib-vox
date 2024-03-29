package main

import (
	"fmt"
	"math"
	"os"

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
func doNothing(s string) {}
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
func DrawGridAtPoint(p rl.Vector3, numLines int, spacing float32) {
	MAX := float32(numLines) * (spacing)
	for k := 0; k <= numLines; k++ {
		rl.DrawLine3D(rl.NewVector3(p.X-spacing*float32(k), p.Y, p.Z+MAX), rl.NewVector3(p.X-spacing*float32(k), p.Y, p.Z-MAX), rl.LightGray)
		rl.DrawLine3D(rl.NewVector3(p.X+spacing*float32(k), p.Y, p.Z+MAX), rl.NewVector3(p.X+spacing*float32(k), p.Y, p.Z-MAX), rl.LightGray)
		rl.DrawLine3D(rl.NewVector3(p.X+MAX, p.Y, p.Z-spacing*float32(k)), rl.NewVector3(p.X-MAX, p.Y, p.Z-spacing*float32(k)), rl.LightGray)
		rl.DrawLine3D(rl.NewVector3(p.X+MAX, p.Y, p.Z+spacing*float32(k)), rl.NewVector3(p.X-MAX, p.Y, p.Z+spacing*float32(k)), rl.LightGray)
	}
	rl.DrawLine3D(rl.NewVector3(p.X, p.Y, p.Z+MAX), rl.NewVector3(p.X, p.Y, p.Z-MAX), rl.DarkGray)
	rl.DrawLine3D(rl.NewVector3(p.X+MAX, p.Y, p.Z), rl.NewVector3(p.X-MAX, p.Y, p.Z), rl.DarkGray)
}
func main() {
	/// pointTool := MultistepTool2Points(func(a, b rl.Vector3, scene *Scene, material rl.Color, materialID uint, operator VoxelOperatorFn) {
	/// 	operator(scene, a, material, 0)
	/// })
	fmt.Printf("args: %v", os.Args)
	filename := "temp.scene.cxdi"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(2.0, 4.0, 6.0)
	camera.Target = rl.NewVector3(0.0, 0.5, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective
	light := VxdiDirectionalLight{
		Direction:      rl.NewVector3(-1, -3, -2),
		ShadowStrength: .5,
		LightStrength:  .5,
	}
	// app.Layer.ImportScene("temp.exp.vxdi")
	cursor1 := NewVoxel()
	cursor1.Material = rl.Red
	cursor2 := NewVoxel()
	cursor2.Material = rl.Green

	orbit := NewOrbit(&camera)

	rl.SetTargetFPS(25)

	// loop globals

	status := "KEYS [Y] [R] [G] [B] TURN LIGHTS ON/OFF" // In Go, a string is usually used for text, but a byte array can hold arbitrary binary data or ASCII text.
	/// var ctrl, leftCtrl int
	/// showHelp := false // Using byte instead of char, since Go doesn't have a char type. Byte is an alias for uint8.
	windowShouldClose := false
	isMousePositionChanged := true
	var dbgMoveNumber uint8 = 0 // Unsigned char in C is equivalent to uint8 in Go.
	app := NewVxdiAppEditor(&camera, light)
	app.Mouse2 = rl.GetMousePosition()  // Direct assignment; Go infers type from the function's return value.
	previousMousePosition := app.Mouse2 // Copy the value from app.Mouse2.

	rl.SetConfigFlags(rl.FlagMsaa4xHint | rl.FlagWindowResizable)

	rl.InitWindow(app.ScreenWidth, app.ScreenHeight, "Vox D3i - Blox simulator")

	app.Mouse3 = app.Layer.GetIntersections(&camera)

	app.Mouse3Int = Vector3Floor(app.Mouse3.Point)
	app.Mouse3IntNext = rl.Vector3Add(app.Mouse3Int, rl.NewVector3(0, 1, 0))

	helpTool := NewTool(2, func(t *Tool) {
		app.Status = fmt.Sprintf("help tool not finished : points:%v, point:%v", t.Points, app.Mouse3IntNext)
	}, func(t *Tool) {
		app.Status = fmt.Sprintf("help tool finished : points:%v, point:%v", t.Points, app.Mouse3IntNext)
	})
	selectTool := NewTool(2, func(t *Tool) {
		app.Status = fmt.Sprintf("select tool not finished : points:%v, point:%v", t.Points, app.Mouse3IntNext)
	}, func(t *Tool) {
		app.Status = fmt.Sprintf("select tool finished : points:%v, point:%v", t.Points, app.Mouse3IntNext)
	})
	voxelTool := NewTool(1, func(t *Tool) {
		app.Status = fmt.Sprintf("voxel tool not finished : points:%v, point:%v", t.Points, app.Mouse3IntNext)
	}, func(t *Tool) {
		lastPoint := t.Points[len(t.Points)-1]
		if rl.IsKeyDown(rl.KeyLeftAlt) {
			app.Status = fmt.Sprintf("delete voxel : app.Mouse3Int:%v,app.Mouse3IntNext:%v,points:%v, point:%v", app.Mouse3Int, app.Mouse3IntNext, t.Points, lastPoint)
			app.Layer.RemoveVoxel(app.Mouse3Int.X, app.Mouse3Int.Y, app.Mouse3Int.Z)
		} else {
			app.Status = fmt.Sprintf("add voxel : app.Mouse3Int:%v,app.Mouse3IntNext:%v,points:%v, point:%v", app.Mouse3Int, app.Mouse3IntNext, t.Points, lastPoint)
			app.Layer.AddVoxel(InitVoxel(app.Mouse3IntNext.X, app.Mouse3IntNext.Y, app.Mouse3IntNext.Z, app.CurrentColor))
		}
	})
	lineTool := NewTool2Steps(
		app,
		rl.KeyLeftAlt,
		RasterizeLine,
	)
	structureTool := NewTool2Steps(
		app,
		rl.KeyLeftAlt,
		RasterizeStructureCube,
	)
	shellTool := NewTool2Steps(
		app,
		rl.KeyLeftAlt,
		RasterizeHollowCube,
	)
	volumeTool := NewTool2Steps(
		app,
		rl.KeyLeftAlt,
		RasterizeSolidCube,
	)
	planeTool := NewTool3Steps(
		app,
		rl.KeyLeftAlt,
		RasterizePlane,
	)
	app.RegisterTool("help", helpTool)
	app.RegisterTool("select", selectTool)
	app.RegisterTool("voxel", voxelTool)
	app.RegisterTool("line", lineTool)
	app.RegisterTool("structure", structureTool)
	app.RegisterTool("shell", shellTool)
	app.RegisterTool("volume", volumeTool)
	app.RegisterTool("plane", planeTool)
	if e, v := os.ReadDir("."); e == nil {
		fmt.Printf("working dir : %v\n", v)
	}
	sh := []Collision{}
	app.Layer.OnChange = func(sc *Scene) {
		sh = app.Layer.GetShadows(&app.DirectionalLight)
	}
	app.Layer.ImportScene(filename)

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
		app.Mouse3 = app.Layer.GetIntersections(&camera)
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
			app.Mouse3Int = app.Mouse3.HitVoxel.Position
			app.Mouse3IntNext = rl.Vector3Add(app.Mouse3.HitVoxel.Position, app.Mouse3.Normal)
		case CollisionHitPlane, CollisionHitGuide:
			app.Mouse3IntNext = Vector3Round(app.Mouse3.Point)
			app.Mouse3Int = app.Mouse3IntNext
		}
		if rl.IsKeyPressed(rl.KeyLeftAlt) {
			app.Mouse3IntNext = app.Mouse3Int
		}

		cursor1.Position = app.Mouse3Int
		cursor2.Position = app.Mouse3IntNext
		app.Tools[app.CurrentTool].OnChanged(app.Tools[app.CurrentTool])

		if box_contains(&drawBox, &app.Mouse2) && rl.IsMouseButtonReleased(rl.MouseButtonLeft) && !rl.IsKeyDown(rl.KeyLeftControl) && !rl.IsKeyDown(rl.KeyLeftShift) {
			app.Tools[app.CurrentTool].Next(app.Mouse3Int)
		}
		orbit.ControlCamera()
		app.AddGuides(app.Mouse3Int)

		rl.UpdateCamera(&camera, rl.CameraCustom)

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.BeginMode3D(camera)

		DrawGridAtPoint(rl.NewVector3(0.5, 0.5, 0.5), 20, 1.0)
		app.Layer.Draw(AppRenderShaded, light)
		app.Guides.Draw(AppRenderWireframe, light)
		app.ConstructionHints.Draw(AppRenderBrighten, light)
		shadowColor := rl.Fade(rl.DarkGray, 0.3)
		// vv := rl.NewVector3(10, 10, 10)
		// vn := rl.NewVector3(0, 1, 0)
		///DrawDoubleSidedPlateWithNormal(&vv, vn, VOXEL_SZ, rl.Red)
		for _, s := range sh {
			// rl.DrawCubeV(s.Point, rl.Vector3CrossProduct(), shadowColor)
			DrawPlateWithNormal(&s.Point, s.Normal, VOXEL_SZ, shadowColor)
			// rl.DrawSphere(s.Point, .5, shadowColor)
		}
		if box_contains(&drawBox, &app.Mouse2) {
			cursor1.DrawWireframe(app.DirectionalLight, 1.2)
			cursor2.DrawWireframe(app.DirectionalLight, 1.2)
		}
		// You can then use Raylib functions to check for intersections, etc.
		if box_contains(&drawBox, &app.Mouse2) {

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
				rl.DrawText(name, int32(app.ScreenWidth)-64, int32(i)*64+48, 16, rl.NewColor(0, 0, 0, 255))
			}(index, toolName)
		}

		cu := int32(24)
		calcColor := func(chroma int32, lum int32) rl.Color {

			color := rl.ColorFromHSV(float32(chroma)*360/20, 1, 0.5)
			switch lum {
			case 0:
				v := uint8(256 * chroma / 21)
				color = rl.NewColor(v, v, v, 255)
			case 1, 2, 3:
				color = rl.ColorBrightness(color, float32(lum-2)*0.5)
			case 4:
				color = rl.Fade(color, 0.5)
			}
			return color
		}
		for chroma := int32(0); chroma <= 20; chroma += 1 {
			for lum := int32(0); lum <= 4; lum += 1 {
				color := calcColor(chroma, lum)
				rl.DrawRectangle(lum*cu, chroma*cu, cu, cu, color)
			}
		}
		chromaIndex := int32(math.Floor(float64(app.Mouse2.Y / float32(cu))))
		lumaIndex := int32(math.Floor(float64(app.Mouse2.X / float32(cu))))
		if box_contains(&leftPanelBox, &app.Mouse2) {
			rl.DrawRectangle(lumaIndex*cu, chromaIndex*cu, cu, cu, rl.NewColor(127, 255, 0, 127))
			if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
				color := calcColor(chromaIndex, lumaIndex)
				fmt.Printf("chosen color chroma-index:%d,lum-index:%d, result: %v", chromaIndex, lumaIndex, color)
				app.CurrentColor = color
				app.AddColor(rl.Vector2{X: float32(lumaIndex), Y: float32(chromaIndex)})
			}
		}
		for k := 0; k < len(app.HistoryColors); k++ {
			h := app.HistoryColors[k]
			rl.DrawRectangleLines(int32(h.X)*cu, int32(h.Y)*cu, cu, cu, rl.NewColor(127, 255, 0, 127))
		}

		// rl.DrawFPS(10, 10)

		app.Status = fmt.Sprintf("%s current tool : %s,CurrentColor: %v, scene: %d, guides: %d, helpers: %d", app.Status, app.CurrentTool, app.CurrentColor, len(app.Layer.Voxels), len(app.Guides.Voxels), len(app.ConstructionHints.Voxels))
		if isMousePositionChanged {
			cccc := fmt.Sprintf("%s %d", status, 2)
			doNothing(cccc)
		}

		rl.DrawText(app.Status, 10, app.ScreenHeight-20, 20, rl.Black)
		/// }
		rl.EndDrawing()
		// time.Sleep(1000000000)
	}
	app.Layer.ExportScene(filename)
	// app.Layer.ExportScene("temp.exp.vxdi")
	rl.CloseWindow()
}
