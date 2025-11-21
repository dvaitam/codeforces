package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func buildPerm(n, step, shift int) []int {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = (i*step + shift) % n
		p[i]++
	}
	return p
}

func factorialUpTo8(n int) int {
	f := 1
	for i := 2; i <= n; i++ {
		f *= i
	}
	return f
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)

		// Quick impossibility checks.
		if n == 1 {
			if k == 1 {
				fmt.Fprintln(out, "YES")
				fmt.Fprintln(out, 1)
			} else {
				fmt.Fprintln(out, "NO")
			}
			continue
		}

		if n == 2 {
			if k == 2 {
				fmt.Fprintln(out, "YES")
				fmt.Fprintln(out, "1 2")
				fmt.Fprintln(out, "2 1")
			} else {
				fmt.Fprintln(out, "NO")
			}
			continue
		}

		if k == 1 || (n%2 == 0 && k%2 == 1) || (n%2 == 1 && k%2 == 1 && k < n) {
			fmt.Fprintln(out, "NO")
			continue
		}

		if n <= 8 {
			lim := factorialUpTo8(n)
			if k > lim {
				fmt.Fprintln(out, "NO")
				continue
			}
		}

		perms := make([][]int, 0, k)
		remaining := k
		blockUsed := false

		// For odd k (and odd n), take one block of all cyclic shifts; the rest will be done with pairs.
		if remaining%2 == 1 {
			blockUsed = true
			for shift := 0; shift < n; shift++ {
				perms = append(perms, buildPerm(n, 1, shift))
			}
			remaining -= n
		}

		pairs := remaining / 2
		step := 1
		if blockUsed {
			step = 2 // avoid duplicating the step used in the full block
		}

		for pairs > 0 && step <= n {
			if gcd(step, n) != 1 {
				step++
				continue
			}
			// Skip step=1 if already used in the block.
			if blockUsed && step == 1 {
				step++
				continue
			}
			for shift := 0; shift < n && pairs > 0; shift++ {
				base := buildPerm(n, step, shift)
				comp := make([]int, n)
				for i, v := range base {
					comp[i] = n + 1 - v
				}
				perms = append(perms, base, comp)
				pairs--
			}
			step++
		}

		if len(perms) != k {
			fmt.Fprintln(out, "NO")
			continue
		}

		fmt.Fprintln(out, "YES")
		for _, p := range perms {
			for i, v := range p {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, v)
			}
			fmt.Fprintln(out)
		}
	}
}
