package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// VoxelOperatorFn defines a function type in Go that matches the C typedef for the voxel operator function.
type VoxelOperatorFn func(position rl.Vector3)
type Builder2Points func(a, b rl.Vector3, scene *Scene, material rl.Color, materialID uint, operator VoxelOperatorFn)

// rasterizeLine function to rasterize a line between two 3D points in Go.
func RasterizeLine(a, b rl.Vector3, operator VoxelOperatorFn) {
	dx := b.X - a.X
	dy := b.Y - a.Y
	dz := b.Z - a.Z
	ab := rl.NewVector3(dx, dy, dz)

	cx0 := min(int(a.X), int(b.X))
	cy0 := min(int(a.Y), int(b.Y))
	cz0 := min(int(a.Z), int(b.Z))
	cx1 := max(int(a.X), int(b.X))
	cy1 := max(int(a.Y), int(b.Y))
	cz1 := max(int(a.Z), int(b.Z))

	for x := cx0; x <= cx1; x++ {
		for y := cy0; y <= cy1; y++ {
			for z := cz0; z <= cz1; z++ {
				p := rl.NewVector3(float32(x), float32(y), float32(z))

				// Calculate vector from a to p
				ap := rl.NewVector3(p.X-a.X, p.Y-a.Y, p.Z-a.Z)

				// Calculate dot product of ab and ap
				dotABAP := ab.X*ap.X + ab.Y*ap.Y + ab.Z*ap.Z

				// Calculate dot product of ab with itself
				dotABAB := ab.X*ab.X + ab.Y*ab.Y + ab.Z*ab.Z

				// Calculate t value (projection factor)
				t := dotABAP / dotABAB

				// Calculate the point on the line
				p0 := rl.NewVector3(a.X+t*ab.X, a.Y+t*ab.Y, a.Z+t*ab.Z)
				if rl.Vector3Distance(p0, p) <= 0.867 {
					operator(p)
				}
			}
		}
	}
}

// RasterizeSolidCube function to rasterize a solid cube given two 3D points.
func RasterizeSolidCube(a, b rl.Vector3, operator VoxelOperatorFn) {
	// Calculate the differences between points
	// dx := b.X - a.X
	// dy := b.Y - a.Y
	// dz := b.Z - a.Z
	// ab := rl.NewVector3(dx, dy, dz)

	cx0 := min(int(a.X), int(b.X))
	cy0 := min(int(a.Y), int(b.Y))
	cz0 := min(int(a.Z), int(b.Z))
	cx1 := max(int(a.X), int(b.X))
	cy1 := max(int(a.Y), int(b.Y))
	cz1 := max(int(a.Z), int(b.Z))

	for x := cx0; x <= cx1; x++ {
		for y := cy0; y <= cy1; y++ {
			for z := cz0; z <= cz1; z++ {
				p := rl.NewVector3(float32(x), float32(y), float32(z))
				operator(p)
			}
		}
	}
}
func RasterizeHollowCube(a, b rl.Vector3, operator VoxelOperatorFn) {
	RasterizeSolidCube(a, rl.NewVector3(a.X, b.Y, b.Z), operator)
	RasterizeSolidCube(rl.NewVector3(b.X, a.Y, a.Z), b, operator)

	RasterizeSolidCube(a, rl.NewVector3(b.X, a.Y, b.Z), operator)
	RasterizeSolidCube(rl.NewVector3(a.X, b.Y, a.Z), b, operator)

	RasterizeSolidCube(a, rl.NewVector3(b.X, b.Y, a.Z), operator)
	RasterizeSolidCube(rl.NewVector3(a.X, a.Y, b.Z), b, operator)
}

// RasterizeStructureCube rasterizes a structured cube by drawing its skeletal structure.
func RasterizeStructureCube(a, b rl.Vector3, operator VoxelOperatorFn) {
	RasterizeSolidCube(a, rl.NewVector3(b.X, a.Y, a.Z), operator)
	RasterizeSolidCube(a, rl.NewVector3(a.X, b.Y, a.Z), operator)
	RasterizeSolidCube(rl.NewVector3(b.X, a.Y, a.Z), rl.NewVector3(b.X, b.Y, a.Z), operator)
	RasterizeSolidCube(rl.NewVector3(a.X, b.Y, a.Z), rl.NewVector3(b.X, b.Y, a.Z), operator)

	// Generate edges for the opposite face
	RasterizeSolidCube(rl.NewVector3(a.X, a.Y, b.Z), rl.NewVector3(b.X, a.Y, b.Z), operator)
	RasterizeSolidCube(rl.NewVector3(a.X, a.Y, b.Z), rl.NewVector3(a.X, b.Y, b.Z), operator)
	RasterizeSolidCube(rl.NewVector3(b.X, a.Y, b.Z), b, operator)
	RasterizeSolidCube(rl.NewVector3(a.X, b.Y, b.Z), b, operator)

	// Generate edges connecting both faces
	RasterizeSolidCube(a, rl.NewVector3(a.X, a.Y, b.Z), operator)
	RasterizeSolidCube(rl.NewVector3(b.X, a.Y, a.Z), rl.NewVector3(b.X, a.Y, b.Z), operator)
	RasterizeSolidCube(rl.NewVector3(a.X, b.Y, a.Z), rl.NewVector3(a.X, b.Y, b.Z), operator)
	RasterizeSolidCube(rl.NewVector3(b.X, b.Y, a.Z), b, operator)
}
