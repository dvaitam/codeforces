package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func parseVal(s string) int64 {
	var res int64
	for i := 0; i < len(s); i++ {
		res = res*10 + int64(s[i]-'0')
	}
	return res
}

func bestSplit(block string) int64 {
	l := len(block)
	var best int64
	for split := 1; split < l; split++ {
		left := parseVal(block[:split])
		right := parseVal(block[split:])
		sum := left + right
		if sum > best {
			best = sum
		}
	}
	return best
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		blocks := strings.Split(s, "+")
		m := len(blocks) - 1
		ans := parseVal(blocks[0]) + parseVal(blocks[m])
		for i := 1; i <= m-1; i++ {
			ans += bestSplit(blocks[i])
		}
		fmt.Fprintln(writer, ans)
	}
}
