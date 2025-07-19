package main

import (
	"bufio"
	"fmt"
	"os"
)

const MX = 300000

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	have := make([][]bool, 10)
	for i := range have {
		have[i] = make([]bool, MX+1)
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] == 1 {
			fmt.Fprintln(writer, 1)
			return
		}
		have[0][a[i]] = true
	}
	// overall gcd
	g := a[0]
	for i := 1; i < n; i++ {
		g = gcd(g, a[i])
	}
	if g > 1 {
		fmt.Fprintln(writer, -1)
		return
	}
	// sieve out multiples
	for i := 2; i <= MX; i++ {
		if !have[0][i] {
			continue
		}
		for j := i + i; j <= MX; j += i {
			have[0][j] = false
		}
	}
	// initial arrays
	ini := make([]int, 0, MX)
	arr := make([]int, 0, MX)
	for i := 1; i <= MX; i++ {
		if have[0][i] {
			ini = append(ini, i)
		}
	}
	arr = append(arr, ini...)

	// iterative dp for subset sizes
	for x := 1; x < 10; x++ {
		nextArr := make([]int, 0, MX)
		seen := have[x]
		prevArr := arr
		for _, v1 := range prevArr {
			for _, v2 := range ini {
				g := gcd(v1, v2)
				if g == 1 {
					fmt.Fprintln(writer, x+1)
					return
				}
				if !seen[g] {
					seen[g] = true
					nextArr = append(nextArr, g)
				}
			}
		}
		arr = nextArr
	}
	// fallback
	fmt.Fprintln(writer, -1)
}
