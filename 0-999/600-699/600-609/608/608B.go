package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var a, b string
	if _, err := fmt.Fscan(reader, &a); err != nil {
		return
	}
	fmt.Fscan(reader, &b)
	n := len(a)
	m := len(b)
	pref := make([]int, m+1)
	for i := 0; i < m; i++ {
		pref[i+1] = pref[i]
		if b[i] == '1' {
			pref[i+1]++
		}
	}
	segLen := m - n + 1
	var ans int64
	for i := 0; i < n; i++ {
		l := i
		r := i + segLen - 1
		ones := pref[r+1] - pref[l]
		if a[i] == '0' {
			ans += int64(ones)
		} else {
			ans += int64(segLen - ones)
		}
	}
	fmt.Println(ans)
}
