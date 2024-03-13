package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const VOXEL_SZ = float32(1)
const VOXEL_SZ2 = float32(0.5)

type Voxel struct {
	Position rl.Vector3
	Material rl.Color
}

func NewVoxel() *Voxel {
	return &Voxel{
		Position: rl.NewVector3(0, 0, 0),
		Material: rl.NewColor(1, 1, 1, 255),
	}
}
func InitVoxel(x float32, y float32, z float32, color rl.Color) *Voxel {
	return &Voxel{
		Position: rl.NewVector3(x, y, z),
		Material: color,
	}
}
func KeyForVoxel(pos rl.Vector3) string {
	return fmt.Sprintf("%f,%f,%f", pos.X, pos.Y, pos.Z)
}
func (voxel *Voxel) GetBoundingBox() rl.BoundingBox {
	pos := Vector3Round(voxel.Position)
	return rl.BoundingBox{
		Min: rl.NewVector3(pos.X-VOXEL_SZ2, pos.Y-VOXEL_SZ2, pos.Z-VOXEL_SZ2),
		Max: rl.NewVector3(pos.X+VOXEL_SZ2, pos.Y+VOXEL_SZ2, pos.Z+VOXEL_SZ2),
	}
}

// DrawTriangleStrip3D draws a triangle strip defined by points in 3D space.
// This function is a conceptual translation and may require custom engine modifications to work as intended.
func DrawTriangleStrip3D(points []rl.Vector3, color rl.Color) {
	if len(points) < 3 {
		return
	}

	// Since Raylib Go doesn't expose rlBegin, rlVertex3f, and rlEnd directly,
	// you might need to use a combination of existing functions or modify the Raylib Go bindings.
	// This is a placeholder loop to demonstrate the iteration over points, adapted to Go's range syntax.
	for i := 2; i < len(points); i++ {
		if i%2 == 0 {
			// Placeholder for drawing the triangle
			// You might need to implement a custom function using lower-level OpenGL calls or extend Raylib Go
			drawTriangle(points[i], points[i-2], points[i-1], color)
		} else {
			// Placeholder for drawing the triangle
			drawTriangle(points[i], points[i-1], points[i-2], color)
		}
	}
}

// drawTriangle is a placeholder function. You'll need to replace this with actual drawing logic,
// possibly by extending Raylib Go or using OpenGL calls directly.
func drawTriangle(p1, p2, p3 rl.Vector3, color rl.Color) {
	// Implement the triangle drawing here. This could involve using raylib.DrawLine3D for edges,
	// or more likely, making custom OpenGL calls if you need to fill the triangles.
	rl.DrawTriangle3D(p1, p2, p3, color)
}
func (voxel *Voxel) DrawWireframe(light VxdiDirectionalLight, sz float32) {
	rl.DrawCubeWires(voxel.Position, sz, sz, sz, voxel.Material)
}
func (voxel *Voxel) DrawCube(light VxdiDirectionalLight, sz float32) {
	rl.DrawCube(voxel.Position, sz, sz, sz, voxel.Material)
}
func DrawDoubleSidedPlateWithNormal(center *rl.Vector3, normal rl.Vector3, size float32, color rl.Color) {
	p1 := rl.Vector3Add(*center, rl.Vector3Scale(normal, 0.1))
	n2 := rl.NewVector3(normal.X*-1, normal.Y*-1, normal.Z*-1)
	p2 := rl.Vector3Add(*center, rl.Vector3Scale(n2, 0.1))
	DrawPlateWithNormal(&p1, normal, size, color)
	DrawPlateWithNormal(&p2, normal, size, color)
}
func Vector3Dominant(v *rl.Vector3) (rl.Vector3, int32) {
	x := float64(v.X)
	y := float64(v.Y)
	z := float64(v.Z)
	if math.Abs(x) > math.Abs(y) && math.Abs(x) > math.Abs(z) {
		if v.X > 0 {
			return rl.NewVector3(1, 0, 0), 0
		} else {
			return rl.NewVector3(-1, 0, 0), 1
		}
	} else if math.Abs(y) > math.Abs(x) && math.Abs(y) > math.Abs(z) {
		if v.Y > 0 {
			return rl.NewVector3(0, 1, 0), 2
		} else {
			return rl.NewVector3(0, -1, 0), 3
		}
	} else {
		if v.Z > 0 {
			return rl.NewVector3(0, 0, 1), 4
		} else {
			return rl.NewVector3(0, 0, -1), 5
		}
	}
}

