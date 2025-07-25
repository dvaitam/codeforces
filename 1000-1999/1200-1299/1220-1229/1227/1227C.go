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
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		var s string
		fmt.Fscan(in, &s)

		target := make([]byte, n)
		for i := 0; i < k-1; i++ {
			target[2*i] = '('
			target[2*i+1] = ')'
		}
		rem := n - 2*(k-1)
		half := rem / 2
		for i := 0; i < half; i++ {
			target[2*(k-1)+i] = '('
		}
		for i := half; i < rem; i++ {
			target[2*(k-1)+i] = ')'
		}

		b := []byte(s)
		ops := make([][2]int, 0)
		for i := 0; i < n; i++ {
			if b[i] == target[i] {
				continue
			}
			j := i + 1
			for j < n && b[j] != target[i] {
				j++
			}
			// reverse b[i..j]
			for l, r := i, j; l < r; l, r = l+1, r-1 {
				b[l], b[r] = b[r], b[l]
			}
			ops = append(ops, [2]int{i + 1, j + 1})
		}

		fmt.Fprintln(out, len(ops))
		for _, op := range ops {
			fmt.Fprintf(out, "%d %d\n", op[0], op[1])
		}
	}
}
