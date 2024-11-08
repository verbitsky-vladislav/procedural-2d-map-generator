package biome

import "math"

type WorldBiome struct {
	LowerBound, UpperBound float64
	Data                   Data
}

type Data struct {
	Name   string
	NameRU string
	Color  string
}

func NewWorldBiome(lowerBound, upperBound float64, data Data) *WorldBiome {
	return &WorldBiome{
		LowerBound: math.Max(0, lowerBound),
		UpperBound: math.Min(1, upperBound),
		Data:       data,
	}
}
