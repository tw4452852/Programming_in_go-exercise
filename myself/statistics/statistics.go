package main

import (
	"fmt"
	"sort"
	"net/http"
	"strings"
	"strconv"
	"log"
	"math"
)

const (
    pageTop    = `<!DOCTYPE HTML><html><head>
<style>.error{color:#FF0000;}</style></head><title>Statistics</title>
<body><h3>Statistics</h3>
<p>Computes basic statistics for a given list of numbers</p>`
    form       = `<form action="/" method="POST">
<label for="numbers">Numbers (comma or space-separated):</label><br />
<input type="text" name="numbers" size="30"><br />
<input type="submit" value="Calculate">
</form>`
    pageBottom = `</body></html>`
    anError    = `<p class="error">%s</p>`
)

type statistics struct {
	numbers	[]float64
	mean	float64
	median	float64
	mode	[]float64
	devi	float64
}

func getStats(numbers []float64) (stats statistics) {
	stats.numbers = numbers
	sort.Float64s(stats.numbers)
	stats.mean = sum(numbers) / float64(len(numbers))
	stats.median = median(numbers)
	stats.mode = mode(numbers)
	stats.devi = deviation(numbers, stats.mean)
	return stats
}

func deviation(numbers []float64, mean float64) float64 {
	var sum float64
	for _, number := range numbers {
		sum += ((number - mean) * (number - mean))
	}

	return math.Sqrt(sum / float64((len(numbers) - 1)))
}

func mode(numbers []float64) (result []float64) {
	var (
		count = make(map[float64]int, len(numbers))
		max int
	)

	for _, number := range numbers {
		count[number]++
		if count[number] > max {
			max = count[number]
		}
	}

	if max > 1 {
		for number, cnt := range count {
			if cnt == max {
				result = append(result, number)
			}
		}
	}
	sort.Float64s(result)

	return result
}

func sum(numbers []float64) (total float64) {
	for _, x := range numbers {
		total += x
	}
	return total
}

func median(numbers []float64) float64 {
	middle := len(numbers) / 2
	result := numbers[middle]
	if len(numbers) % 2 == 0 {
		result = (result + numbers[middle - 1]) / 2
	}
	return result
}

func homePage(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	fmt.Fprint(writer, pageTop, form)
	if err != nil {
		fmt.Fprintf(writer, anError, err)
	} else {
		if numbers, message, ok := processRequest(request); ok {
			stats := getStats(numbers)
			fmt.Fprint(writer, formatStats(stats))
		} else if message != "" {
			fmt.Fprintf(writer, anError, message)
		}
	}
	fmt.Fprint(writer, pageBottom)
}

func processRequest(request *http.Request) ([]float64, string, bool) {
	var numbers []float64
	if slice, found := request.Form["numbers"]; found && len(slice) > 0 {
		text := strings.Replace(slice[0], ",", " ", -1)
		for _, field := range strings.Fields(text) {
			if x, err := strconv.ParseFloat(field, 64); err != nil {
				return numbers, "'" + field + "' is invalid", false
			} else {
				numbers = append(numbers, x)
			}
		}
	}
	if len(numbers) == 0 {
		return numbers, "", false
	}
	return numbers, "", true
}

func formatStats(stats statistics) string {
    return fmt.Sprintf(`<table border="1">
	<tr><th colspan="2">Results</th></tr>
	<tr><td>Numbers</td><td>%v</td></tr>
	<tr><td>Count</td><td>%d</td></tr>
	<tr><td>Mean</td><td>%f</td></tr>
	<tr><td>Median</td><td>%f</td></tr>
	<tr><td>Mode</td><td>%v</td></tr>
	<tr><td>Std. Dev</td><td>%f</td></tr>
	</table>`, stats.numbers, len(stats.numbers), stats.mean, stats.median, stats.mode, stats.devi)
}

func main() {
	http.HandleFunc("/", homePage)
	if err := http.ListenAndServe(":9001", nil); err != nil {
		log.Fatal("failed to start server", err)
	}
}
