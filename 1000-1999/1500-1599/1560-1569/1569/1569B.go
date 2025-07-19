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

	var t int
	fmt.Fscan(reader, &t)
	for t > 0 {
		t--
		var n int
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)

		one, two := 0, 0
		first, prev := -1, -1
		defeats := make([]int, n)
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				one++
			} else {
				two++
				if first < 0 {
					first = i
					prev = i
				} else {
					defeats[i] = prev
					defeats[first] = i
					prev = i
				}
			}
		}
		if two == 1 || two == 2 {
			fmt.Fprintln(writer, "NO")
			continue
		}
		fmt.Fprintln(writer, "YES")

		ans := make([][]rune, n)
		for i := 0; i < n; i++ {
			ans[i] = make([]rune, n)
			for j := 0; j < n; j++ {
				ans[i][j] = '?'
			}
		}
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				if i == j {
					ans[i][j] = 'X'
				} else if s[i] == '1' || s[j] == '1' {
					ans[i][j] = '='
					ans[j][i] = '='
				} else {
					ans[i][j] = '+'
					ans[j][i] = '-'
					if defeats[j] == i {
						ans[i][j], ans[j][i] = ans[j][i], ans[i][j]
					}
				}
			}
		}
		for i := 0; i < n; i++ {
			writer.WriteString(string(ans[i]))
			writer.WriteByte('\n')
		}
	}
}
