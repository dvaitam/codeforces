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
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	M := make([][]int64, n)
	for i := range M {
		M[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(reader, &M[i][j])
		}
	}

	// compute the first element using indices 0,1,2
	a := make([]int64, n)
	if n >= 3 {
		val := M[0][1] * M[0][2] / M[1][2]
		a0 := int64(math.Sqrt(float64(val)))
		for (a0+1)*(a0+1) <= val {
			a0++
		}
		for a0*a0 > val {
			a0--
		}
		a[0] = a0
	} else if n == 2 {
		// not expected by constraints but handle gracefully
		a[0] = int64(math.Sqrt(float64(M[0][1])))
	} else if n == 1 {
		a[0] = 0
	}

	for i := 1; i < n; i++ {
		if a[0] != 0 {
			a[i] = M[0][i] / a[0]
		} else {
			a[i] = 0
		}
	}

	for i, v := range a {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, v)
	}
	writer.WriteByte('\n')
}
