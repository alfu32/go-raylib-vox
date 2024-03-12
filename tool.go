package main

import rl "github.com/gen2brain/raylib-go/raylib"

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
