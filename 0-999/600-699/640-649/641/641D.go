package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n int
	fmt.Fscan(reader, &n)
	px := make([]float64, n)
	pn := make([]float64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &px[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &pn[i])
	}
	a := make([]float64, n)
	b := make([]float64, n)
	var fPrev, gPrev float64
	for i := 0; i < n; i++ {
		x := px[i]
		p := pn[i]
		sum := x + p - fPrev + gPrev
		t := sum*sum - 4*(x-fPrev*(x+p))
		if t < 0 {
			t = 0
		}
		ai := (sum - math.Sqrt(t)) / 2
		bi := x + p - ai
		a[i] = ai
		b[i] = bi
		fPrev += ai
		gPrev += bi
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			writer.WriteRune(' ')
		}
		fmt.Fprintf(writer, "%.7f", a[i])
	}
	writer.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			writer.WriteRune(' ')
		}
		fmt.Fprintf(writer, "%.7f", b[i])
	}
	writer.WriteByte('\n')
}
