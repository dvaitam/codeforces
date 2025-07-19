package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n int
	fmt.Fscan(reader, &n)
	a := make([]int, n)
	var z []int
	negc := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] == 0 {
			z = append(z, i)
		} else if a[i] < 0 {
			negc++
		}
	}
	// find negative to remove if odd count: the one with max value (closest to zero)
	p := -1
	for i := 0; i < n; i++ {
		if a[i] < 0 && (p < 0 || a[i] > a[p]) {
			p = i
		}
	}
	if negc%2 == 0 {
		p = -1
	}
	// operations
	// type 1 x y, type 2 x
	if len(z) > 0 {
		// merge zeros
		for i := 1; i < len(z); i++ {
			fmt.Fprintf(writer, "1 %d %d\n", z[i]+1, z[0]+1)
		}
		if len(z) != n {
			if len(z)+1 == n && negc == 1 {
				// only one non-zero and it's negative
				for i := 0; i < n; i++ {
					if a[i] != 0 {
						fmt.Fprintf(writer, "1 %d %d\n", i+1, z[0]+1)
					}
				}
			} else {
				if p >= 0 {
					fmt.Fprintf(writer, "1 %d %d\n", p+1, z[0]+1)
				}
				// remove zero
				fmt.Fprintf(writer, "2 %d\n", z[0]+1)
				// find f: any non-zero and not p
				f := -1
				for i := 0; i < n; i++ {
					if a[i] != 0 && i != p {
						f = i
						break
					}
				}
				for i := 0; i < n; i++ {
					if a[i] != 0 && i != f && i != p {
						fmt.Fprintf(writer, "1 %d %d\n", i+1, f+1)
					}
				}
			}
		}
	} else {
		if p < 0 {
			for i := 1; i < n; i++ {
				fmt.Fprintf(writer, "1 %d %d\n", i+1, 1)
			}
		} else {
			// remove p
			fmt.Fprintf(writer, "2 %d\n", p+1)
			// choose f
			f := 0
			if p == 0 {
				f = 1
			}
			for i := 0; i < n; i++ {
				if i != f && i != p {
					fmt.Fprintf(writer, "1 %d %d\n", i+1, f+1)
				}
			}
		}
	}
}
