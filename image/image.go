package image

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"tilemap-generator/mapgen/world"
)

func CreateImageFromWorld(world *world.World) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(world.Width), int(world.Height)))

	for y := int64(0); y < world.Height; y++ {
		for x := int64(0); x < world.Width; x++ {
			// Получаем цвет биома
			hexColor := world.Matrix[y][x].Color
			c, err := parseHexColor(hexColor)
			if err != nil {
				fmt.Println("Error parsing hex color:", err)
				continue
			}
			// Устанавливаем пиксель
			img.Set(int(x), int(y), c)
		}
	}

	return img
}

func parseHexColor(hex string) (color.Color, error) {
	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return nil, err
	}
	return color.RGBA{R: r, G: g, B: b, A: 255}, nil
}

func SaveImage(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}
