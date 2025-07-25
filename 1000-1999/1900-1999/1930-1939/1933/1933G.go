package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func color(p int, r, c int64) int {
	switch p {
	case 0:
		return int((r & 1) ^ ((c >> 1) & 1))
	case 1:
		return int((r & 1) ^ ((c >> 1) & 1) ^ 1)
	case 2:
		return int((c & 1) ^ ((r >> 1) & 1))
	case 3:
		return int((c & 1) ^ ((r >> 1) & 1) ^ 1)
	case 4:
		return int((c & 1) ^ (((r >> 1) ^ r) & 1))
	case 5:
		return int((c & 1) ^ (((r >> 1) ^ r) & 1) ^ 1)
	case 6:
		return int(((c >> 1) & 1) ^ (c & 1) ^ (r & 1))
	default: // case 7
		return int(((c >> 1) & 1) ^ (c & 1) ^ (r & 1) ^ 1)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m, q int64
		fmt.Fscan(reader, &n, &m, &q)
		valid := [8]bool{true, true, true, true, true, true, true, true}
		fmt.Fprintln(writer, 8)
		for ; q > 0; q-- {
			var r, c int64
			var shape string
			fmt.Fscan(reader, &r, &c, &shape)
			val := 0
			if shape == "square" {
				val = 1
			}
			for i := 0; i < 8; i++ {
				if valid[i] && color(i, r-1, c-1) != val {
					valid[i] = false
				}
			}
			cnt := 0
			for i := 0; i < 8; i++ {
				if valid[i] {
					cnt++
				}
			}
			fmt.Fprintln(writer, int64(cnt)%mod)
		}
	}
}
