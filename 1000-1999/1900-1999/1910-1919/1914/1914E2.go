package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	a   int64
	b   int64
	sum int64
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}

		diff := int64(0)
		cont := make([]pair, 0)
		for i := 0; i < n; i++ {
			if a[i] > 0 && b[i] > 0 {
				cont = append(cont, pair{a[i], b[i], a[i] + b[i]})
			} else if a[i] > 0 {
				diff += a[i]
			} else {
				diff -= b[i]
			}
		}

		sort.Slice(cont, func(i, j int) bool {
			return cont[i].sum > cont[j].sum
		})

		for idx, p := range cont {
			if idx%2 == 0 {
				diff += p.a
			} else {
				diff -= p.b
			}
		}
		if len(cont)%2 == 1 {
			diff -= 1
		}
		fmt.Fprintln(writer, diff)
	}
}
