package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
	}
	if m/2 < k {
		fmt.Println(m / 2)
	} else {
		fmt.Println(k)
	}
}
