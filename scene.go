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

type AppRenderType int

const (
	AppRenderDefault = AppRenderType(0b10000000) << iota
	AppRenderBrighten
	AppRenderDarken
	AppRenderShaded
	AppRenderTransparent
	AppRenderWireframe
)

type Scene struct {
	Voxels      map[string]*Voxel
	keys        []string
	Light       VxdiDirectionalLight
	IsPersisted bool
	Filename    string
	OnChange    func(sc *Scene)
}

// NewScene creates and returns a new instance of a Scene.
func NewScene(is_persisted bool, light VxdiDirectionalLight) *Scene {
	return &Scene{
		Voxels:      make(map[string]*Voxel),
		keys:        make([]string, 0),
		Light:       light,
		IsPersisted: is_persisted,
		Filename:    "temp",
		OnChange:    func(sc *Scene) {},
	}
}

// AddVoxel adds a new voxel to the scene. If a voxel already exists at the given coordinates, it updates the existing voxel.
func (s *Scene) Clear() {
	s.Voxels = make(map[string]*Voxel)
	s.OnChange(s)
}

// AddVoxel adds a new voxel to the scene. If a voxel already exists at the given coordinates, it updates the existing voxel.
func (s *Scene) AddVoxel(v *Voxel) {
	key := fmt.Sprintf("%d,%d,%d", int32(v.Position.X), int32(v.Position.Y), int32(v.Position.Z))
	// fmt.Printf("Adding Voxel[%s] %v, len %d\n", key, v, len(s.Voxels))
	if _, exists := s.Voxels[key]; !exists {
		s.Voxels[key] = v
		s.keys = append(s.keys, key)
		//s.OnChange(s)
	}

	// fmt.Printf("Added Voxel %v, new len %d\n", v, len(s.Voxels))
}
func (s *Scene) AddVoxelAtPoint(p *rl.Vector3, mat rl.Color) {
	s.AddVoxel(InitVoxel(p.X, p.Y, p.Z, mat))
}
func remove[T comparable](l []T, item T) []T {
	out := make([]T, 0)
	for _, element := range l {
		if element != item {
			out = append(out, element)
		}
	}
	return out
}

// RemoveVoxel removes a voxel from the scene by its coordinates. If no voxel exists at those coordinates, it does nothing.
func (s *Scene) RemoveVoxel(x, y, z float32) {
	key := fmt.Sprintf("%d,%d,%d", int32(x), int32(y), int32(z))
	if _, exists := s.Voxels[key]; exists {
		delete(s.Voxels, key)
		s.keys = remove[string](s.keys, key)
		//s.OnChange(s)
	}

}

// GetVoxel retrieves a voxel from the scene by its coordinates. It returns the voxel and a boolean indicating if it was found.
func (s *Scene) GetVoxel(x, y, z float32) (*Voxel, bool) {
	key := fmt.Sprintf("%d,%d,%d", int32(x), int32(y), int32(z))
	voxel, exists := s.Voxels[key]
	return voxel, exists
}

func (scene *Scene) Draw(typ AppRenderType, light VxdiDirectionalLight) {
	for _, k := range scene.keys {
		if v, ok := scene.Voxels[k]; ok {
			switch typ {
			case AppRenderDarken: // type hint
				rl.DrawCube(Vector3Round(v.Position), VOXEL_SZ, VOXEL_SZ, VOXEL_SZ, rl.ColorBrightness(v.Material, -0.3))
				rl.DrawCubeWires(Vector3Round(v.Position), VOXEL_SZ, VOXEL_SZ, VOXEL_SZ, v.Material)
			case AppRenderBrighten: // type hint
				rl.DrawCube(Vector3Round(v.Position), VOXEL_SZ, VOXEL_SZ, VOXEL_SZ, rl.ColorBrightness(v.Material, 0.3))
				rl.DrawCubeWires(Vector3Round(v.Position), VOXEL_SZ, VOXEL_SZ, VOXEL_SZ, v.Material)
			case AppRenderTransparent: // type hint
				rl.DrawCube(Vector3Round(v.Position), VOXEL_SZ, VOXEL_SZ, VOXEL_SZ, rl.Fade(v.Material, 0.3))
				rl.DrawCubeWires(Vector3Round(v.Position), VOXEL_SZ, VOXEL_SZ, VOXEL_SZ, v.Material)
			case AppRenderShaded: // type objects
				v.DrawShaded(light, VOXEL_SZ)
				rl.DrawCubeWires(Vector3Round(v.Position), VOXEL_SZ, VOXEL_SZ, VOXEL_SZ, rl.DarkGray)
			case AppRenderWireframe: // type guides
				rl.DrawCubeWires(Vector3Round(v.Position), VOXEL_SZ/4, VOXEL_SZ/4, VOXEL_SZ/4, v.Material)
			}
		}
	}
}

