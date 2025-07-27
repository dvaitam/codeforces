package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func canClear(a []int) bool {
	n := len(a)
	s := make([]int, n+1)
	prefOk := make([]bool, n+1)
	prefOk[0] = true
	for i := 1; i <= n; i++ {
		if i%2 == 1 {
			s[i] = s[i-1] + a[i-1]
		} else {
			s[i] = s[i-1] - a[i-1]
		}
		if i%2 == 1 {
			prefOk[i] = prefOk[i-1] && s[i] >= 0
		} else {
			prefOk[i] = prefOk[i-1] && s[i] <= 0
		}
	}
	if prefOk[n] && s[n] == 0 {
		return true
	}
	const INF int = int(1e18)
	oddMin := make([]int, n+2)
	evenMax := make([]int, n+2)
	oddMin[n+1] = INF
	evenMax[n+1] = -INF
	for i := n; i >= 0; i-- {
		if i%2 == 1 {
			if oddMin[i+1] != INF {
				oddMin[i] = min(oddMin[i+1], s[i])
			} else {
				oddMin[i] = s[i]
			}
			evenMax[i] = evenMax[i+1]
		} else {
			if evenMax[i+1] != -INF {
				evenMax[i] = max(evenMax[i+1], s[i])
			} else {
				evenMax[i] = s[i]
			}
			oddMin[i] = oddMin[i+1]
		}
	}
	for i := 1; i < n; i++ {
		if !prefOk[i-1] {
			continue
		}
		var nsI, nsIp1 int
		if i%2 == 1 {
			nsI = s[i-1] + a[i]
			nsIp1 = nsI - a[i-1]
			if nsI < 0 || nsIp1 > 0 {
				continue
			}
		} else {
			nsI = s[i-1] - a[i]
			nsIp1 = nsI + a[i-1]
			if nsI > 0 || nsIp1 < 0 {
				continue
			}
		}
		delta := nsIp1 - s[i+1]
		if delta != -s[n] {
			continue
		}
		if i+2 <= n {
			if oddMin[i+2] != INF && oddMin[i+2] < s[n] {
				continue
			}
			if evenMax[i+2] != -INF && evenMax[i+2] > s[n] {
				continue
			}
		}
		return true
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		if canClear(a) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
