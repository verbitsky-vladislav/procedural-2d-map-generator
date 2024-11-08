package main

import (
	"fmt"
	"log"
	"math/rand"
	"tilemap-generator/image"
	"tilemap-generator/mapgen/biome"
	"tilemap-generator/mapgen/generator"
	"tilemap-generator/mapgen/world"
	"time"
)

type WorldBiomeConfig struct {
	LowerBound float64
	UpperBound float64
}

type Biome struct {
	Params WorldBiomeConfig
	Data   biome.Data
}

var BIOMES = []Biome{
	{
		Params: WorldBiomeConfig{UpperBound: 0.08, LowerBound: 0.0},
		Data: biome.Data{
			Name:   "Liquid",
			NameRU: "Океан", // Русское название для Liquid
			Color:  "#4292c4",
		},
	},
	{
		Params: WorldBiomeConfig{LowerBound: 0.08, UpperBound: 0.11},
		Data: biome.Data{
			Name:   "Liquid",
			NameRU: "Морская вода", // Русское название для Liquid
			Color:  "#4c9ccd",
		},
	},
	{
		Params: WorldBiomeConfig{LowerBound: 0.11, UpperBound: 0.14},
		Data: biome.Data{
			Name:   "Liquid",
			NameRU: "Мелководье", // Русское название для Liquid
			Color:  "#51a5d8",
		},
	},
	{
		Params: WorldBiomeConfig{LowerBound: 0.14, UpperBound: 0.17},
		Data: biome.Data{
			Name:   "Liquid",
			NameRU: "Лагуна", // Русское название для Liquid
			Color:  "#56aade",
		},
	},
	{
		Params: WorldBiomeConfig{LowerBound: 0.17, UpperBound: 0.22},
		Data: biome.Data{
			Name:   "Coast",
			NameRU: "Побережье", // Русское название для Coast
			Color:  "#c5ac6d",
		},
	},
	{
		Params: WorldBiomeConfig{LowerBound: 0.22, UpperBound: 0.25},
		Data: biome.Data{
			Name:   "Coast",
			NameRU: "Песчаные пляжи", // Русское название для Coast
			Color:  "#ccb475",
		},
	},
	{
		Params: WorldBiomeConfig{LowerBound: 0.25, UpperBound: 0.28},
		Data: biome.Data{
			Name:   "Coast",
			NameRU: "Коралловые рифы", // Русское название для Coast
			Color:  "#d2ba7d",
		},
	},
	{
		Params: WorldBiomeConfig{LowerBound: 0.28, UpperBound: 0.34},
		Data: biome.Data{
			Name:   "Fields",
			NameRU: "Зеленые поля", // Русское название для Fields
			Color:  "#67c72b",
		},
	},
	{
		Params: WorldBiomeConfig{LowerBound: 0.34, UpperBound: 0.46},
		Data: biome.Data{
			Name:   "Fields",
			NameRU: "Луга", // Русское название для Fields
			Color:  "#5dbc21",
		},
	},
	{
		Params: WorldBiomeConfig{LowerBound: 0.46, UpperBound: 0.65},
		Data: biome.Data{
			Name:   "Fields",
			NameRU: "Широкие поля", // Русское название для Fields
			Color:  "#56ae1e",
		},
	},
	{
		Params: WorldBiomeConfig{LowerBound: 0.65, UpperBound: 0.72},
		Data: biome.Data{
			Name:   "Mounts",
			NameRU: "Горы", // Русское название для Mounts
			Color:  "#333333",
		},
	},
	{
		Params: WorldBiomeConfig{LowerBound: 0.72, UpperBound: 0.79},
		Data: biome.Data{
			Name:   "Mounts",
			NameRU: "Высокие горы", // Русское название для Mounts
			Color:  "#444444",
		},
	},
	{
		Params: WorldBiomeConfig{LowerBound: 0.79, UpperBound: 1.00},
		Data: biome.Data{
			Name:   "Mounts",
			NameRU: "Заснеженные вершины", // Русское название для Mounts
			Color:  "#555555",
		},
	},
}

func generateSeededInt64Slice(size int, seed int64) []int64 {
	// Устанавливаем сид с помощью времени или произвольного значения
	rand.Seed(seed)

	// Генерируем срез случайных чисел типа int64
	result := make([]int64, size)
	for i := 0; i < size; i++ {
		result[i] = rand.Int63() // Генерирует случайное число int64
	}

	return result
}

func main() {
	cfg := world.Config{
		Width:                1000,
		Height:               1000,
		FrequencyChange:      0.30,
		BorderSmoothness:     0.50,
		HeightRedistribution: 1.00,
		Falloff:              0.00,
		HeightAveraging:      true,
	}

	seed := time.Now().Unix()

	g := generator.NewGenerator(cfg, make([]biome.WorldBiome, 0))

	for _, b := range BIOMES {
		log.Printf("Biom: upperbound: %v, lowerbound: %v", b.Params.UpperBound, b.Params.LowerBound)
		g.AddBiome(b.Params.LowerBound, b.Params.UpperBound, b.Data)
	}

	w := g.Generate(generator.WorldGeneratorParams{
		Seed:     generateSeededInt64Slice(512, seed),
		SeedSize: 512,
		OffsetX:  0,
		OffsetY:  0,
	})

	img := image.CreateImageFromWorld(w)
	//Сохраняем изображение в файл
	err := image.SaveImage(img, "biome_map.png")
	if err != nil {
		fmt.Println("Error saving image:", err)
	} else {
		fmt.Println("Image saved successfully!")
	}

}
