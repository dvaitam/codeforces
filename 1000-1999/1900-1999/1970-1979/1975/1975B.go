package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		solveB(in, out)
	}
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func primeFactors(x int) []int {
	factors := make([]int, 0)
	for i := 2; i*i <= x; i++ {
		if x%i == 0 {
			factors = append(factors, i)
			for x%i == 0 {
				x /= i
			}
		}
	}
	if x > 1 {
		factors = append(factors, x)
	}
	return factors
}

func solveB(in *bufio.Reader, out *bufio.Writer) {
	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	pf := primeFactors(a[0])
	for _, p := range pf {
		g := 0
		for _, v := range a {
			if v%p != 0 {
				g = gcd(g, v)
			}
		}
		if g == 0 {
			fmt.Fprintln(out, "Yes")
			return
		}
		if g > 1 {
			fmt.Fprintln(out, "Yes")
			return
		}
	}
	fmt.Fprintln(out, "No")
}
