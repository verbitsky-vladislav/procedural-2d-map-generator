package utils

import (
	"github.com/aquilax/go-perlin"
	"math"
	"tilemap-generator/mapgen/world"
)

const (
	PERLIN_AMP_FALLOFF = 0.5
	PERLIN_AVG_POWER   = 1.1
)

type PerlinParameters struct {
	Seed   []int64
	Config world.Config
	X, Y   int64
}

type Perlin struct{}

func (p *Perlin) Generate(params PerlinParameters) float64 {
	config := p.normalizeConfig(params.Config)

	// Генерация случайного сидера для библиотеки go-perlin
	newPerlin := perlin.NewPerlin(2, 2, 3, 256)

	// Преобразуем координаты и параметры
	x := float64(params.X) / float64(config.Width)  //* config.FrequencyChange
	y := float64(params.Y) / float64(config.Height) //* config.FrequencyChange

	// Генерируем значение шума в точке (x, y)
	r := newPerlin.Noise2D(x, y)

	// Масштабируем значение шума, чтобы оно попадало в диапазон биомов
	r = p.scaleHeight(r)

	// Применяем перераспределение высоты
	r = p.applyHeightRedistribution(r, *config)

	// Возвращаем результат
	return r
}

func (p *Perlin) scaleHeight(r float64) float64 {
	// Масштабируем значение шума так, чтобы оно попадало в диапазон от 0 до 1
	return (r + 1) / 2 // Преобразуем из диапазона (-1, 1) в (0, 1)
}

func (p *Perlin) applyHeightRedistribution(r float64, config world.Config) float64 {
	if config.HeightAveraging {
		if r > 0.5 {
			r = math.Pow(r, (1.5-r)/PERLIN_AVG_POWER)
		} else if r < 0.5 {
			r = math.Pow(r, (1.5-r)*PERLIN_AVG_POWER)
		}
	}
	return math.Pow(r, config.HeightRedistribution)
}

func (p *Perlin) normalizeConfig(config world.Config) *world.Config {
	config.FrequencyChange = math.Round(p.clamp(&config.FrequencyChange, 0.5)*31 + 1)
	config.BorderSmoothness = math.Round((1-p.clamp(&config.BorderSmoothness, 0.3))*14 + 1)
	config.HeightRedistribution = p.clamp(&config.HeightRedistribution, 1.0, []float64{0.5, 1.5}...)
	config.Falloff = p.clamp(&config.Falloff, 0.0, []float64{0.0, 0.9}...)

	if &config.HeightAveraging == nil {
		config.HeightAveraging = true
	}
	return &config
}

func (p *Perlin) clamp(value *float64, defaultValue float64, limit ...float64) float64 {
	if len(limit) == 0 {
		limit = []float64{0, 1}
	}
	if value == nil {
		value = &defaultValue
	}
	return math.Max(limit[0], math.Min(limit[1], *value))
}
