package main

import (
	"bufio"
	"fmt"
	"os"
)

func calc(h, p int64) int64 {
	if p == 1 {
		return (int64(1) << h) - 1
	}
	var res int64
	cur := int64(1)
	for cur < p && res < h {
		res++
		cur *= 2
	}
	if res >= h {
		return h
	}
	remaining := (int64(1) << h) - (int64(1) << res)
	res += (remaining + p - 1) / p
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for i := 0; i < t; i++ {
		var h, p int64
		fmt.Fscan(reader, &h, &p)
		fmt.Fprintln(writer, calc(h, p))
	}
}
