package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)
	total := 4 * n
	a := make([]int, total+1)
	b := make([]int, total+1)
	c := make([]int, total+1)
	for i := 1; i <= total; i++ {
		fmt.Fscan(reader, &a[i], &b[i], &c[i])
	}
	s := make([]int, n+1)
	used := make([]bool, total+1)
	order := make([]int, total)
	for i := 0; i < total; i++ {
		order[i] = i + 1
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(order), func(i, j int) {
		order[i], order[j] = order[j], order[i]
	})
	ans := make([]int, 0, total)
	add := func(i, sign int) {
		s[a[i]] -= sign
		s[b[i]] += sign
		s[c[i]] += sign
	}
	fit := func(i int) bool {
		add(i, 1)
		if s[b[i]] > 5 || s[c[i]] > 5 {
			add(i, -1)
			return false
		}
		return true
	}
	for len(ans) < total {
		for _, j := range order {
			if !used[j] && fit(j) {
				used[j] = true
				ans = append(ans, j)
			}
		}
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fprintln(writer, "YES")
	for i, x := range ans {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprintf(writer, "%d", x)
	}
	writer.WriteByte('\n')
}
