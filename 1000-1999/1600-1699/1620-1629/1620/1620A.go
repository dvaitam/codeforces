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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var s string
		fmt.Fscan(in, &s)
		cnt := 0
		for i := 0; i < len(s); i++ {
			if s[i] == 'N' {
				cnt++
			}
		}
		if cnt == 1 {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
		}
	}
}
