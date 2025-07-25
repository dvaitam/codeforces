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
		var a, b, l int
		fmt.Fscan(in, &a, &b, &l)
		kset := make(map[int]struct{})
		for powA := 1; powA <= l; {
			for powB := 1; powA*powB <= l; {
				if l%(powA*powB) == 0 {
					k := l / (powA * powB)
					kset[k] = struct{}{}
				}
				if powB > l/b {
					break
				}
				powB *= b
			}
			if powA > l/a {
				break
			}
			powA *= a
		}
		fmt.Fprintln(out, len(kset))
	}
}
