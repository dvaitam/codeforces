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
		var a, b, k int
		fmt.Fscan(reader, &a, &b, &k)
		boys := make([]int, k)
		girls := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &boys[i])
		}
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &girls[i])
		}
		cntBoy := make([]int, a+1)
		cntGirl := make([]int, b+1)
		for i := 0; i < k; i++ {
			cntBoy[boys[i]]++
			cntGirl[girls[i]]++
		}
		total := int64(k) * int64(k-1) / 2
		for i := 1; i <= a; i++ {
			x := cntBoy[i]
			total -= int64(x) * int64(x-1) / 2
		}
		for i := 1; i <= b; i++ {
			x := cntGirl[i]
			total -= int64(x) * int64(x-1) / 2
		}
		fmt.Fprintln(writer, total)
	}
}
