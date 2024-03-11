package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Collision struct {
	Hit          bool
	Distance     float32
	Point        rl.Vector3
	Normal       rl.Vector3
	Voxel        Voxel
	VoxelIndex   string
	CollisionHit CollisionHit
}

type CollisionHit int

const (
	CollisionHitNone CollisionHit = iota
	CollisionHitVoxel
	CollisionHitPlane
	CollisionHitGuide
	// Add other collision hit types as needed
)

var kMaxPoint = rl.NewVector3(math.MaxFloat32, math.MaxFloat32, math.MaxFloat32)
var kZeroPoint = rl.NewVector3(0, 0, 0)
var kVoxelNone = Voxel{} // Assuming a default or "none" voxel state can be represented
