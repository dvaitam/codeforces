package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func decode(b string) string {
	set := make(map[rune]struct{})
	for _, ch := range b {
		set[ch] = struct{}{}
	}
	letters := make([]rune, 0, len(set))
	for ch := range set {
		letters = append(letters, ch)
	}
	sort.Slice(letters, func(i, j int) bool { return letters[i] < letters[j] })
	m := make(map[rune]rune)
	n := len(letters)
	for i, ch := range letters {
		m[ch] = letters[n-1-i]
	}
	res := make([]rune, len(b))
	for i, ch := range b {
		res[i] = m[ch]
	}
	return string(res)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var b string
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &b)
		fmt.Fprintln(out, decode(b))
	}
}
