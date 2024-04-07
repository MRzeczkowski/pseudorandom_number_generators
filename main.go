package main

import (
	"fmt"
	"math"
	"path/filepath"
	"sort"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

// Using glibc values
const A uint64 = 1103515245
const C uint64 = 12345
const M uint64 = 1 << 31 // Alternative to `uint64(math.Pow(2, 31))`, this allows it to be `const`

const Seed uint64 = 42

const X0 = 0.0
const Gamma = 1.0

const N = 10_000_000

func main() {
	fmt.Printf("Generating %d numbers\n\n", N)

	generateLgc()
	generateCauchy()
	generateCauchyNoTangent()
}

func generateLgc() {

	start := time.Now()
	numbers := make([]float64, N)
	for i := 0; i < N; i++ {
		numbers[i] = float64(lgcGenerator())
	}
	elapsed := time.Since(start)
	fmt.Printf("Generating numbers using LGC took %s\n", elapsed)

	var values plotter.Values
	values = append(values, numbers...)

	mean, stdDev := calculateLgcStats(numbers)
	fmt.Printf("LGC stats (normalized):\n\tMean: %f\n\tStandard deviation: %f\n\n", mean, stdDev)

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

	start := time.Now()
	numbers := make([]float64, N)
	for i := 0; i < N; i++ {
		numbers[i] = X0 + Gamma*math.Tan(getRandomInRange(-0.5, 0.5)*math.Pi)
	}
	elapsed := time.Since(start)
	fmt.Printf("Generating Cauchy numbers took %s\n", elapsed)

	q1, median, q3, IQR := calculateCauchyStats(numbers)
	fmt.Printf("Cauchy stats:\n\t1st quartile: %f\n\tMedian: %f\n\t3rd quartile: %f\n\tInterquartile range: %f\n\n", q1, median, q3, IQR)

	values := getCauchyValuesForHistogram(numbers)

	histPlot(values, "Cauchy Generator")
}

const TwoOverPi = 2.0 / math.Pi

func generateCauchyNoTangent() {

	start := time.Now()
	numbers := make([]float64, N)
	const twoOverPiMinusHalf = TwoOverPi - 0.5
	for i := 0; i < N; i++ {
		x := getRandomInRange(-1, 1)
		y := getRandomInRange(0, 1)

		if y > 2.0*(fu(x)-twoOverPiMinusHalf) {
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

		numbers[i] = X0 + Gamma*x
	}
	elapsed := time.Since(start)
	fmt.Printf("Generating Cauchy numbers without using tangent took %s\n", elapsed)

	q1, median, q3, IQR := calculateCauchyStats(numbers)
	fmt.Printf("Cauchy no tangent stats:\n\t1st quartile: %f\n\tMedian: %f\n\t3rd quartile: %f\n\tInterquartile range: %f\n\n", q1, median, q3, IQR)

	values := getCauchyValuesForHistogram(numbers)

	histPlot(values, "Cauchy Generator no tangent")
}

func fu(x float64) float64 {
	if x >= -1.0 && x <= 1.0 {
		return TwoOverPi / (1.0 + x*x)
	}

	return 0.0
}

func getCauchyValuesForHistogram(numbers []float64) plotter.Values {

	var values plotter.Values

	for i := 0; i < N; i++ {
		// Limiting to values from -4 to 4, otherwise the histogram has too many extremely small and large values and it's impossible to see the distribution.
		if numbers[i] >= -4.0 && numbers[i] <= 4.0 {
			values = append(values, numbers[i])
		}
	}

	return values
}

func calculateCauchyStats(numbers []float64) (q1, median, q3, IQR float64) {
	sort.Float64s(numbers)

	median = numbers[N/2]
	if N%2 == 0 {
		median = (median + numbers[N/2-1]) / 2
	}

	q1 = numbers[N/4]
	q3 = numbers[(3*N)/4]
	IQR = q3 - q1

	return q1, median, q3, IQR
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
		variance += math.Pow(number-mean, 2)
	}
	variance /= float64(len(numbers))
	stdDev = math.Sqrt(variance)

	return mean / float64(M), stdDev / float64(M)
}

func histPlot(values plotter.Values, title string) {
	p := plot.New()

	p.Title.Text = title

	// Taking 1% of values as bins
	bins := int(float32(values.Len()) * 0.01)

	hist, err := plotter.NewHist(values, bins)
	if err != nil {
		panic(err)
	}
	p.Add(hist)

	err = p.Save(1440, 720, filepath.Join("plots", title+".png"))
	if err != nil {
		panic(err)
	}
}