// sceneRayIntersectPoint checks if a ray intersects any voxel in the scene and returns the nearest intersection if any.
func (scene *Scene) IntersectPoint(ray rl.Ray) Collision {
	result := Collision{
		Hit:           false,
		Distance:      math.MaxFloat32,
		Point:         kMaxPoint,
		Normal:        kZeroPoint,
		HitVoxel:      kVoxelNone,
		HitVoxelIndex: "none",
		CollisionHit:  CollisionHitNone,
	}

	for i, cube := range scene.Voxels { // Assuming scene.Voxels is a slice of Voxel
		boundingBox := cube.GetBoundingBox()
		collision := rl.GetRayCollisionBox(ray, boundingBox)
		if collision.Hit && collision.Distance < result.Distance {
			result = Collision{
				Hit:           collision.Hit,
				Distance:      collision.Distance,
				Point:         collision.Point,
				Normal:        collision.Normal,
				HitVoxel:      *cube,
				HitVoxelIndex: i,
				CollisionHit:  CollisionHitVoxel,
			}
		}
	}
	return result
}

// sceneRayIntersectPoint checks if a ray intersects any voxel in the scene and returns the nearest intersection if any.
func (scene *Scene) IntersectPoints(ray rl.Ray) []Collision {
	result := []Collision{}
	result0 := []Collision{}

	for i, voxel := range scene.Voxels { // Assuming scene.Voxels is a slice of Voxel
		boundingBox := voxel.GetBoundingBox()
		collision := rl.GetRayCollisionBox(ray, boundingBox)
		if collision.Hit {
			result = append(result, Collision{
				Hit:           collision.Hit,
				Distance:      collision.Distance,
				Point:         collision.Point,
				Normal:        collision.Normal,
				HitVoxel:      *voxel,
				HitVoxelIndex: i,
				CollisionHit:  CollisionHitVoxel,
			})
		}
	}

	// Fallback check against a plane
	fallbackBox := rl.NewBoundingBox(rl.NewVector3(-200.0, 0.4, -200.0), rl.NewVector3(200.0, 0.5, 200.0))
	collision := rl.GetRayCollisionBox(ray, fallbackBox)
	result = append(result, Collision{
		Hit:           false,
		Distance:      collision.Distance,
		Point:         collision.Point,
		Normal:        collision.Normal,
		HitVoxel:      kVoxelNone,
		HitVoxelIndex: "none",
		CollisionHit:  CollisionHitPlane,
	})
	for _, c := range result {
		if c.CollisionHit == CollisionHitVoxel && c.Distance <= 1 {
			continue
		}
		result0 = append(result0, c)
	}
	return result0
}

// sceneGetIntersections checks for intersections between a ray (cast from the mouse position) and scene objects or a fallback plane.
func (scene *Scene) GetIntersections(camera *rl.Camera3D) Collision {
	ray := rl.GetMouseRay(rl.GetMousePosition(), *camera)
	result := scene.IntersectPoint(ray)

	if !result.Hit {
		// Fallback check against a plane
		fallbackBox := rl.NewBoundingBox(rl.NewVector3(-200.0, 0.4, -200.0), rl.NewVector3(200.0, 0.5, 200.0))
		collision := rl.GetRayCollisionBox(ray, fallbackBox)
		if collision.Hit {
			result = Collision{
				Hit:           false,
				Distance:      collision.Distance,
				Point:         collision.Point,
				Normal:        collision.Normal,
				HitVoxel:      kVoxelNone,
				HitVoxelIndex: "none",
				CollisionHit:  CollisionHitPlane,
			}
		}
	}
	return result
}

// sceneGetIntersections checks for intersections between a ray (cast from the mouse position) and scene objects or a fallback plane.
func (scene *Scene) GetVoxelShadow(voxel *Voxel, light *VxdiDirectionalLight) []Collision {
	ray := rl.NewRay(voxel.Position, light.Direction)
	result := scene.IntersectPoints(ray)

	if len(result) == 0 {
		// Fallback check against a plane
		fallbackBox := rl.NewBoundingBox(rl.NewVector3(-200.0, 0.4, -200.0), rl.NewVector3(200.0, 0.5, 200.0))
		collision := rl.GetRayCollisionBox(ray, fallbackBox)
		result = append(result, Collision{
			Hit:           false,
			Distance:      collision.Distance,
			Point:         collision.Point,
			Normal:        collision.Normal,
			HitVoxel:      kVoxelNone,
			HitVoxelIndex: "none",
			CollisionHit:  CollisionHitPlane,
		})
		return result
	} else {
		result0 := []Collision{}
		for _, k := range result {
			if k.HitVoxel.Equals(voxel) {
				result0 = append(result0, k)
			}
			result0 = append(result0, k)
		}
		return result0
	}
}

