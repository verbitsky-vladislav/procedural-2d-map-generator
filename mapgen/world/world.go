package world

type Config struct {
	Width, Height int

	/**
	 * Frequency of biomes change
	 * Default: 0.3
	 * Min: 0.0, Max: 1.0
	 */
	FrequencyChange float64

	/**
	 * Smoothness of biomes borders
	 * Default: 0.5
	 * Min: 0.0, Max: 1.0
	 */
	BorderSmoothness float64

	/**
	 * Redistribution of biomes height
	 * Default: 1.0
	 * Min: 0.5, Max: 1.5
	 */
	HeightRedistribution float64

	/**
	 * Averaging of biomes height
	 * Default: true
	 */
	HeightAveraging bool

	/**
	 * Scale of falloff area
	 * Default: 0.0
	 */
	Falloff float64
}

type Point struct {
	X, Y int
}
