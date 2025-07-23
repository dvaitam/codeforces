package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func charVal(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'A' && c <= 'Z':
		return int(c-'A') + 10
	case c >= 'a' && c <= 'z':
		return int(c-'a') + 36
	case c == '-':
		return 62
	case c == '_':
		return 63
	}
	return 0
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}

	var ways [64]int64
	for i := 0; i < 64; i++ {
		w := int64(1)
		for b := 0; b < 6; b++ {
			if (i & (1 << b)) == 0 {
				w = (w * 3) % MOD
			}
		}
		ways[i] = w
	}

	ans := int64(1)
	for i := 0; i < len(s); i++ {
		ans = (ans * ways[charVal(s[i])]) % MOD
	}

	fmt.Fprintln(writer, ans)
}