func (scene *Scene) GetShadows(light *VxdiDirectionalLight) []Collision {
	shadows := []Collision{}
	// shadowedVoxels := []string{}
	for _, key := range scene.keys {
		voxel := scene.Voxels[key]
		if math.Abs(float64(voxel.Position.Y)) < 1 {
			continue
		}
		if voxel.Material.A < 200 {
			continue
		}
		collisions := scene.GetVoxelShadow(voxel, light)
		for _, col := range collisions {
			if col.Point.Y < -1 {
				continue
			}
			shadows = append(shadows, col)
		}
	}
	return shadows
}

// Assume voxelGetBoundingBox is implemented elsewhere, returning an rl.BoundingBox for a given Voxel.

// Additional methods for Scene

// ExportScene writes the current state of the scene to a text file.
func (s *Scene) ExportScene(filename string) error {
	fmt.Printf("Export Scene bin %s \n", filename)
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("error opening file %s\n", filename)
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	lineNumber := 1
	for _, voxel := range s.Voxels {
		line := fmt.Sprintf("%.0f,%.0f,%.0f,%d,%d,%d,%d\n", voxel.Position.X, voxel.Position.Y, voxel.Position.Z, voxel.Material.R, voxel.Material.G, voxel.Material.B, voxel.Material.A)
		_, err := writer.WriteString(line)
		fmt.Printf("writing %s to file %s at line %d\n", line, filename, lineNumber)
		if err != nil {
			fmt.Printf("error writing file %s\n", filename)
			return err
		}
		lineNumber++
	}
	return writer.Flush()
}

// ImportScene reads a scene from a text file and replaces the current scene's content.
func (s *Scene) ImportScene(filename string) error {
	fmt.Printf("Import Scene bin %s \n", filename)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error opening file %s\n", filename)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	s.Voxels = make(map[string]*Voxel) // Reset current scene
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) < 7 {
			fmt.Printf("error malformed line %d : [%s] %s\n", lineNumber, line, filename)
			lineNumber++
			continue // Skip malformed lines
		}
		x, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {

			fmt.Printf("error parsing x=%s in line %d : %s\n", parts[0], lineNumber, line)
			lineNumber++
			continue
		}
		y, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {

			fmt.Printf("error parsing y=%s in line %d : %s\n", parts[0], lineNumber, line)
			lineNumber++
			continue
		}
		z, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {

			fmt.Printf("error parsing z=%s in line %d : %s\n", parts[2], lineNumber, line)
			lineNumber++
			continue
		}
		r, err := strconv.ParseUint(parts[3], 10, 8)
		if err != nil {

			fmt.Printf("error parsing r=%s in line %d : %s\n", parts[3], lineNumber, line)
			lineNumber++
			continue
		}
		g, err := strconv.ParseUint(parts[4], 10, 8)
		if err != nil {

			fmt.Printf("error parsing g=%s in line %d : %s\n", parts[4], lineNumber, line)
			lineNumber++
			continue
		}
		b, err := strconv.ParseUint(parts[5], 10, 8)
		if err != nil {

			fmt.Printf("error parsing b=%s in line %d : %s\n", parts[5], lineNumber, line)
			lineNumber++
			continue
		}
		a, err := strconv.ParseUint(parts[6], 10, 8)
		if err != nil {

			fmt.Printf("error parsing a=%s in line %d : %s\n", parts[6], lineNumber, line)
			lineNumber++
			continue
		}
		vx := Voxel{
			Position: rl.Vector3{X: float32(x), Y: float32(y), Z: float32(z)},
			Material: rl.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)},
		}
		s.AddVoxel(
			&vx,
		)
		fmt.Printf("found  %v at line %d : %s\n", vx, lineNumber, line)
		lineNumber++
	}
	s.OnChange(s)
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
	for _, v := range s.Voxels {
		err = encoder.Encode(v)
		if err != nil {
			return err
		}
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
	voxels := make(map[string]Voxel)
	err = decoder.Decode(&voxels)
	if err != nil {
		return err
	}

	// s.Voxels = voxels // Replace current scene with the loaded one
	for _, voxel := range voxels {
		s.AddVoxel(&voxel)
		fmt.Printf("loaded %v", voxel)
	}
	fmt.Printf("loaded %v", s.keys)
	s.OnChange(s)
	return nil
}
