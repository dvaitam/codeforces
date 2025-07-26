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
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	if n <= 2 {
		fmt.Fprintln(writer, 0)
		return
	}

	const B = 316
	ans := 1

	// handle small differences |d| <= B
	for d := -B; d <= B; d++ {
		m := make(map[int]int)
		for i := 1; i <= n; i++ {
			val := a[i] - (i-1)*d
			if cnt := m[val] + 1; cnt > ans {
				ans = cnt
			}
			m[val]++
		}
	}

	// handle large differences using local enumeration
	for i := 1; i <= n; i++ {
		m := make(map[int]int)
		for j := i + 1; j <= n && j <= i+B; j++ {
			diff := a[j] - a[i]
			den := j - i
			if diff%den == 0 {
				d := diff / den
				if d > B || d < -B {
					if cnt := m[d] + 2; cnt > ans {
						ans = cnt
					}
					m[d]++
				}
			}
		}
	}

	fmt.Fprintln(writer, n-ans)
}
