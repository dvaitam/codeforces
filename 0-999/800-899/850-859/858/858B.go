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

	type Pair struct{ k, f int }
	mem := make([]Pair, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &mem[i].k, &mem[i].f)
	}

	floors := make(map[int]struct{})
	for x := 1; x <= 100; x++ {
		ok := true
		for _, p := range mem {
			if (p.k-1)/x+1 != p.f {
				ok = false
				break
			}
		}
		if ok {
			fl := (n-1)/x + 1
			floors[fl] = struct{}{}
		}
	}

	if len(floors) == 1 {
		for fl := range floors {
			fmt.Fprintln(out, fl)
		}
	} else {
		fmt.Fprintln(out, -1)
	}
}
