package generator

import (
	"fmt"
	"testing"
	"tilemap-generator/mapgen/world"
	"time"
)

// Параметры для бенчмарка
var wgConfig = world.Config{
	Width:                1000,
	Height:               1000,
	FrequencyChange:      0.30,
	BorderSmoothness:     0.50,
	HeightRedistribution: 1.00,
	Falloff:              0.00,
	HeightAveraging:      true,
}
var worldGen = &WorldGenerator{Config: wgConfig}

var currentSeed = int(time.Now().Unix())

// BenchmarkGenerateWorld бенчмарк для генерации мира с разными размерами матрицы
func BenchmarkGenerateWorld(b *testing.B) {
	for _, size := range []struct {
		width, height int64
	}{
		{100, 100},   // Маленький размер карты
		{500, 500},   // Средний размер карты
		{1000, 1000}, // Большой размер карты
	} {
		worldGen.Config.Width = size.width
		worldGen.Config.Height = size.height

		// Настраиваем параметры генерации
		params := WorldGeneratorParams{
			Seed:    currentSeed,
			OffsetX: 0,
			OffsetY: 0,
		}
		b.Run(fmt.Sprintf("WorldSize_%dx%d", size.width, size.height), func(b *testing.B) {

			// Выполняем бенчмарк
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = worldGen.Generate(params)
			}
		})
	}
}