func DrawPlateWithNormal(pos *rl.Vector3, normal rl.Vector3, size float32, color rl.Color) {
	revval := func(v float32) float32 {
		if math.Abs(float64(v)) < 0.01 {
			if v < 0 {
				return -1
			} else {
				return 1
			}
		} else {
			if v < 0 {
				return -0.02
			} else {
				return 0.02
			}
		}
	}
	// Normalize the normal vector
	n := rl.Vector3Normalize(normal)
	dom, _ := Vector3Dominant(&n)
	szx := revval(dom.X)*size/2 - 0.005
	szy := revval(dom.Y)*size/2 - 0.005
	szz := revval(dom.Z)*size/2 - 0.005
	//pos0 := Vector3Round(*pos)
	//pos = &pos0

	// Define vertices for each face of the cube
	vertices := []rl.Vector3{
		// Front face
		{X: pos.X - szx, Y: pos.Y - szy, Z: pos.Z + szz},
		{X: pos.X + szx, Y: pos.Y - szy, Z: pos.Z + szz},
		{X: pos.X - szx, Y: pos.Y + szy, Z: pos.Z + szz},
		{X: pos.X + szx, Y: pos.Y + szy, Z: pos.Z + szz},
		// Back face
		{X: pos.X - szx, Y: pos.Y - szy, Z: pos.Z - szz},
		{X: pos.X - szx, Y: pos.Y + szy, Z: pos.Z - szz},
		{X: pos.X + szx, Y: pos.Y - szy, Z: pos.Z - szz},
		{X: pos.X + szx, Y: pos.Y + szy, Z: pos.Z - szz},
		// Top face
		{X: pos.X - szx, Y: pos.Y + szy, Z: pos.Z + szz},
		{X: pos.X + szx, Y: pos.Y + szy, Z: pos.Z + szz},
		{X: pos.X - szx, Y: pos.Y + szy, Z: pos.Z - szz},
		{X: pos.X + szx, Y: pos.Y + szy, Z: pos.Z - szz},
		// Bottom face
		{X: pos.X - szx, Y: pos.Y - szy, Z: pos.Z - szz},
		{X: pos.X + szx, Y: pos.Y - szy, Z: pos.Z - szz},
		{X: pos.X - szx, Y: pos.Y - szy, Z: pos.Z + szz},
		{X: pos.X + szx, Y: pos.Y - szy, Z: pos.Z + szz},
		// Left face
		{X: pos.X - szx, Y: pos.Y - szy, Z: pos.Z - szz},
		{X: pos.X - szx, Y: pos.Y - szy, Z: pos.Z + szz},
		{X: pos.X - szx, Y: pos.Y + szy, Z: pos.Z - szz},
		{X: pos.X - szx, Y: pos.Y + szy, Z: pos.Z + szz},
		// Right face
		{X: pos.X + szx, Y: pos.Y - szy, Z: pos.Z - szz},
		{X: pos.X + szx, Y: pos.Y + szy, Z: pos.Z - szz},
		{X: pos.X + szx, Y: pos.Y - szy, Z: pos.Z + szz},
		{X: pos.X + szx, Y: pos.Y + szy, Z: pos.Z + szz},
	}

	// Draw the cube using triangle strips with shaded colors
	DrawTriangleStrip3D(vertices[0:4], color)   //voxel.Paterial_color);// Front face
	DrawTriangleStrip3D(vertices[4:8], color)   //voxel.Paterial_color);// Back face
	DrawTriangleStrip3D(vertices[8:12], color)  // Top face
	DrawTriangleStrip3D(vertices[12:16], color) // Bottom face
	DrawTriangleStrip3D(vertices[16:20], color) // Left face
	DrawTriangleStrip3D(vertices[20:24], color) // Right face
	// Draw each face of the cube with shaded colors
}
func (voxel *Voxel) DrawShaded(light VxdiDirectionalLight, sz float32) {
	sz2 := sz / 2
	n := rl.Vector3Normalize(light.Direction)
	pos := Vector3Round(voxel.Position)

	// Define vertices for each face of the cube
	vertices := []rl.Vector3{
		// Front face
		{X: pos.X - sz2, Y: pos.Y - sz2, Z: pos.Z + sz2},
		{X: pos.X + sz2, Y: pos.Y - sz2, Z: pos.Z + sz2},
		{X: pos.X - sz2, Y: pos.Y + sz2, Z: pos.Z + sz2},
		{X: pos.X + sz2, Y: pos.Y + sz2, Z: pos.Z + sz2},
		// Back face
		{X: pos.X - sz2, Y: pos.Y - sz2, Z: pos.Z - sz2},
		{X: pos.X - sz2, Y: pos.Y + sz2, Z: pos.Z - sz2},
		{X: pos.X + sz2, Y: pos.Y - sz2, Z: pos.Z - sz2},
		{X: pos.X + sz2, Y: pos.Y + sz2, Z: pos.Z - sz2},
		// Top face
		{X: pos.X - sz2, Y: pos.Y + sz2, Z: pos.Z + sz2},
		{X: pos.X + sz2, Y: pos.Y + sz2, Z: pos.Z + sz2},
		{X: pos.X - sz2, Y: pos.Y + sz2, Z: pos.Z - sz2},
		{X: pos.X + sz2, Y: pos.Y + sz2, Z: pos.Z - sz2},
		// Bottom face
		{X: pos.X - sz2, Y: pos.Y - sz2, Z: pos.Z - sz2},
		{X: pos.X + sz2, Y: pos.Y - sz2, Z: pos.Z - sz2},
		{X: pos.X - sz2, Y: pos.Y - sz2, Z: pos.Z + sz2},
		{X: pos.X + sz2, Y: pos.Y - sz2, Z: pos.Z + sz2},
		// Left face
		{X: pos.X - sz2, Y: pos.Y - sz2, Z: pos.Z - sz2},
		{X: pos.X - sz2, Y: pos.Y - sz2, Z: pos.Z + sz2},
		{X: pos.X - sz2, Y: pos.Y + sz2, Z: pos.Z - sz2},
		{X: pos.X - sz2, Y: pos.Y + sz2, Z: pos.Z + sz2},
		// Right face
		{X: pos.X + sz2, Y: pos.Y - sz2, Z: pos.Z - sz2},
		{X: pos.X + sz2, Y: pos.Y + sz2, Z: pos.Z - sz2},
		{X: pos.X + sz2, Y: pos.Y - sz2, Z: pos.Z + sz2},
		{X: pos.X + sz2, Y: pos.Y + sz2, Z: pos.Z + sz2},
	}

	// Draw the cube using triangle strips with shaded colors
	DrawTriangleStrip3D(vertices[0:4], rl.ColorBrightness(voxel.Material, -n.X*light.ShadowStrength))   //voxel.Paterial_color);// Front face
	DrawTriangleStrip3D(vertices[4:8], rl.ColorBrightness(voxel.Material, n.X*light.LightStrength))     //voxel.Paterial_color);// Back face
	DrawTriangleStrip3D(vertices[8:12], rl.ColorBrightness(voxel.Material, -n.Y*light.ShadowStrength))  // Top face
	DrawTriangleStrip3D(vertices[12:16], rl.ColorBrightness(voxel.Material, n.Y*light.LightStrength))   // Bottom face
	DrawTriangleStrip3D(vertices[16:20], rl.ColorBrightness(voxel.Material, -n.Z*light.ShadowStrength)) // Left face
	DrawTriangleStrip3D(vertices[20:24], rl.ColorBrightness(voxel.Material, n.Z*light.LightStrength))   // Right face
	// Draw each face of the cube with shaded colors
}

// Draw the cube using triangle strips with shaded colors
// The actual drawing functions need to be adapted to your environment and the capabilities of the Raylib Go bindings.
// Raylib Go does not directly support DrawTriangleStrip3D or ColorBrightness functions.
// You would need to implement similar functionality in Go, possibly using other Raylib functions
// such as rl.DrawTriangle3D or custom shader logic to achieve the shading effects based on light direction and strength.

func (a *Voxel) Equals(b *Voxel) bool {
	return math.Abs(float64(a.Position.X)-float64(a.Position.X)) < 1 &&
		math.Abs(float64(a.Position.Y)-float64(a.Position.Y)) < 1 &&
		math.Abs(float64(a.Position.Z)-float64(a.Position.Z)) < 1
}
