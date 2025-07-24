package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}

	sort.Ints(a)
	sort.Ints(b)

	best := m
	for i := 0; i < n; i++ {
		x := (b[i] - a[0]) % m
		if x < 0 {
			x += m
		}
		tmp := make([]int, n)
		for j := 0; j < n; j++ {
			val := (a[j] + x) % m
			if val < 0 {
				val += m
			}
			tmp[j] = val
		}
		sort.Ints(tmp)
		equal := true
		for j := 0; j < n; j++ {
			if tmp[j] != b[j] {
				equal = false
				break
			}
		}
		if equal && x < best {
			best = x
		}
	}

	fmt.Fprintln(writer, best)
}
