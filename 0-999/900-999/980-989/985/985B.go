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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	rows := make([][]byte, n)
	cnt := make([]int, m)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		rows[i] = []byte(s)
		for j := 0; j < m; j++ {
			if rows[i][j] == '1' {
				cnt[j]++
			}
		}
	}
	for i := 0; i < n; i++ {
		ok := true
		for j := 0; j < m; j++ {
			if cnt[j]-int(rows[i][j]-'0') == 0 {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
			return
		}
	}
	fmt.Fprintln(writer, "NO")
}
