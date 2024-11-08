package utils

import (
	"math"
	"tilemap-generator/mapgen/world"
)

const (
	// PERLIN_YWRAPB определяет степень для сдвига, чтобы создать значение для обёртки по оси Y
	// Это значение будет использоваться для того, чтобы сделать перлин-шум цикличным по Y.
	PERLIN_YWRAPB = 4

	// PERLIN_YWRAP устанавливает значение обёртки по оси Y, сдвигая 1 влево на PERLIN_YWRAPB битов
	// Это эквивалентно умножению на 2 в степени 4, то есть 16.
	PERLIN_YWRAP = 1 << PERLIN_YWRAPB

	// PERLIN_ZWRAPB определяет степень для сдвига, чтобы создать значение для обёртки по оси Z
	// Это значение будет использоваться для того, чтобы сделать перлин-шум цикличным по Z.
	PERLIN_ZWRAPB = 8

	// PERLIN_ZWRAP устанавливает значение обёртки по оси Z, сдвигая 1 влево на PERLIN_ZWRAPB битов
	// Это эквивалентно умножению на 2 в степени 8, то есть 256.
	PERLIN_ZWRAP = 1 << PERLIN_ZWRAPB

	// PERLIN_AMP_FALLOFF задает коэффициент затухания амплитуды для каждой последующей октавы
	// Это позволяет создать более мягкие, детализированные шумовые паттерны.
	PERLIN_AMP_FALLOFF = 0.5

	// PERLIN_AVG_POWER устанавливает среднюю мощность шума Перлина
	// Используется для корректировки общего уровня шума и масштабирования значений.
	PERLIN_AVG_POWER = 1.1
)

type PerlinParameters struct {
	Seed   []int
	Config world.Config

	X, Y int
}

type Perlin struct{}

func (p *Perlin) Generate(params PerlinParameters) float64 {
	config := p.normalizeConfig(params.Config)
	frequencyChange := config.FrequencyChange
	borderSmoothness := config.BorderSmoothness
	heightAveraging := config.HeightAveraging
	heightRedistribution := config.HeightRedistribution
	falloff := config.Falloff

	size := len(params.Seed) - 1
	cx := float64(params.X) / float64(params.Config.Width) * frequencyChange
	cy := float64(params.Y) / float64(params.Config.Height) * frequencyChange

	xi := int(math.Floor(cx))
	yi := int(math.Floor(cy))
	xf := cx - float64(xi)
	yf := cy - float64(yi)

	r := 0.0
	ampl := 0.5

	for o := 0; o < int(borderSmoothness); o++ {
		of := xi + (yi << PERLIN_YWRAPB)

		rxf := p.scaledCosine(xf)
		ryf := p.scaledCosine(yf)

		n1 := float64(params.Seed[of&size])
		n1 += rxf * (float64(params.Seed[(of+1)&size]) - n1)

		n2 := float64(params.Seed[(of+PERLIN_YWRAP)&size])
		n2 += rxf * (float64(params.Seed[(of+PERLIN_YWRAP+1)&size]) - n2)

		n1 += ryf * (n2 - n1)
		r += n1 * ampl
		ampl *= PERLIN_AMP_FALLOFF
		of += PERLIN_ZWRAP

		xi <<= 1
		xf *= 2
		if xf >= 1.0 {
			xi++
			xf--
		}

		yi <<= 1
		yf *= 2
		if yf >= 1.0 {
			yi++
			yf--
		}
	}

	if heightAveraging {
		if r > 0.5 {
			r = math.Pow(r, (1.5-r)/PERLIN_AVG_POWER)
		} else if r < 0.5 {
			r = math.Pow(r, (1.5-r)*PERLIN_AVG_POWER)
		}
	}

	r = math.Pow(r, heightRedistribution)

	if falloff > 0.0 {
		r *= p.heightFalloff(float64(params.X), params.Config.Width, falloff) *
			p.heightFalloff(float64(params.Y), params.Config.Height, falloff)
	}

	return r
}

func (p *Perlin) clamp(value float64, defaultValue float64, limit [2]float64) float64 {
	v := value
	if math.IsNaN(value) {
		v = defaultValue
	}

	return math.Max(limit[0], math.Min(limit[1], v))
}

func (p *Perlin) scaledCosine(i float64) float64 {
	return 0.5 * (1.0 - math.Cos(i*math.Pi))
}

func (p *Perlin) smootherStep(x float64) float64 {
	return 3*math.Pow(x, 2) - 2*math.Pow(x, 3)
}

func (p *Perlin) heightFalloff(
	offset float64, length int, falloff float64,
) float64 {
	radius := float64(length / 2)
	distance := math.Abs(radius - offset)
	target := radius * (1 - falloff)

	if distance < target {
		return 1.0
	}

	x := (distance - target) / radius / (1 - target/radius)
	x = math.Min(1, math.Max(0, x))

	return 1 - p.smootherStep(x)
}

func (p *Perlin) normalizeConfig(config world.Config) *world.Config {
	frequencyChange := math.Round(p.clamp(config.FrequencyChange, 0.3, [2]float64{0.0, 1.0})*31 + 1)
	borderSmoothness := math.Round((1-p.clamp(config.BorderSmoothness, 0.5, [2]float64{0.0, 1.0}))*14 + 1)
	heightRedistribution := 2.0 - p.clamp(config.HeightRedistribution, 1.0, [2]float64{0.5, 1.5})
	falloff := p.clamp(config.Falloff, 0.0, [2]float64{0.0, 0.9})

	heightAveraging := true
	if !config.HeightAveraging {
		heightAveraging = config.HeightAveraging
	}

	return &world.Config{
		FrequencyChange:      frequencyChange,
		BorderSmoothness:     borderSmoothness,
		HeightRedistribution: heightRedistribution,
		Falloff:              falloff,
		HeightAveraging:      heightAveraging,
	}
}
