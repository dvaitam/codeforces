package main

import (
	"bufio"
	"fmt"
	"os"
)

func nextQ(s []byte, q []int, c byte) int {
	j := q[len(s)-1]
	for j > 0 && s[j] == c {
		j = q[j-1]
	}
	if s[j] != c {
		return j + 1
	}
	return 0
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		s := []byte{'0'}
		q := []int{0}
		for i := 2; i <= n; i++ {
			fmt.Fprintf(out, "2 %d\n", i)
			out.Flush()
			var ans int
			fmt.Fscan(in, &ans)
			q0 := nextQ(s, q, '0')
			q1 := nextQ(s, q, '1')
			if ans == q0 {
				s = append(s, '0')
				q = append(q, q0)
			} else {
				s = append(s, '1')
				q = append(q, q1)
			}
		}
		fmt.Fprintf(out, "0 %s\n", string(s))
		out.Flush()
		var verdict int
		fmt.Fscan(in, &verdict)
		if verdict == -1 {
			return
		}
	}
}
