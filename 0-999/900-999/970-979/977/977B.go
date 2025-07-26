package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)
	var s string
	fmt.Fscan(reader, &s)
	counts := make(map[string]int)
	for i := 0; i < n-1; i++ {
		pair := s[i : i+2]
		counts[pair]++
	}
	best := ""
	bestCnt := 0
	for i := 0; i < n-1; i++ {
		pair := s[i : i+2]
		if counts[pair] > bestCnt {
			bestCnt = counts[pair]
			best = pair
		}
	}
	fmt.Println(best)
}
