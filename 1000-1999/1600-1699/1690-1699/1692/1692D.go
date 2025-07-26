package main

import (
	"bufio"
	"fmt"
	"os"
)

func isPalindrome(minutes int) bool {
	h := minutes / 60
	m := minutes % 60
	return h/10 == m%10 && h%10 == m/10
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var tt int
	fmt.Fscan(reader, &tt)
	for ; tt > 0; tt-- {
		var s string
		var x int
		fmt.Fscan(reader, &s, &x)
		h := int(s[0]-'0')*10 + int(s[1]-'0')
		m := int(s[3]-'0')*10 + int(s[4]-'0')
		start := h*60 + m
		seen := make(map[int]bool)
		t := start
		count := 0
		for !seen[t] {
			seen[t] = true
			if isPalindrome(t) {
				count++
			}
			t = (t + x) % 1440
		}
		fmt.Fprintln(writer, count)
	}
}
