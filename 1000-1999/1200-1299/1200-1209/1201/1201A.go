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

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	answers := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &answers[i])
	}
	scores := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &scores[i])
	}

	total := 0
	for j := 0; j < m; j++ {
		cnt := [5]int{}
		for i := 0; i < n; i++ {
			idx := int(answers[i][j] - 'A')
			if idx >= 0 && idx < 5 {
				cnt[idx]++
			}
		}
		max := 0
		for k := 0; k < 5; k++ {
			if cnt[k] > max {
				max = cnt[k]
			}
		}
		total += max * scores[j]
	}
	fmt.Fprintln(out, total)
}
