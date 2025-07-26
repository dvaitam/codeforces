package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, K, T, R int
	if _, err := fmt.Fscan(in, &N); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &K); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &R); err != nil {
		return
	}

	// Skip initial SINR values (R*K*T lines with N floats)
	for t := 0; t < T; t++ {
		for k := 0; k < K; k++ {
			for r := 0; r < R; r++ {
				for n := 0; n < N; n++ {
					var dummy float64
					fmt.Fscan(in, &dummy)
				}
			}
		}
	}

	// Skip interference factors (N*R*K lines with N floats)
	for m := 0; m < N; m++ {
		for r := 0; r < R; r++ {
			for k := 0; k < K; k++ {
				for n := 0; n < N; n++ {
					var dummy float64
					fmt.Fscan(in, &dummy)
				}
			}
		}
	}

	var J int
	fmt.Fscan(in, &J)
	// Skip frame information
	for i := 0; i < J; i++ {
		var a, b, c, d, e int
		fmt.Fscan(in, &a, &b, &c, &d, &e)
	}

	// Output zero power allocation (valid but trivial)
	zeros := make([]string, N)
	for i := 0; i < N; i++ {
		zeros[i] = "0"
	}
	line := strings.Join(zeros, " ")
	for t := 0; t < T; t++ {
		for k := 0; k < K; k++ {
			for r := 0; r < R; r++ {
				fmt.Fprintln(out, line)
			}
		}
	}
}
