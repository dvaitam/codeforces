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
		var n int
		fmt.Fscan(reader, &n)
		needU, needR, needD, needL := false, false, false, false
		for i := 0; i < n; i++ {
			var x, y int
			fmt.Fscan(reader, &x, &y)
			if y > 0 {
				needU = true
			}
			if y < 0 {
				needD = true
			}
			if x > 0 {
				needR = true
			}
			if x < 0 {
				needL = true
			}
		}
		cnt := 0
		if needU {
			cnt++
		}
		if needR {
			cnt++
		}
		if needD {
			cnt++
		}
		if needL {
			cnt++
		}
		if cnt <= 3 {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
