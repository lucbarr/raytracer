package main

import (
	"image/color"
	"math"
)

// Scene defines a scene
type Scene struct {
	Camera *Camera
	Light  *Light
	Sphere *Sphere
}

// Light defines the light model parameters
type Light struct {
	Source  Vec3
	Ambient Pixel
}

// ApplyModel Calculates lambertial model given a ray and a sphere
func (l *Light) ApplyModel(ray Vec3, s *Sphere) Pixel {
	return l.Ambient
}

// Camera defines the camera
type Camera struct {
	Obs        Vec3
	ImagePlane Quad
}

// Vec3 defines a 3D vector
type Vec3 struct {
	x, y, z float64
}

// Quad defines a quadrilateral
//  0 -------------1
//  |							 |
//  |              |
//  |              |
//  3--------------2
type Quad [4]Vec3

// Sphere defines the sphere
type Sphere struct {
	Center Vec3
	Radius float64
}

// Pixel defines an RGB pixel
type Pixel color.RGBA

func NewPixel(r, g, b uint8) Pixel {
	return Pixel{
		R: r,
		G: g,
		B: b,
		A: 255,
	}
}

// Sub retruns a-b
func Sub(a, b Vec3) Vec3 {
	a.x -= b.x
	a.y -= b.y
	a.z -= b.z
	return a
}

// Add return the sum of all vectors vecs
func Add(vecs ...Vec3) Vec3 {
	ret := Vec3{}
	for _, vec := range vecs {
		ret.x += vec.x
		ret.y += vec.y
		ret.z += vec.z
	}
	return ret
}

// Mul returns a vector multiplied by a scalar
func Mul(a Vec3, s float64) Vec3 {
	a.x *= s
	a.y *= s
	a.z *= s
	return a
}

// Len returns the lenght of a vector
func (v Vec3) Len() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

// Render returns the rendered image RGB Pixels
func (s *Scene) Render(w, h int) ([][]Pixel, error) {
	quad := s.Camera.ImagePlane

	sup := Sub(quad[1], quad[0])
	lat := Sub(quad[3], quad[0])

	v0 := quad[0]

	pixels := make([][]Pixel, h)
	for i := range pixels {
		pixels[i] = make([]Pixel, w)
	}

	var dx, dy Vec3
	for i := 1; i <= w; i++ {
		for j := 1; j <= h; j++ {
			dx = Mul(sup, float64(i)/float64(w))
			dy = Mul(lat, float64(j)/float64(h))

			pixelPoint := Add(v0, dx, dy)
			ray := Sub(pixelPoint, s.Camera.Obs)

			pixel := s.Light.ApplyModel(ray, s.Sphere)
			pixels[i-1][j-1] = pixel
		}
	}

	return pixels, nil
}
