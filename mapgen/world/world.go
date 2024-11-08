package world

import (
	"fmt"
	"tilemap-generator/mapgen/biome"
)

type Point struct {
	X, Y int64
}

type Config struct {
	Width, Height int64

	// Частота изменения биомов по карте.
	// Определяет, насколько часто будут встречаться переходы между биомами.
	// Значение от 0.0 до 1.0, где 0.0 - изменения почти не происходят, а 1.0 - переходы между биомами максимально частые.
	// Default: 0.3 (умеренная частота изменения).
	FrequencyChange float64

	// Плавность переходов между биомами, т.е. насколько резкими или мягкими будут границы биомов.
	// 0.0 означает резкие, чёткие границы, а 1.0 — плавные, едва заметные переходы.
	// Default: 0.5 (средняя плавность).
	BorderSmoothness float64

	// Перераспределение высот биомов.
	// Значение определяет, насколько высоты биомов могут быть изменены или распределены по карте.
	// Чем выше значение, тем сильнее могут варьироваться высоты биомов в пределах заданного диапазона.
	// Default: 1.0 (среднее перераспределение высот).
	// Min: 0.5 (меньшее изменение высот), Max: 1.5 (большее изменение высот).
	HeightRedistribution float64

	// Параметр Falloff влияет на область смягчения или исчезновения биомов.
	// Он определяет, как сильно границы биомов будут "размазываться" на определённом расстоянии от границ.
	// Значение 0.0 указывает на отсутствие области перехода, а большие значения создают более широкие области для сглаживания.
	// Default: 0.0 (плавное падение).
	Falloff float64

	// Если включена эта опция, высоты биомов будут усредняться, чтобы уменьшить резкие перепады между биомами.
	// Это может помочь создать более естественные ландшафты с мягкими переходами.
	// Default: true (включено).
	HeightAveraging bool
}

type World struct {
	Width, Height int64
	Seed          int
	Matrix        [][]biome.Data // todo rewrite from T to something else
}

func NewWorld(matrix [][]biome.Data, seed int) *World {
	// Проверка на пустую матрицу и пустые строки
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		// В случае пустой матрицы или пустых строк, можно вернуть мир с размерами 0x0
		fmt.Println("Error: Matrix is empty or contains empty rows.")
		return &World{
			Width:  0,
			Height: 0,
			Seed:   seed,
			Matrix: matrix,
		}
	}

	return &World{
		Width:  int64(len(matrix[0])), // Длина первой строки (ширина карты)
		Height: int64(len(matrix)),    // Количество строк (высота карты)
		Seed:   seed,
		Matrix: matrix,
	}
}

func (w *World) Each(callback func(point Point, biome biome.Data) bool) {
	for y := int64(0); y < w.Height; y++ {
		for x := int64(0); x < w.Width; x++ {
			res := callback(Point{x, y}, w.Matrix[y][x])
			if !res {
				return
			}
		}
	}
}

func (w *World) GetAt(point Point) biome.Data {
	return w.Matrix[point.Y][point.X]
}

func (w *World) ReplaceAt(point Point, data biome.Data) {
	w.Matrix[point.Y][point.X] = data
}
