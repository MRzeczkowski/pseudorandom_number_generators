package main

import (
	"fmt"
	"math"
	"path/filepath"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

// Using glibc values
const A uint64 = 1103515245
const C uint64 = 12345
const M uint64 = 1 << 31 // Alternative to `uint64(math.Pow(2, 31))`, this allows it to be `const`

const Seed uint64 = 42

const N = 1_000_000

func main() {
	generateLgc()
	generateCauchy()
	generateCauchyNoTangent()
}

func generateLgc() {

	numbers := make([]float64, N)
	for i := 0; i < N; i++ {
		numbers[i] = float64(lgcGenerator())
	}

	var values plotter.Values
	values = append(values, numbers...)

	mean, stdDev := calculateLgcStats(numbers)
	fmt.Printf("Normalized LGC stats:\n\tMean: %f\n\tStandard deviation: %f\n", mean, stdDev)

	histPlot(values, "Linear Congruential Generator")
}

func lcg(seed uint64) func() uint64 {
	r := seed
	return func() uint64 {
		r = (A*r + C) % M
		return r
	}
}

var lgcGenerator = lcg(Seed)

func generateCauchy() {
	numbers := make([]float64, N)
	var values plotter.Values
	for i := 0; i < N; i++ {
		numbers[i] = math.Tan(getRandomInRange(-0.5, 0.5) * math.Pi)

		// Limiting to values from -4 to 4, otherwise the histogram has too many extremely small and large values and it's impossible to see the distribution.
		if numbers[i] >= -4.0 && numbers[i] <= 4.0 {
			values = append(values, numbers[i])
		}
	}

	histPlot(values, "Cauchy Generator")
}

func generateCauchyNoTangent() {

	const twoOverPi = 2.0 / math.Pi

	fu := func(x float64) float64 {
		if x >= -1.0 && x <= 1.0 {
			return twoOverPi * (1.0 / (1.0 + math.Pow(x, 2.0)))
		}

		return 0.0
	}

	numbers := make([]float64, N)
	var values plotter.Values
	for i := 0; i < N; i++ {
		x := getRandomInRange(-1, 1)
		y := getRandomInRange(0, 1)

		if y/2.0 > fu(x)-(twoOverPi-0.5) {
			x = getRandomInRange(-1, 1)
		}

		// Adding tails, inspired by this article - https://devzine.pl/2011/02/21/generator-liczb-pseudolosowych-cz-3-rozklad-cauchyego/
		if getRandomInRange(0, 1) < 0.5 {
			if x != 0 {
				x = 1 / x
			} else {
				x = math.MaxFloat64
			}
		}

		numbers[i] = x

		// Limiting to values from -4 to 4, otherwise the histogram has too many extremely small and large values and it's impossible to see the distribution.
		if numbers[i] >= -4.0 && numbers[i] <= 4.0 {
			values = append(values, numbers[i])
		}
	}

	histPlot(values, "Cauchy Generator no tangent")
}

func getRandomInRange(min, max float64) float64 {
	x := lgcGenerator()
	xNorm := float64(x) / float64(M)
	return xNorm*(max-min) + min
}

func calculateLgcStats(numbers []float64) (mean, stdDev float64) {
	sum := 0.0
	for _, number := range numbers {
		sum += float64(number)
	}
	mean = sum / float64(len(numbers))

	variance := 0.0
	for _, number := range numbers {
		variance += math.Pow(float64(number)-mean, 2)
	}
	variance /= float64(len(numbers))
	stdDev = math.Sqrt(variance)

	return mean / float64(M), stdDev / float64(M)
}

func histPlot(values plotter.Values, title string) {
	p := plot.New()

	p.Title.Text = title

	hist, err := plotter.NewHist(values, values.Len()/1_000)
	if err != nil {
		panic(err)
	}
	p.Add(hist)

	err = p.Save(1920, 1080, filepath.Join("plots", title+".png"))
	if err != nil {
		panic(err)
	}
}
