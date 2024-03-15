package main

import rl "github.com/gen2brain/raylib-go/raylib"

func (scene *Scene) TraceVoxel(vx *Voxel, light VxdiDirectionalLight) Collision {
	ray := rl.NewRay(vx.Position, light.Direction)
	result := scene.IntersectPoint(ray)
	return result
}
