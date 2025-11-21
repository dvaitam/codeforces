package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int = 1_000_000_007
const inv2 int = (mod + 1) / 2

var memo = map[string]int{"1": 0}

func normalize(s string) string {
	i := 0
	for i < len(s) && s[i] == '0' {
		i++
	}
	if i == len(s) {
		return "0"
	}
	return s[i:]
}

func normalizeBytes(b []byte) string {
	i := 0
	for i < len(b) && b[i] == '0' {
		i++
	}
	if i == len(b) {
		return "0"
	}
	return string(b[i:])
}

func floorDiv2(s string) string {
	if len(s) <= 1 {
		return "0"
	}
	return normalize(s[:len(s)-1])
}

func ceilDiv2(s string) string {
	b := []byte(s)
	carry := byte(1)
	for i := len(b) - 1; i >= 0 && carry == 1; i-- {
		if b[i] == '1' {
			b[i] = '0'
		} else {
			b[i] = '1'
			carry = 0
		}
	}
	if carry == 1 {
		b = append([]byte{'1'}, b...)
	}
	return normalizeBytes(b[:len(b)-1])
}

func solve(s string) int {
	if val, ok := memo[s]; ok {
		return val
	}
	var res int
	if s[len(s)-1] == '0' {
		child := floorDiv2(s)
		res = (1 + solve(child)) % mod
	} else {
		down := floorDiv2(s)
		up := ceilDiv2(s)
		sum := solve(down) + solve(up)
		sum %= mod
		res = (1 + sum*inv2%mod) % mod
	}
	memo[s] = res
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)
		answer := solve(s)
		fmt.Fprintln(out, answer)
	}
}
