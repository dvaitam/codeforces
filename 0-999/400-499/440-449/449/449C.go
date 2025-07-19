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
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	used := make([]bool, n+1)
	grp := make([][2]int, 0, n/2)

	// pair odd numbers by common divisors
	var p int
	for i := 3; i <= n/2; i += 2 {
		if !used[i] {
			p = i
			for j := i * 3; j <= n; j += i {
				if !used[j] {
					if p != 0 {
						grp = append(grp, [2]int{p, j})
						used[p], used[j] = true, true
						p = 0
					} else {
						p = j
					}
				}
			}
			if p != 0 {
				pair := i * 2
				if pair <= n && !used[p] && !used[pair] {
					grp = append(grp, [2]int{p, pair})
					used[p], used[pair] = true, true
				}
			}
		}
	}
	// pair even numbers
	p = 0
	for i := 2; i <= n; i += 2 {
		if !used[i] {
			if p != 0 {
				grp = append(grp, [2]int{p, i})
				used[p], used[i] = true, true
				p = 0
			} else {
				p = i
			}
		}
	}

	fmt.Fprintln(writer, len(grp))
	for _, pr := range grp {
		fmt.Fprintf(writer, "%d %d\n", pr[0], pr[1])
	}
}
