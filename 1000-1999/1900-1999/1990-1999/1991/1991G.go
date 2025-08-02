package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var t, n, m, k, q int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &n, &m, &k, &q, &s)
		if n == k && m == k {
			for i := 0; i < q; i++ {
				fmt.Println("1 1")
			}
		} else if n == k {
			nw := 0
			for _, c := range s {
				if c == 'H' {
					nw++
					fmt.Printf("%d 1\n", nw)
					nw %= n
				} else {
					fmt.Printf("1 %d\n", m)
				}
			}
		} else if m == k {
			nw := 0
			for _, c := range s {
				if c == 'V' {
					nw++
					fmt.Printf("1 %d\n", nw)
					nw %= m
				} else {
					fmt.Printf("%d 1\n", n)
				}
			}
		} else {
			x, y, u, v := n, m, k, k
			for _, c := range s {
				if c == 'H' {
					fmt.Printf("%d 1\n", x)
					x--
					if x == 0 {
						if y <= k {
							v, y, x = y, m, k
						} else {
							x = n
						}
					} else if x == k {
						x, u = u, k
					}
				} else {
					fmt.Printf("1 %d\n", y)
					y--
					if y == 0 {
						if x <= k {
							u, x, y = x, n, k
						} else {
							y = m
						}
					} else if y == k {
						y, v = v, k
					}
				}
			}
		}
	}
}
