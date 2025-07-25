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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int64
		fmt.Fscan(reader, &n, &m)
		ans := int64(0)
		for g := int64(1); g <= m; g++ {
			kmax := (n + g) / (g * g)
			kmin := int64(1)
			if g == 1 {
				kmin = 2
			}
			if kmax >= kmin {
				ans += kmax - kmin + 1
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
