package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type VxdiAppEditor struct {
	CurrentColor                 rl.Color
	HistoryColors                []rl.Vector2
	CurrentCameraProjection      rl.CameraProjection
	CurrentCameraProjectionIndex uint
	CameraProjections            [2]rl.CameraProjection
	CameraProjectionNames        [2]string
	CurrentCameraMode            rl.CameraMode
	CurrentCameraModeIndex       uint
	CameraModes                  [5]rl.CameraMode
	CameraModeNames              [5]string
	ScreenWidth                  int32
	ScreenHeight                 int32
	Layer                        *Scene
	Guides                       *Scene // Assume Scene is defined elsewhere
	ConstructionHints            *Scene // Assume Scene is defined elsewhere
	DirectionalLight             VxdiDirectionalLight
	CurrentTool                  string
	Tools                        map[string]*Tool
	ToolNames                    []string
	Status                       string
	Mouse2                       rl.Vector2
	Mouse3                       Collision
	Mouse3Int                    rl.Vector3
	Mouse3IntNext                rl.Vector3
}

func NewVxdiAppEditor(camera *rl.Camera3D, light VxdiDirectionalLight) *VxdiAppEditor {
	camera.Position = rl.NewVector3(10.0, 10.0, 10.0)
	camera.Target = rl.NewVector3(0, 0, 0)
	camera.Up = rl.NewVector3(0, 1, 0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective

	layer0 := NewScene(false, light)            // Assuming SceneInit is defined elsewhere
	guides := NewScene(false, light)            // Assuming SceneInit is defined elsewhere
	constructionHints := NewScene(false, light) // Assuming SceneInit is defined elsewhere

	app := VxdiAppEditor{
		CurrentCameraProjection:      rl.CameraPerspective,
		CurrentCameraProjectionIndex: 0,
		CameraProjections:            [2]rl.CameraProjection{rl.CameraPerspective, rl.CameraOrthographic},
		CameraProjectionNames:        [2]string{"PERSPECTIVE", "ORTHOGRAPHIC"},
		CurrentCameraMode:            rl.CameraFree,
		CurrentCameraModeIndex:       0,
		CameraModes:                  [5]rl.CameraMode{rl.CameraFree, rl.CameraOrbital, rl.CameraFirstPerson, rl.CameraThirdPerson},
		CameraModeNames:              [5]string{"FREE", "ORBITAL", "FIRST_PERSON", "THIRD_PERSON"},
		ScreenWidth:                  800,
		ScreenHeight:                 450,
		DirectionalLight:             light,
		CurrentTool:                  "select",
		Tools:                        make(map[string]*Tool, 0),
		ToolNames:                    make([]string, 0),
		Layer:                        layer0,
		Guides:                       guides,
		ConstructionHints:            constructionHints,
	}

	// Assume FillColorCircle and other necessary initialization here
	app.CurrentColor = rl.Red // Example, assuming color initialization
	app.AddColor(rl.Vector2{X: 2, Y: 0})

	return &app
}
func (app *VxdiAppEditor) AddColor(color rl.Vector2) {
	if len(app.HistoryColors) == 10 {
		for k := 0; k < 9; k++ {
			app.HistoryColors[k] = app.HistoryColors[k+1]
			app.HistoryColors[9] = color
		}
	} else {
		app.HistoryColors = append(app.HistoryColors, color)
	}
}
func (app *VxdiAppEditor) RegisterTool(name string, tool *Tool) {
	tool.Name = name
	app.Tools[name] = tool
	app.ToolNames = append(app.ToolNames, name)
}
func (app *VxdiAppEditor) RenderUI() {
}
