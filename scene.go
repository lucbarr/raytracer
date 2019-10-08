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
// constraints:
// ray(t) = o + vt
// F(x) = len2(x-c) + len2(r)
func (l *Light) ApplyModel(v Vec3, o Vec3, s *Sphere) Pixel {
	c := s.Center

	oc := Sub(o, c)
	oc2 := oc.Len2()
	r2 := s.Radius * s.Radius
	v2 := v.Len2()
	dot := Dot(v, oc)

	delta := 4*dot*dot - 4*(v2*(oc2-r2))

	if delta < 0 {
		return Pixel{
			R: 0,
			G: 0,
			B: 0,
			A: 255,
		}
	}

	t1 := ((-2 * dot) + math.Sqrt(delta)) / 2 * v2
	t2 := ((-2 * dot) - math.Sqrt(delta)) / 2 * v2

	vt1 := Mul(v, t1)
	vt2 := Mul(v, t2)

	x1 := Add(o, vt1)
	x2 := Add(o, vt2)

	dx1 := Sub(x1, o)
	dx2 := Sub(x2, o)

	var closest Vec3
	if dx1.Len2() < dx2.Len() {
		closest = dx1
	} else {
		closest = dx2
	}

	n := Sub(closest, c)

	return l.applyLight(n)
}

func (l *Light) applyLight(n Vec3) Pixel {
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

// NewPixel returns a new pixel
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

// Dot returns the dot product between a and b
func Dot(a, b Vec3) float64 {
	return (a.x * b.x) + (a.y * b.y) + (a.z * b.z)
}

// Len2 returns the lenght of a vector
func (v Vec3) Len2() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z
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

			pixel := s.Light.ApplyModel(ray, s.Camera.Obs, s.Sphere)
			pixels[i-1][j-1] = pixel
		}
	}

	return pixels, nil
}
