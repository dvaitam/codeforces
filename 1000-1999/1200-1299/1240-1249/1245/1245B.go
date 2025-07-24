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
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		var a, b, c int
		fmt.Fscan(reader, &a, &b, &c)
		var s string
		fmt.Fscan(reader, &s)

		ans := make([]byte, n)
		wins := 0
		for i, ch := range s {
			switch ch {
			case 'R':
				if b > 0 {
					ans[i] = 'P'
					b--
					wins++
				}
			case 'P':
				if c > 0 {
					ans[i] = 'S'
					c--
					wins++
				}
			case 'S':
				if a > 0 {
					ans[i] = 'R'
					a--
					wins++
				}
			}
		}

		if wins < (n+1)/2 {
			fmt.Fprintln(writer, "NO")
			continue
		}

		for i := 0; i < n; i++ {
			if ans[i] == 0 {
				if a > 0 {
					ans[i] = 'R'
					a--
				} else if b > 0 {
					ans[i] = 'P'
					b--
				} else {
					ans[i] = 'S'
					c--
				}
			}
		}

		fmt.Fprintln(writer, "YES")
		fmt.Fprintln(writer, string(ans))
	}
}
