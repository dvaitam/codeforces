package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// perfectSquare checks if x is a perfect square
func perfectSquare(x int64) bool {
	s := int64(math.Sqrt(float64(x) + 0.5))
	return s*s == x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var N int
	fmt.Fscan(reader, &N)
	N64 := int64(N)
	included := make([]bool, N+1)
	for i := 1; i <= N; i++ {
		included[i] = true
	}

	standard := func() {
		for n := N; n%4 != 0; n-- {
			included[n] = false
		}
		idx := (N - N%4) / 2
		included[idx] = false
	}

	switch N % 4 {
	case 3:
		if perfectSquare(N64 + 1) {
			included[N/2+1] = false
			included[N] = false
		} else if perfectSquare(2 * (N64 / 2) * (N64/2 - 1)) {
			included[N/2-2] = false
			included[N] = false
		} else if perfectSquare((N64/2 - 1) * N64) {
			included[N/2-2] = false
			included[N-2] = false
		} else {
			included[2] = false
			included[N/2] = false
			included[N] = false
		}
	case 2:
		if perfectSquare(2 * (N64/2 + 1)) {
			included[N/2+1] = false
		} else if perfectSquare(2 * (N64 / 2) * (N64/2 - 1)) {
			included[N/2-2] = false
		} else {
			included[2] = false
			included[N/2] = false
		}
	default:
		standard()
	}

	// ensure 1 is always included
	if N >= 1 {
		included[1] = true
	}

	var remaining []int
	for i := 1; i <= N; i++ {
		if included[i] {
			remaining = append(remaining, i)
		}
	}

	fmt.Fprintln(writer, len(remaining))
	for i, v := range remaining {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, v)
	}
	writer.WriteByte('\n')
}
