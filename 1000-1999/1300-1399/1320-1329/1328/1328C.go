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
		var n int
		fmt.Fscan(in, &n)
		var x string
		fmt.Fscan(in, &x)
		a := make([]byte, n)
		b := make([]byte, n)
		// first digit guaranteed to be '2'
		a[0], b[0] = '1', '1'
		broken := false
		for i := 1; i < n; i++ {
			switch x[i] {
			case '0':
				a[i], b[i] = '0', '0'
			case '1':
				if !broken {
					a[i], b[i] = '1', '0'
					broken = true
				} else {
					a[i], b[i] = '0', '1'
				}
			case '2':
				if !broken {
					a[i], b[i] = '1', '1'
				} else {
					a[i], b[i] = '0', '2'
				}
			}
		}
		fmt.Fprintln(out, string(a))
		fmt.Fprintln(out, string(b))
	}
}
