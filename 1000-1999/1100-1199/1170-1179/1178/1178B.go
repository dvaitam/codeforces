package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	n := len(s)
	// pref[i] = number of "vv" pairs up to index i
	pref := make([]int64, n)
	for i := 1; i < n; i++ {
		pref[i] = pref[i-1]
		if s[i] == 'v' && s[i-1] == 'v' {
			pref[i]++
		}
	}
	total := int64(0)
	if n > 0 {
		total = pref[n-1]
	}
	var ans int64
	for i := 0; i < n; i++ {
		if s[i] == 'o' {
			left := pref[i]
			right := total - pref[i]
			ans += left * right
		}
	}
	fmt.Println(ans)
}
