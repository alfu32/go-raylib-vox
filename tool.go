package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Tool struct {
	Name      string
	Points    []rl.Vector3
	Max       int
	Current   int
	OnFinish  func(tool *Tool)
	OnChanged func(tool *Tool)
}

func NewTool(max int, on_changed func(*Tool), on_finish func(*Tool)) *Tool {
	tool := &Tool{
		Points:    make([]rl.Vector3, max),
		Max:       max,
		Current:   0,
		OnFinish:  on_finish,
		OnChanged: on_changed,
	}
	return tool
}
func NewTool2Steps(app *VxdiAppEditor, modKey int, rasterizer_fn func(a, b rl.Vector3, operator VoxelOperatorFn)) *Tool {
	tool := NewTool(
		2,
		func(tool *Tool) {
			app.ConstructionHints.Clear()
			if tool.Current == 0 {
				fmt.Printf("tool progress : got point number [%d/%d]\n", tool.Current, tool.Max)
				app.ConstructionHints.AddVoxelAtPoint(&app.Mouse3Int, app.CurrentColor)
			} else if tool.Current == 1 {
				rasterizer_fn(tool.Points[0], app.Mouse3Int, func(point rl.Vector3) {
					app.ConstructionHints.AddVoxelAtPoint(&point, app.CurrentColor)
				})
			}
		},
		func(tool *Tool) {
			fmt.Printf("tool : got last point [%d/%d]\n", tool.Current, tool.Max)
			app.ConstructionHints.Clear()
			app.Guides.Clear()
			app.CurrentTool = app.ToolNames[0]
			if rl.IsKeyDown(rl.KeyLeftAlt) {
				fmt.Printf("tool : removing [%d/%d]\n", tool.Current, tool.Max)

				rasterizer_fn(tool.Points[0], app.Mouse3Int, func(point rl.Vector3) {
					app.Layer.RemoveVoxel(point.X, point.Y, point.Z)
				})
			} else {
				fmt.Printf("tool : adding [%d/%d]\n", tool.Current, tool.Max)

				rasterizer_fn(tool.Points[0], app.Mouse3Int, func(point rl.Vector3) {
					app.Layer.AddVoxelAtPoint(&point, app.CurrentColor)
				})
			}
			app.Layer.OnChange(app.Layer)
		},
	)
	return tool
}

func NewTool3Steps(app *VxdiAppEditor, modKey int, rasterizer_fn func(a, b, c rl.Vector3, operator VoxelOperatorFn)) *Tool {
	tool := NewTool(
		3,
		func(tool *Tool) {
			app.ConstructionHints.Clear()
			if tool.Current == 0 {
				fmt.Printf("tool progress : got point number [%d/%d]\n", tool.Current, tool.Max)
				app.ConstructionHints.AddVoxelAtPoint(&app.Mouse3Int, app.CurrentColor)
			} else if tool.Current == 1 {
				rasterizer_fn(tool.Points[0], app.Mouse3Int, app.Mouse3Int, func(point rl.Vector3) {
					app.ConstructionHints.AddVoxelAtPoint(&point, app.CurrentColor)
				})
			} else if tool.Current == 2 {
				rasterizer_fn(tool.Points[0], tool.Points[1], app.Mouse3Int, func(point rl.Vector3) {
					app.ConstructionHints.AddVoxelAtPoint(&point, app.CurrentColor)
				})
			}
		},
		func(tool *Tool) {
			fmt.Printf("tool : got last point [%d/%d]\n", tool.Current, tool.Max)
			app.ConstructionHints.Clear()
			app.Guides.Clear()
			app.CurrentTool = app.ToolNames[0]
			if rl.IsKeyDown(rl.KeyLeftAlt) {
				fmt.Printf("tool : removing [%d/%d]\n", tool.Current, tool.Max)

				rasterizer_fn(tool.Points[0], tool.Points[1], app.Mouse3Int, func(point rl.Vector3) {
					app.Layer.RemoveVoxel(point.X, point.Y, point.Z)
				})
			} else {
				fmt.Printf("tool : adding [%d/%d]\n", tool.Current, tool.Max)

				rasterizer_fn(tool.Points[0], tool.Points[1], app.Mouse3Int, func(point rl.Vector3) {
					app.Layer.AddVoxelAtPoint(&point, app.CurrentColor)
				})
			}
			app.Layer.OnChange(app.Layer)
		},
	)
	return tool
}
func (tool *Tool) Next(point rl.Vector3) {
	if tool.Current < tool.Max {
		tool.Points[tool.Current] = point
		tool.Current++
	}

	if tool.Current == tool.Max {
		tool.OnFinish(tool)
		tool.Reset()
		// Assume APP_CONSTRUCTION_MODE_VOXEL and scene__clear are defined and handled elsewhere
	} else {
		tool.OnChanged(tool)
	}
}
func (tool *Tool) Reset() {
	tool.Current = 0
}
