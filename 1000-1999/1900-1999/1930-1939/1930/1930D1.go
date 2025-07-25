package main

import (
	"bufio"
	"fmt"
	"os"
)

func minOnes(p string) int {
	n := len(p)
	covered := make([]bool, n)
	ans := 0
	for i := 0; i < n; i++ {
		if p[i] == '1' && !covered[i] {
			pos := i + 1
			if pos >= n {
				pos = i
			}
			ans++
			for j := pos - 1; j <= pos+1; j++ {
				if j >= 0 && j < n {
					covered[j] = true
				}
			}
		}
	}
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(reader, &n)
		fmt.Fscan(reader, &s)
		total := 0
		for i := 0; i < n; i++ {
			for j := i + 1; j <= n; j++ {
				total += minOnes(s[i:j])
			}
		}
		fmt.Println(total)
	}
}
