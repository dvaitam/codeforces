package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	fmt.Fscan(in, &q)

	// generate n distinct magic words consisting of X only
	words := make([]string, n)
	for i := 0; i < n; i++ {
		length := 1 << uint(i%20) // limit growth
		words[i] = strings.Repeat("X", length)
		fmt.Fprintln(out, words[i])
	}

	powerMap := make(map[int][2]int)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			p := len(words[i]) + len(words[j])
			if _, ok := powerMap[p]; !ok {
				powerMap[p] = [2]int{i + 1, j + 1}
			}
		}
	}

	for k := 0; k < q; k++ {
		var p int
		fmt.Fscan(in, &p)
		if pair, ok := powerMap[p]; ok {
			fmt.Fprintln(out, pair[0], pair[1])
		} else {
			fmt.Fprintln(out, -1, -1)
		}
	}
}
