package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// item holds length and probability percentage
type item struct {
	length int64
	p      int64
}

// byOrder implements sort.Interface for []item based on custom comparator
type byOrder []item

func (a byOrder) Len() int      { return len(a) }
func (a byOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byOrder) Less(i, j int) bool {
	// compare a[i] and a[j]: return true if a[i] should come before a[j]
	ai, aj := a[i], a[j]
	// use int64 for calculation to avoid overflow
	left := ai.p * ai.length * (100 - aj.p)
	right := aj.p * aj.length * (100 - ai.p)
	return left > right
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]item, n)
	for i := 0; i < n; i++ {
		var length, p int64
		fmt.Fscan(reader, &length, &p)
		arr[i] = item{length: length, p: p}
	}
	sort.Sort(byOrder(arr))

	var S, t float64
	for _, it := range arr {
		// accumulate S and t in float64
		S += 10000.0*float64(it.length) + t*(100.0-float64(it.p))
		t += float64(it.p) * float64(it.length)
	}
	// final result
	result := S / 10000.0
	fmt.Fprintf(writer, "%.10f\n", result)
}
