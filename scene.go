package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Scene struct {
	Voxels      map[string]*Voxel
	Light       VxdiLight
	IsPersisted bool
	Filename    string
}

// NewScene creates and returns a new instance of a Scene.
func NewScene(is_persisted bool, light VxdiLight) *Scene {
	return &Scene{
		Voxels:      make(map[string]*Voxel),
		Light:       light,
		IsPersisted: is_persisted,
		Filename:    "temp",
	}
}

// AddVoxel adds a new voxel to the scene. If a voxel already exists at the given coordinates, it updates the existing voxel.
func (s *Scene) Clear() {
	s.Voxels = make(map[string]*Voxel)
}

// AddVoxel adds a new voxel to the scene. If a voxel already exists at the given coordinates, it updates the existing voxel.
func (s *Scene) AddVoxel(v *Voxel) {
	key := fmt.Sprintf("%d,%d,%d", int32(v.Position.X), int32(v.Position.Y), int32(v.Position.Z))
	fmt.Printf("Adding Voxel[%s] %v, len %d\n", key, v, len(s.Voxels))
	s.Voxels[key] = v
	fmt.Printf("Added Voxel %v, new len %d\n", v, len(s.Voxels))
}
func (s *Scene) AddVoxelAtPoint(p *rl.Vector3, mat rl.Color) {
	s.AddVoxel(InitVoxel(p.X, p.Y, p.Z, mat))
}

// RemoveVoxel removes a voxel from the scene by its coordinates. If no voxel exists at those coordinates, it does nothing.
func (s *Scene) RemoveVoxel(x, y, z float32) {
	key := fmt.Sprintf("%d,%d,%d", int32(x), int32(y), int32(z))
	delete(s.Voxels, key)
}

// GetVoxel retrieves a voxel from the scene by its coordinates. It returns the voxel and a boolean indicating if it was found.
func (s *Scene) GetVoxel(x, y, z float32) (*Voxel, bool) {
	key := fmt.Sprintf("%d,%d,%d", int32(x), int32(y), int32(z))
	voxel, exists := s.Voxels[key]
	return voxel, exists
}
func (scene *Scene) Draw(typ uint8, light VxdiLight) {
	for _, v := range scene.Voxels {
		switch typ {
		case 0: // type hint
			rl.DrawCubeWires(Vector3Round(v.Position), VOXEL_SZ, VOXEL_SZ, VOXEL_SZ, rl.DarkGray)
		case 1: // type objects
			v.DrawShaded(light, VOXEL_SZ)
			rl.DrawCubeWires(Vector3Round(v.Position), VOXEL_SZ, VOXEL_SZ, VOXEL_SZ, rl.DarkGray)
		case 2: // type guides
			rl.DrawCubeWires(Vector3Round(v.Position), VOXEL_SZ/4, VOXEL_SZ/4, VOXEL_SZ/4, rl.DarkGray)
		}

	}
}

// sceneRayIntersectPoint checks if a ray intersects any voxel in the scene and returns the nearest intersection if any.
func (scene *Scene) IntersectPoint(ray rl.Ray) Collision {
	result := Collision{
		Hit:          false,
		Distance:     math.MaxFloat32,
		Point:        kMaxPoint,
		Normal:       kZeroPoint,
		Voxel:        kVoxelNone,
		VoxelIndex:   "none",
		CollisionHit: CollisionHitNone,
	}

	for i, cube := range scene.Voxels { // Assuming scene.Voxels is a slice of Voxel
		boundingBox := cube.GetBoundingBox()
		collision := rl.GetRayCollisionBox(ray, boundingBox)
		if collision.Hit && collision.Distance < result.Distance {
			result = Collision{
				Hit:          collision.Hit,
				Distance:     collision.Distance,
				Point:        collision.Point,
				Normal:       collision.Normal,
				Voxel:        *cube,
				VoxelIndex:   i,
				CollisionHit: CollisionHitVoxel,
			}
		}
	}
	return result
}

// sceneGetIntersections checks for intersections between a ray (cast from the mouse position) and scene objects or a fallback plane.
func (scene *Scene) GetIntersections(camera *rl.Camera3D) Collision {
	ray := rl.GetMouseRay(rl.GetMousePosition(), *camera)
	result := scene.IntersectPoint(ray)

	if !result.Hit {
		// Fallback check against a plane
		fallbackBox := rl.NewBoundingBox(rl.NewVector3(-100.0, -0.1, -100.0), rl.NewVector3(100.0, 0.0, 100.0))
		collision := rl.GetRayCollisionBox(ray, fallbackBox)
		if collision.Hit {
			result = Collision{
				Hit:          false,
				Distance:     collision.Distance,
				Point:        collision.Point,
				Normal:       collision.Normal,
				Voxel:        kVoxelNone,
				VoxelIndex:   "none",
				CollisionHit: CollisionHitPlane,
			}
		}
	}
	return result
}

// Assume voxelGetBoundingBox is implemented elsewhere, returning an rl.BoundingBox for a given Voxel.

// Additional methods for Scene

// ExportScene writes the current state of the scene to a text file.
func (s *Scene) ExportScene(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, voxel := range s.Voxels {
		line := fmt.Sprintf("%f,%f,%f,%d,%d,%d,%d\n", voxel.Position.X, voxel.Position.Y, voxel.Position.Z, voxel.Material.R, voxel.Material.G, voxel.Material.B, voxel.Material.A)
		_, err := writer.WriteString(line)
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

// ImportScene reads a scene from a text file and replaces the current scene's content.
func (s *Scene) ImportScene(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	s.Voxels = make(map[string]*Voxel) // Reset current scene
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 4 {
			continue // Skip malformed lines
		}
		x, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			continue
		}
		y, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			continue
		}
		z, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			continue
		}
		r, err := strconv.ParseUint(parts[3], 10, 8)
		if err != nil {
			continue
		}
		g, err := strconv.ParseUint(parts[4], 10, 8)
		if err != nil {
			continue
		}
		b, err := strconv.ParseUint(parts[5], 10, 8)
		if err != nil {
			continue
		}
		a, err := strconv.ParseUint(parts[6], 10, 8)
		if err != nil {
			continue
		}

		s.AddVoxel(
			&Voxel{
				Position: rl.Vector3{X: float32(x), Y: float32(y), Z: float32(z)},
				Material: rl.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)},
			},
		)
	}
	return scanner.Err()
}

// SaveScene saves the current state of the scene to a binary file using gob encoding.
func (s *Scene) SaveScene(filename string) error {
	fmt.Printf("SaveScene bin %s \n", filename)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	encoder := gob.NewEncoder(writer)
	err = encoder.Encode(s.Voxels)
	if err != nil {
		return err
	}
	fmt.Printf("flushing scene bin %s \n", filename)
	return writer.Flush()
}

// LoadScene loads a scene from a binary file using gob encoding and replaces the current scene's content.
func (s *Scene) LoadScene(filename string) error {
	fmt.Printf("LoadScene bin %s \n", filename)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := gob.NewDecoder(reader)

	// Assuming the structure of Voxels map is the same as when it was saved
	voxels := make(map[string]*Voxel)
	err = decoder.Decode(&voxels)
	if err != nil {
		return err
	}

	s.Voxels = voxels // Replace current scene with the loaded one
	return nil
}
