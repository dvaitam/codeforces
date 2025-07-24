package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func digitSumStr(s string) int {
	sum := 0
	for i := 0; i < len(s); i++ {
		sum += int(s[i] - '0')
	}
	return sum
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	x, _ := strconv.ParseInt(s, 10, 64)
	best := x
	bestSum := digitSumStr(s)

	for i := 0; i < len(s); i++ {
		if s[i] == '0' {
			continue
		}
		t := s[:i] + string(s[i]-1) + strings.Repeat("9", len(s)-i-1)
		num, _ := strconv.ParseInt(t, 10, 64)
		if num <= 0 {
			continue
		}
		sum := digitSumStr(t)
		if sum > bestSum || (sum == bestSum && num > best) {
			best = num
			bestSum = sum
		}
	}
	fmt.Println(best)
}
