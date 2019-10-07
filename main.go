package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

const (
	height = 400
	width  = 400
)

func main() {
	quad := [4]Vec3{
		Vec3{100, 200, 200},
		Vec3{100, 200, 200},
		Vec3{100, 200, 200},
		Vec3{100, 200, 200},
	}

	cam := &Camera{
		Obs:        Vec3{0, 0, 0},
		ImagePlane: Quad(quad),
	}

	light := &Light{
		Source:  Vec3{1000, 0, 0},
		Ambient: NewPixel(100, 0, 0),
	}

	sphere := &Sphere{
		Center: Vec3{500, 0, 0},
		Radius: 100,
	}

	scene := &Scene{
		Camera: cam,
		Light:  light,
		Sphere: sphere,
	}

	pixels, err := scene.Render(width, height)
	if err != nil {
		fmt.Println(err)
		return
	}

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for i, row := range pixels {
		for j, pixel := range row {
			img.Set(i, j, color.RGBA{
				R: pixel.R,
				G: pixel.G,
				B: pixel.B,
				A: pixel.A,
			})
		}
	}

	file, _ := os.Create("img.png")
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		fmt.Println(err)
		return
	}
}
