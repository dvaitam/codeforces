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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	ask := func(a, b int) int {
		fmt.Fprintf(out, "? %d %d\n", a, b)
		out.Flush()
		var g int
		if _, err := fmt.Fscan(in, &g); err != nil {
			os.Exit(0)
		}
		return g
	}

	for ; t > 0; t-- {
		ans := 0
		for i := 0; i < 30; i++ {
			a := (1 << i) - ans
			b := 3*(1<<i) - ans
			g := ask(a, b)
			if g%(1<<(i+1)) == 0 {
				ans |= 1 << i
			}
		}
		fmt.Fprintf(out, "! %d\n", ans)
		out.Flush()
	}
}
