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
	// positions and spit distances
	pos := make(map[int]int, n)
	x := make([]int, n)
	d := make([]int, n)
	for i := 0; i < n; i++ {
		var xi, di int
		fmt.Fscan(reader, &xi, &di)
		x[i] = xi
		d[i] = di
		pos[xi] = i
	}
	// check for mutual spit
	for i := 0; i < n; i++ {
		target := x[i] + d[i]
		if j, ok := pos[target]; ok {
			if x[j]+d[j] == x[i] {
				fmt.Fprintln(writer, "YES")
				return
			}
		}
	}
	fmt.Fprintln(writer, "NO")
}
