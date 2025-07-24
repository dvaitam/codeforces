package main

import (
	"bufio"
	"fmt"
	"os"
)

func generateNames() []string {
	names := make([]string, 0, 52)
	for i := 0; i < 52; i++ {
		first := 'A' + rune(i/26)
		second := 'a' + rune(i%26)
		names = append(names, fmt.Sprintf("%c%c", first, second))
	}
	return names
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	m := n - k + 1
	s := make([]string, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &s[i])
	}

	pool := generateNames()
	res := make([]string, n)
	idx := 0

	for i := 0; i < k-1; i++ {
		res[i] = pool[idx]
		idx++
	}
	for i := 0; i < m; i++ {
		if s[i] == "YES" {
			res[i+k-1] = pool[idx]
			idx++
		} else {
			res[i+k-1] = res[i]
		}
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, res[i])
	}
	fmt.Fprintln(out)
}
