package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

var (
	height  = 400
	width   = 400
	ambient = NewPixel(50, 0, 50)
	k       = NewPixel(139, 0, 139)
	radius  = 100.0
)

func main() {
	quad := [4]Vec3{
		Vec3{-100, 100, 100},
		Vec3{100, 100, 100},
		Vec3{100, -100, 100},
		Vec3{-100, -100, 100},
	}

	cam := &Camera{
		Obs:        Vec3{0, 0, 0},
		ImagePlane: Quad(quad),
	}

	light := &Light{
		Ambient: ambient,
	}

	sphere := &Sphere{
		Center: Vec3{0, 0, 200},
		Radius: radius,
		K:      k,
	}

	scene := &Scene{
		Camera: cam,
		Light:  light,
		Sphere: sphere,
	}

	lights := []Vec3{
		Vec3{1, 1, 0},
		Vec3{0, 1, 0},
		Vec3{-1, 1, 0},
		Vec3{-1, 0, 0},
		Vec3{-1, -1, 0},
		Vec3{0, -1, 0},
		Vec3{1, -1, 0},
		Vec3{1, 0, 0},
	}

	for i, l := range lights {
		light.Source = l

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

		file, err := os.Create(fmt.Sprintf("img%v.png", i))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		err = png.Encode(file, img)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
