package main

import (
	"bufio"
	"fmt"
	"os"
)

func prefixFunction(s string) []int {
	pi := make([]int, len(s))
	for i := 1; i < len(s); i++ {
		j := pi[i-1]
		for j > 0 && s[i] != s[j] {
			j = pi[j-1]
		}
		if s[i] == s[j] {
			j++
		}
		pi[i] = j
	}
	return pi
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	if n <= 0 {
		return
	}

	var res string
	fmt.Fscan(reader, &res)

	for i := 1; i < n; i++ {
		var w string
		fmt.Fscan(reader, &w)
		l := len(w)
		if len(res) < l {
			l = len(res)
		}
		tail := res[len(res)-l:]
		s := w + "#" + tail
		pi := prefixFunction(s)
		overlap := pi[len(s)-1]
		res += w[overlap:]
	}

	fmt.Fprint(writer, res)
}
