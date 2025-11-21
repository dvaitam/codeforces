package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func factorial(n int) int {
	res := 1
	for i := 2; i <= n; i++ {
		res *= i
	}
	return res
}

// nextPermutation generates the next lexicographic permutation of p.
// Returns false if it was the last permutation.
func nextPermutation(p []int) bool {
	n := len(p)
	i := n - 2
	for i >= 0 && p[i] >= p[i+1] {
		i--
	}
	if i < 0 {
		return false
	}
	j := n - 1
	for p[j] <= p[i] {
		j--
	}
	p[i], p[j] = p[j], p[i]
	for l, r := i+1, n-1; l < r; l, r = l+1, r-1 {
		p[l], p[r] = p[r], p[l]
	}
	return true
}

func permKey(p []int) string {
	var b strings.Builder
	for i, v := range p {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.Itoa(v))
	}
	return b.String()
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

		// Basic impossibility checks.
		if n == 1 {
			if k == 1 {
				fmt.Fprintln(out, "YES")
				fmt.Fprintln(out, 1)
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
			if k > factorial(n) {
				fmt.Fprintln(out, "NO")
				continue
			}
		}

		perms := make([][]int, 0, k)
		used := make(map[string]struct{})
		remaining := k

		// For odd k (and odd n with k >= n), use n cyclic shifts.
		if remaining%2 == 1 {
			for shift := 0; shift < n; shift++ {
				p := make([]int, n)
				for i := 0; i < n; i++ {
					p[i] = (i+shift)%n + 1
				}
				perms = append(perms, p)
				used[permKey(p)] = struct{}{}
			}
			remaining -= n
		}

		pairsNeeded := remaining / 2
		cur := make([]int, n)
		for i := 0; i < n; i++ {
			cur[i] = i + 1
		}

		for pairsNeeded > 0 {
			key := permKey(cur)
			if _, ok := used[key]; !ok {
				comp := make([]int, n)
				for i, v := range cur {
					comp[i] = n + 1 - v
				}
				ckey := permKey(comp)
				if _, ok2 := used[ckey]; !ok2 {
					perms = append(perms, append([]int(nil), cur...))
					perms = append(perms, comp)
					used[key] = struct{}{}
					used[ckey] = struct{}{}
					pairsNeeded--
					if pairsNeeded == 0 {
						break
					}
				}
			}
			if !nextPermutation(cur) {
				break
			}
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
