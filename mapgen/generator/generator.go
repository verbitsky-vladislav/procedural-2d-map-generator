package generator

import (
	"log"
	"tilemap-generator/mapgen/biome"
	"tilemap-generator/mapgen/utils"
	"tilemap-generator/mapgen/world"
)

type WorldGeneratorParams struct {
	Seed     []int64
	SeedSize int64
	OffsetX  int64
	OffsetY  int64
}

type WorldGenerator struct {
	Config world.Config
	Biomes []biome.WorldBiome
}

func NewGenerator(config world.Config, biomes []biome.WorldBiome) *WorldGenerator {
	return &WorldGenerator{
		Config: config,
		Biomes: biomes,
	}
}

func (wg *WorldGenerator) AddBiome(lowerBound, upperBound float64, data biome.Data) biome.WorldBiome {
	b := biome.NewWorldBiome(lowerBound, upperBound, data)
	wg.Biomes = append(wg.Biomes, *b)

	return *b
}

func (wg *WorldGenerator) ClearBiomes() {
	wg.Biomes = make([]biome.WorldBiome, 0)
}

func (wg *WorldGenerator) GetBiomes() []biome.WorldBiome {
	return wg.Biomes
}

func (wg *WorldGenerator) PeakBiome(height float64) *biome.WorldBiome {
	for _, b := range wg.Biomes {
		// Проверка, попадает ли высота в диапазон биома
		if height >= b.LowerBound && height < b.UpperBound {
			// Логируем, какой биом был выбран
			log.Printf("Height %.5f is in range [%.5f, %.5f), selecting biome: %s\n", height, b.LowerBound, b.UpperBound, b.Data.NameRU)
			return &b // Возвращаем ссылку на биом
		}
	}

	return nil // Если биом не найден
}

func (wg *WorldGenerator) Generate(params WorldGeneratorParams) *world.World {
	// Проверяем на nil или пустой срез для seed
	currentSeed := params.Seed
	if len(currentSeed) == 0 {
		currentSeed = utils.Seed(params.SeedSize)
	} else {
	}

	// Создаем матрицу
	matrix := make([][]biome.Data, wg.Config.Height) // Создаем внешний срез (по числу строк)

	perlin := utils.Perlin{}

	// Генерация карты
	for y := int64(0); y < wg.Config.Height; y++ {
		matrix[y] = make([]biome.Data, wg.Config.Width) // Инициализируем строку
		for x := int64(0); x < wg.Config.Width; x++ {
			// Генерация высоты для каждой клетки
			height := perlin.Generate(utils.PerlinParameters{
				Seed:   currentSeed,
				X:      x + params.OffsetX,
				Y:      y + params.OffsetY,
				Config: wg.Config,
			})

			// Определяем биом для данной высоты
			b := wg.PeakBiome(height)
			if b != nil {
				matrix[y][x] = b.Data
			}
		}
	}

	// Возвращаем новый мир
	return world.NewWorld(matrix, currentSeed)
}
