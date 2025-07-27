package main

import (
	"bufio"
	"fmt"
	"os"
)

func simulate(s string, firstBias bool) int {
	first, second := 0, 0
	remainingFirst, remainingSecond := 5, 5
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			remainingFirst--
			if s[i] == '1' || (s[i] == '?' && firstBias) {
				first++
			}
		} else {
			remainingSecond--
			if s[i] == '1' || (s[i] == '?' && !firstBias) {
				second++
			}
		}
		if first > second+remainingSecond {
			return i + 1
		}
		if second > first+remainingFirst {
			return i + 1
		}
	}
	return 10
}

func solve(s string) int {
	res1 := simulate(s, true)
	res2 := simulate(s, false)
	if res1 < res2 {
		return res1
	}
	return res2
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var s string
		if _, err := fmt.Fscan(reader, &s); err != nil {
			return
		}
		fmt.Fprintln(writer, solve(s))
	}
}
