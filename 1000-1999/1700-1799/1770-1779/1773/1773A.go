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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		if n == 1 {
			fmt.Fprintln(out, "Impossible")
			continue
		}

		// build inverse permutation c = a^{-1}
		c := make([]int, n+1)
		for i := 0; i < n; i++ {
			c[a[i]] = i + 1
		}

		if n == 2 {
			if a[0] == 1 && a[1] == 2 {
				fmt.Fprintln(out, "Possible")
				fmt.Fprintln(out, "2 1")
				fmt.Fprintln(out, "2 1")
			} else {
				fmt.Fprintln(out, "Impossible")
			}
			continue
		}

		q := make([]int, n+1)
		order := make([]int, n)
		for i := 0; i < n; i++ {
			order[i] = i + 1
		}
		for i := 0; i < n; i++ {
			q[order[i]] = order[(i+1)%n]
		}

		conf := make([]int, 0)
		for i := 1; i <= n; i++ {
			if q[i] == c[i] {
				conf = append(conf, i)
			}
		}

		leftover := make([]int, 0)
		for len(conf) > 0 {
			u := conf[len(conf)-1]
			conf = conf[:len(conf)-1]
			foundIdx := -1
			for i := len(conf) - 1; i >= 0; i-- {
				v := conf[i]
				if c[v] != u && c[u] != v {
					foundIdx = i
					break
				}
			}
			if foundIdx == -1 {
				leftover = append(leftover, u)
				continue
			}
			v := conf[foundIdx]
			conf = append(conf[:foundIdx], conf[foundIdx+1:]...)
			q[u], q[v] = q[v], q[u]
		}

		leftover = append(leftover, conf...)

		if len(leftover) == 1 {
			i := leftover[0]
			found := -1
			for j := 1; j <= n; j++ {
				if j == i || j == c[i] {
					continue
				}
				if q[j] == i || q[j] == c[i] {
					continue
				}
				found = j
				break
			}
			if found == -1 {
				fmt.Fprintln(out, "Impossible")
				continue
			}
			j := found
			q[i], q[j] = q[j], q[i]
		} else if len(leftover) > 1 {
			// try to rotate leftovers
			for idx := 0; idx < len(leftover); idx++ {
				i := leftover[idx]
				found := -1
				for j := idx + 1; j < len(leftover); j++ {
					v := leftover[j]
					if c[v] != i && c[i] != v {
						found = j
						break
					}
				}
				if found != -1 {
					v := leftover[found]
					leftover = append(leftover[:found], leftover[found+1:]...)
					q[i], q[v] = q[v], q[i]
					idx--
				}
			}
		}

		valid := true
		for i := 1; i <= n; i++ {
			if q[i] == i || q[i] == c[i] {
				valid = false
				break
			}
		}

		if !valid {
			fmt.Fprintln(out, "Impossible")
			continue
		}

		invQ := make([]int, n+1)
		for i := 1; i <= n; i++ {
			invQ[q[i]] = i
		}

		p := make([]int, n+1)
		valid = true
		for j := 1; j <= n; j++ {
			p[j] = c[invQ[j]]
			if p[j] == j {
				valid = false
				break
			}
		}

		if !valid {
			fmt.Fprintln(out, "Impossible")
			continue
		}

		fmt.Fprintln(out, "Possible")
		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, p[i])
		}
		fmt.Fprintln(out)
		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, q[i])
		}
		fmt.Fprintln(out)
	}
}
