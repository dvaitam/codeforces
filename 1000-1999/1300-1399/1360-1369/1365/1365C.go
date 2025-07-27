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

	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}

	pos := make([]int, n+1)
	for i := 0; i < n; i++ {
		pos[a[i]] = i
	}

	cnt := make(map[int]int)
	best := 0
	for i := 0; i < n; i++ {
		shift := (i - pos[b[i]] + n) % n
		cnt[shift]++
		if cnt[shift] > best {
			best = cnt[shift]
		}
	}

	fmt.Fprintln(out, best)
}
