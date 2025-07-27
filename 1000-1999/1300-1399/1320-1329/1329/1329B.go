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
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var d, m int64
		fmt.Fscan(reader, &d, &m)
		ans := int64(1)
		for bit := 0; (int64(1) << bit) <= d; bit++ {
			L := int64(1) << bit
			R := (int64(1) << (bit + 1)) - 1
			if R > d {
				R = d
			}
			cnt := R - L + 1
			ans = ans * (cnt + 1) % m
		}
		ans = (ans - 1 + m) % m
		fmt.Fprintln(writer, ans)
	}
}
