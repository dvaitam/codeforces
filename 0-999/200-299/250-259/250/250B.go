package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		fmt.Fprintln(out, expandIPv6(s))
	}
}

func expandIPv6(short string) string {
	var blocks []string
	if strings.Contains(short, "::") {
		parts := strings.SplitN(short, "::", 2)
		leftStr, rightStr := parts[0], parts[1]
		var left, right []string
		if len(leftStr) > 0 {
			left = strings.Split(leftStr, ":")
		}
		if len(rightStr) > 0 {
			right = strings.Split(rightStr, ":")
		}
		missing := 8 - (len(left) + len(right))
		blocks = make([]string, 0, 8)
		blocks = append(blocks, left...)
		for i := 0; i < missing; i++ {
			blocks = append(blocks, "0")
		}
		blocks = append(blocks, right...)
	} else {
		blocks = strings.Split(short, ":")
	}

	for i := range blocks {
		blocks[i] = padBlock(blocks[i])
	}
	return strings.Join(blocks, ":")
}

func padBlock(block string) string {
	block = strings.ToLower(block)
	if len(block) == 0 {
		return "0000"
	}
	if len(block) < 4 {
		block = strings.Repeat("0", 4-len(block)) + block
	}
	return block
}
