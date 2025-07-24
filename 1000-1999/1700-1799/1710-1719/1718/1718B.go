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

	fib := []int64{1, 1}
	sum := []int64{1, 2}
	for sum[len(sum)-1] < 1e11 {
		next := fib[len(fib)-1] + fib[len(fib)-2]
		fib = append(fib, next)
		sum = append(sum, sum[len(sum)-1]+next)
	}

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var k int
		if _, err := fmt.Fscan(reader, &k); err != nil {
			return
		}
		counts := make([]int64, k)
		var total int64
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &counts[i])
			total += counts[i]
		}

		n := -1
		for i, v := range sum {
			if v == total {
				n = i
				break
			}
			if v > total {
				break
			}
		}
		if n == -1 {
			fmt.Fprintln(writer, "NO")
			continue
		}

		last := -1
		ok := true
		for i := n; i >= 0; i-- {
			idx := -1
			for j := 0; j < k; j++ {
				if j == last {
					continue
				}
				if counts[j] >= fib[i] {
					if idx == -1 || counts[j] > counts[idx] {
						idx = j
					}
				}
			}
			if idx == -1 {
				ok = false
				break
			}
			counts[idx] -= fib[i]
			last = idx
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
