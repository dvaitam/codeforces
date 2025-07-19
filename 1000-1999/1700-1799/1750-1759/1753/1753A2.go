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

	var T int
	fmt.Fscan(reader, &T)
	for T > 0 {
		T--
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n+2)
		sum := 0
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
			if a[i] < 0 {
				sum -= a[i]
			} else {
				sum += a[i]
			}
		}
		if sum%2 != 0 {
			writer.WriteString("-1\n")
			continue
		}
		p := make([]int, 0, n+2)
		p = append(p, 1)
		i := 1
		rev := 0
		for {
			for i <= n && a[i] == 0 {
				i++
			}
			if i > n {
				break
			}
			l := i
			i++
			for i <= n && a[i] == 0 {
				i++
			}
			if i > n {
				break
			}
			r := i
			i++
			if (r-l)%2 == 1 {
				if a[l] == a[r] {
					continue
				}
				sgnl := (l % 2) ^ rev
				if sgnl != 0 {
					if p[len(p)-1] != r {
						p = append(p, r)
					}
					p = append(p, r+1)
				} else {
					if p[len(p)-1] != l {
						p = append(p, l)
					}
					p = append(p, l+1)
				}
			} else {
				if a[l] != a[r] {
					continue
				}
				sgnl := (l % 2) ^ rev
				if sgnl != 0 {
					if p[len(p)-1] != r {
						p = append(p, r-1)
					}
					p = append(p, r+1)
					rev ^= 1
				} else {
					if p[len(p)-1] != l {
						p = append(p, l)
					}
					p = append(p, l+1)
				}
			}
		}
		if p[len(p)-1] != n+1 {
			p = append(p, n+1)
		}
		m := len(p) - 1
		fmt.Fprintln(writer, m)
		for j := 0; j < m; j++ {
			fmt.Fprintln(writer, p[j], p[j+1]-1)
		}
	}
}
