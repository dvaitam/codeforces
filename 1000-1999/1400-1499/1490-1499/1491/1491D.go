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

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		if u > v {
			fmt.Fprintln(writer, "NO")
			continue
		}
		onesU, onesV := 0, 0
		ok := true
		for i := 0; i < 30; i++ {
			if (u>>i)&1 == 1 {
				onesU++
			}
			if (v>>i)&1 == 1 {
				onesV++
			}
			if onesV > onesU {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
