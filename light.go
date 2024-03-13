package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type VxdiDirectionalLight struct {
	Direction      rl.Vector3
	ShadowStrength float32
	LightStrength  float32
}
