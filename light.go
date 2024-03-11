package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type VxdiLight struct {
	Direction      rl.Vector3
	ShadowStrength float32
	LightStrength  float32
}
