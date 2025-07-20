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

	var t int
	fmt.Fscan(in, &t)

	for ; t > 0; t-- {
		var n, x int
		fmt.Fscan(in, &n, &x)

		ans := 1                              // corresponds to `ans = 1` in the C++ initializer
		s := map[int]struct{}{1: {}}          // initial set {1}

		for i := 0; i < n; i++ {
			var v int
			fmt.Fscan(in, &v)

			// build the new candidates in a separate map
			s2 := make(map[int]struct{})
			for val := range s {
				if x%(val*v) == 0 {
					s2[val*v] = struct{}{}
				}
			}

			// merge s2 into s
			for k := range s2 {
				s[k] = struct{}{}
			}

			// if we've reached x, start a new segment
			if _, ok := s[x]; ok {
				ans++
				s = map[int]struct{}{1: {}, v: {}}
			}
		}

		fmt.Fprintln(out, ans)
	}
}
