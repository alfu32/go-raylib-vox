package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// VxdiMultistepToolMutateVxdiAppEditorFn defines a function type for mutating the app editor.
type VxdiMultistepToolMutateVxdiAppEditorFn func(tool *VxdiMultistepTool, app *VxdiAppEditor, moving_point rl.Vector3)

// VxdiMultistepTool represents a tool that operates in multiple steps.
type VxdiMultistepTool struct {
	OnPointAcquired func(tool *VxdiMultistepTool, app *VxdiAppEditor)
	OnFinish        func(tool *VxdiMultistepTool, app *VxdiAppEditor)
	LastInputIndex  int
	NumInputs       int
	Inputs          []rl.Vector3
}

// NewVxdiMultistepTool initializes a new VxdiMultistepTool with the specified number of inputs and callbacks.
func NewVxdiMultistepTool(numInputs int, onPointAcquired, onFinish func(tool *VxdiMultistepTool, app *VxdiAppEditor)) *VxdiMultistepTool {
	return &VxdiMultistepTool{
		NumInputs:       numInputs,
		Inputs:          make([]rl.Vector3, numInputs),
		OnPointAcquired: onPointAcquired,
		OnFinish:        onFinish,
	}
}
func NewVxdiMultistepTool2Points(modKey int, progressFn func(a, b rl.Vector3), onFn func(a, b rl.Vector3), offFn func(a, b rl.Vector3)) *VxdiMultistepTool {
	tool := VxdiMultistepTool{
		NumInputs: 2,
		Inputs:    make([]rl.Vector3, 2),
		OnPointAcquired: func(tool *VxdiMultistepTool, app *VxdiAppEditor) {
			fmt.Printf("volume tool : got point number [%d/%d]\n", tool.LastInputIndex, tool.NumInputs)
			app.ConstructionHints.Clear()
			if tool.NumInputs == 1 {
				app.ConstructionHints.AddVoxelAtPoint(&tool.Inputs[0], app.CurrentColor)
			} else if tool.NumInputs == 2 {
				progressFn(tool.Inputs[0], app.Mouse3Int)
			}
		},
		OnFinish: func(tool *VxdiMultistepTool, app *VxdiAppEditor) {
			fmt.Printf("tool : got last point [%d/%d]\n", tool.LastInputIndex, tool.NumInputs)
			if rl.IsKeyDown(rl.KeyLeftAlt) {
				fmt.Printf("tool : removing [%d/%d]\n", tool.LastInputIndex, tool.NumInputs)
				onFn(tool.Inputs[0], tool.Inputs[1])
			} else {
				fmt.Printf("tool : adding [%d/%d]\n", tool.LastInputIndex, tool.NumInputs)
				offFn(tool.Inputs[0], tool.Inputs[1])
			}

		},
	}
	return &tool
}

// Reset reinitializes the tool for a new operation.
func (tool *VxdiMultistepTool) Reset() {
	tool.LastInputIndex = 0
}

// ReceivePoint processes a new point received by the tool.
func (tool *VxdiMultistepTool) ReceivePoint(app *VxdiAppEditor, scene *Scene, point rl.Vector3) {
	if tool.LastInputIndex < tool.NumInputs {
		tool.Inputs[tool.LastInputIndex] = point
		tool.LastInputIndex++
	}

	if tool.LastInputIndex == tool.NumInputs {
		tool.OnFinish(tool, app)
		tool.Reset()
		// Assume APP_CONSTRUCTION_MODE_VOXEL and scene__clear are defined and handled elsewhere
	}
}

// ReceiveMovingPoint processes a moving point, useful for visual feedback before finalizing.
func (tool *VxdiMultistepTool) ReceiveMovingPoint(app *VxdiAppEditor, point rl.Vector3) {
	if tool.LastInputIndex < tool.NumInputs {
		tool.Inputs[tool.LastInputIndex] = point

		if tool.LastInputIndex == tool.NumInputs-1 {
			tool.OnFinish(tool, app)
		} else {
			tool.OnPointAcquired(tool, app)
		}
	}
}

func (tool *VxdiMultistepTool) FactoryAquireFunction(before_complete func(), on_complete func()) {

}
