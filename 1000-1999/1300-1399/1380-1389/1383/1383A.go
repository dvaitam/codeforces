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
		var n int
		fmt.Fscan(in, &n)
		var sA, sB string
		fmt.Fscan(in, &sA)
		fmt.Fscan(in, &sB)
		a := []byte(sA)
		b := []byte(sB)
		impossible := false
		for i := 0; i < n; i++ {
			if a[i] > b[i] {
				impossible = true
				break
			}
		}
		if impossible {
			fmt.Fprintln(out, -1)
			continue
		}
		moves := 0
		for c := byte('a'); c <= byte('t'); c++ { // iterate 20 letters
			minY := byte('z')
			for i := 0; i < n; i++ {
				if a[i] == c && b[i] > c {
					if b[i] < minY {
						minY = b[i]
					}
				}
			}
			if minY != byte('z') {
				moves++
				for i := 0; i < n; i++ {
					if a[i] == c && b[i] > c {
						a[i] = minY
					}
				}
			}
		}
		fmt.Fprintln(out, moves)
	}
}
