package main

import (
	"bufio"
	"fmt"
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

	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	freq := make([]int, n+1)
	for _, v := range a {
		if v >= 1 && v <= n {
			freq[v]++
		}
	}

	missing := make([]int, 0)
	for v := 1; v <= n; v++ {
		if freq[v] == 0 {
			missing = append(missing, v)
		}
	}

	missIdx := 0
	used := make([]bool, n+1)
	changes := 0

	for i := 0; i < n; i++ {
		v := a[i]
		if v < 1 || v > n {
			// out of range, replace immediately
			a[i] = missing[missIdx]
			missIdx++
			changes++
			continue
		}
		if !used[v] {
			if freq[v] == 1 {
				used[v] = true
				continue
			}
			if missIdx < len(missing) && missing[missIdx] < v {
				// better to replace
				freq[v]--
				a[i] = missing[missIdx]
				missIdx++
				changes++
			} else {
				used[v] = true
				freq[v]--
			}
		} else {
			// already used this value
			freq[v]--
			a[i] = missing[missIdx]
			missIdx++
			changes++
		}
	}

	fmt.Fprintln(writer, changes)
	for i, v := range a {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, v)
	}
	fmt.Fprintln(writer)
}
