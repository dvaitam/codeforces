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
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	fibSet := make(map[int]struct{})
	a, b := 1, 1
	for a <= n {
		fibSet[a] = struct{}{}
		a, b = b, a+b
	}

	result := make([]byte, n)
	for i := 1; i <= n; i++ {
		if _, ok := fibSet[i]; ok {
			result[i-1] = 'O'
		} else {
			result[i-1] = 'o'
		}
	}
	fmt.Fprintln(out, string(result))
}
