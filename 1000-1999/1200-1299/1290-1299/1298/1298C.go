package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var s string
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	fmt.Fscan(reader, &s)
	ans := 0
	cnt := 0
	for i := 0; i < len(s); i++ {
		if s[i] == 'x' {
			cnt++
			if cnt >= 3 {
				ans++
			}
		} else {
			cnt = 0
		}
	}
	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, ans)
	writer.Flush()
}
