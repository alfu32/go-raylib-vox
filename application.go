package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type AppConstructionMode int

const (
	AppConstructionModeHelp AppConstructionMode = iota + 0x100
	AppConstructionModeSelect
	AppConstructionModeVoxel
	AppConstructionModeLine
	AppConstructionModeStructure
	AppConstructionModeShell
	AppConstructionModeVolume
	AppConstructionModePlate
)

type VxdiAppEditor struct {
	ConstructionMode             AppConstructionMode
	CurrentColor                 rl.Color
	CurrentColorIndex            uint
	Colors                       [360]rl.Color
	CurrentCameraProjection      rl.CameraProjection
	CurrentCameraProjectionIndex uint
	CameraProjections            [2]rl.CameraProjection
	CameraProjectionNames        [2]string
	CurrentCameraMode            rl.CameraMode
	CurrentCameraModeIndex       uint
	CameraModes                  [5]rl.CameraMode
	CameraModeNames              [5]string
	ScreenWidth                  int
	ScreenHeight                 int
	Guides                       Scene // Assume Scene is defined elsewhere
	ConstructionHints            Scene // Assume Scene is defined elsewhere
	LightDirection               VxdiLight
	TextBuffer                   string
}

func NewVxdiAppEditor(camera *rl.Camera3D, light VxdiLight) VxdiAppEditor {
	camera.Position = rl.NewVector3(10.0, 10.0, 10.0)
	camera.Target = rl.NewVector3(0, 0, 0)
	camera.Up = rl.NewVector3(0, 1, 0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective

	guides := NewScene(false, light)            // Assuming SceneInit is defined elsewhere
	constructionHints := NewScene(false, light) // Assuming SceneInit is defined elsewhere

	app := VxdiAppEditor{
		ConstructionMode:             AppConstructionModeVoxel,
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
		LightDirection:               light,
		Guides:                       guides,
		ConstructionHints:            constructionHints,
	}

	// Assume FillColorCircle and other necessary initialization here
	app.CurrentColor = app.Colors[4] // Example, assuming color initialization
	app.CurrentColorIndex = 4

	return app
}
